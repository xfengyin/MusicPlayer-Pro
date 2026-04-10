<template>
  <div class="search-bar">
    <n-input
      ref="inputRef"
      v-model:value="searchQuery"
      type="text"
      placeholder="搜索歌曲、歌手、专辑..."
      :input-props="{ autocomplete: 'off' }"
      @focus="onFocus"
      @blur="onBlur"
      @keydown.enter="onSearch"
      clearable
      class="search-input"
    >
      <template #prefix>
        <svg viewBox="0 0 24 24" width="20" height="20" fill="var(--text-subdued)"><path d="M10.533 1.279c-5.18 0-9.407 4.14-9.407 9.279s4.226 9.279 9.407 9.279c2.234 0 4.29-.77 5.907-2.058l4.353 4.353a1 1 0 101.414-1.414l-4.344-4.344a9.157 9.157 0 002.077-5.816c0-5.14-4.226-9.28-9.407-9.28zm-7.407 9.279c0-4.006 3.302-7.28 7.407-7.28s7.407 3.274 7.407 7.28-3.302 7.279-7.407 7.279-7.407-3.273-7.407-7.28z"/></svg>
      </template>
    </n-input>

    <!-- Loading -->
    <div v-if="isLoading" class="loading-indicator">
      <n-spin size="small" />
      <span>搜索中...</span>
    </div>

    <!-- Search History Dropdown -->
    <transition name="slide-fade">
      <div
        v-if="showHistory && searchHistory.length > 0 && isFocused && !hasQuery"
        class="search-dropdown"
      >
        <div class="dropdown-header">
          <span>搜索历史</span>
          <button class="clear-btn" @click="clearHistory">清空</button>
        </div>
        <div class="history-list">
          <div
            v-for="(item, index) in searchHistory"
            :key="index"
            class="history-item"
            @click="selectHistory(item)"
          >
            <svg viewBox="0 0 24 24" width="16" height="16" fill="var(--text-subdued)"><path d="M12 2C6.477 2 2 6.477 2 12s4.477 10 10 10 10-4.477 10-10S17.523 2 12 2zm0 18a8 8 0 110-16 8 8 0 010 16zm-.5-13a.5.5 0 01.5.5v4.5h3.5a.5.5 0 010 1H11.5a.5.5 0 01-.5-.5V7.5a.5.5 0 01.5-.5z"/></svg>
            <span>{{ item }}</span>
            <button class="remove-history" @click.stop="removeHistory(index)">
              <svg viewBox="0 0 24 24" width="14" height="14" fill="currentColor"><path d="M18 6L6 18M6 6l12 12" stroke="currentColor" stroke-width="2" stroke-linecap="round"/></svg>
            </button>
          </div>
        </div>
      </div>
    </transition>

    <!-- Search Results Dropdown -->
    <transition name="slide-fade">
      <div
        v-if="showResults && searchResults.length > 0 && isFocused"
        class="search-dropdown"
      >
        <div class="results-list">
          <div
            v-for="song in searchResults"
            :key="song.id"
            class="result-item"
            @click="selectResult(song)"
          >
            <div class="result-cover-wrapper">
              <img
                v-if="song.cover"
                :src="song.cover"
                :alt="song.title"
                class="result-cover"
                @error="onImageError"
              />
              <div v-else class="result-cover-placeholder">
                <svg viewBox="0 0 24 24" width="20" height="20" fill="currentColor"><path d="M12 3v10.55c-.59-.34-1.27-.55-2-.55C7.79 13 6 14.79 6 17s1.79 4 4 4 4-1.79 4-4V7h4V3h-6z"/></svg>
              </div>
              <div class="result-cover-play">
                <svg viewBox="0 0 24 24" width="14" height="14" fill="#000"><path d="M8 5.14v13.72a1 1 0 001.5.86l11.04-6.86a1 1 0 000-1.72L9.5 4.28a1 1 0 00-1.5.86z"/></svg>
              </div>
            </div>
            <div class="result-info">
              <div class="result-title">{{ song.title }}</div>
              <div class="result-artist">{{ song.artist }}</div>
            </div>
          </div>
        </div>
        <div v-if="hasMoreResults" class="show-more">
          <button class="show-more-btn" @click="onShowMore">查看全部结果</button>
        </div>
      </div>
    </transition>

    <!-- No Results -->
    <transition name="slide-fade">
      <div
        v-if="showResults && searchResults.length === 0 && !isLoading && hasQuery"
        class="search-dropdown"
      >
        <div class="no-results">
          <p>未找到相关歌曲</p>
          <span>试试其他关键词</span>
        </div>
      </div>
    </transition>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch, onMounted } from 'vue'
import { NInput, NButton, NSpin } from 'naive-ui'
import { searchMusic } from '../api/music'
import type { Song } from '../types/music'
import { usePlayerStore } from '../stores/player'

interface Props {
  maxHistory?: number
  debounceMs?: number
  maxResults?: number
}

const props = withDefaults(defineProps<Props>(), {
  maxHistory: 10,
  debounceMs: 300,
  maxResults: 10,
})

const emit = defineEmits<{
  search: [query: string]
  select: [song: Song]
  'history-clear': []
  'show-more': [query: string]
}>()

const inputRef = ref()
const searchQuery = ref('')
const isFocused = ref(false)
const showHistory = ref(false)
const showResults = ref(false)
const isLoading = ref(false)
const searchHistory = ref<string[]>([])
const searchResults = ref<Song[]>([])
const hasMoreResults = ref(false)

const playerStore = usePlayerStore()
let debounceTimer: ReturnType<typeof setTimeout> | null = null
const hasQuery = computed(() => searchQuery.value.trim().length > 0)

watch(searchQuery, (newQuery) => {
  if (debounceTimer) clearTimeout(debounceTimer)

  if (newQuery.trim().length > 0) {
    debounceTimer = setTimeout(() => {
      performSearch(newQuery)
    }, props.debounceMs)
  } else {
    searchResults.value = []
    showResults.value = false
    if (searchHistory.value.length > 0 && isFocused.value) {
      showHistory.value = true
    }
  }
})

function onFocus() {
  isFocused.value = true
  if (searchHistory.value.length > 0 && !hasQuery.value) {
    showHistory.value = true
  }
}

function onBlur() {
  setTimeout(() => {
    isFocused.value = false
    showHistory.value = false
  }, 200)
}

function onSearch() {
  const query = searchQuery.value.trim()
  if (query) {
    addToHistory(query)
    emit('search', query)
    showHistory.value = false
    performSearch(query)
  }
}

async function performSearch(query: string) {
  if (!query.trim()) return

  isLoading.value = true
  showResults.value = true
  showHistory.value = false

  try {
    const response = await searchMusic(query, 1, props.maxResults + 1)
    searchResults.value = response.songs.slice(0, props.maxResults)
    hasMoreResults.value = response.songs.length > props.maxResults
  } catch (error) {
    console.error('Search failed:', error)
    searchResults.value = []
  } finally {
    isLoading.value = false
  }
}

function selectHistory(query: string) {
  searchQuery.value = query
  onSearch()
}

function removeHistory(index: number) {
  searchHistory.value.splice(index, 1)
  saveHistory()
}

function clearHistory() {
  searchHistory.value = []
  saveHistory()
  emit('history-clear')
}

async function selectResult(song: Song) {
  playerStore.addToPlaylist(song)
  await playerStore.playSong(song)
  emit('select', song)
  searchQuery.value = ''
  showResults.value = false
}

function onShowMore() {
  const query = searchQuery.value.trim()
  if (query) emit('show-more', query)
}

function loadHistory() {
  try {
    const saved = localStorage.getItem('searchHistory')
    if (saved) searchHistory.value = JSON.parse(saved)
  } catch (e) {
    searchHistory.value = []
  }
}

function saveHistory() {
  try {
    localStorage.setItem('searchHistory', JSON.stringify(searchHistory.value))
  } catch (e) {
    console.error('Failed to save search history:', e)
  }
}

function addToHistory(query: string) {
  const index = searchHistory.value.indexOf(query)
  if (index > -1) searchHistory.value.splice(index, 1)
  searchHistory.value.unshift(query)
  if (searchHistory.value.length > props.maxHistory) searchHistory.value.pop()
  saveHistory()
}

function onImageError(e: Event) {
  const target = e.target as HTMLImageElement
  target.style.display = 'none'
}

onMounted(() => {
  loadHistory()
})
</script>

<style scoped>
.search-bar {
  position: relative;
  width: 100%;
  max-width: 500px;
}

.search-input {
  border-radius: var(--radius-pill) !important;
  background: var(--bg-elevated) !important;
}

.search-input :deep(.n-input) {
  border-radius: var(--radius-pill) !important;
}

.search-input :deep(.n-input-wrapper) {
  border-radius: var(--radius-pill) !important;
  padding: 8px 16px;
}

.loading-indicator {
  position: absolute;
  top: calc(100% + 8px);
  left: 0;
  right: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  padding: 12px;
  color: var(--text-subdued);
  font-size: var(--font-size-small);
}

/* Dropdown */
.search-dropdown {
  position: absolute;
  top: calc(100% + 8px);
  left: 0;
  right: 0;
  background: var(--bg-elevated);
  border-radius: var(--radius-lg);
  box-shadow: var(--shadow-dialog);
  z-index: 100;
  max-height: 400px;
  overflow-y: auto;
  border: none;
}

.dropdown-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 16px;
  border-bottom: 1px solid rgba(255, 255, 255, 0.06);
  font-size: var(--font-size-small);
  color: var(--text-subdued);
}

.clear-btn {
  background: none;
  border: none;
  color: var(--text-subdued);
  font-size: var(--font-size-small);
  cursor: pointer;
  font-weight: var(--font-weight-bold);
}

.clear-btn:hover {
  color: var(--text-base);
}

.history-list,
.results-list {
  padding: 4px 0;
}

.history-item,
.result-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 10px 16px;
  cursor: pointer;
  transition: background 0.2s;
}

.history-item:hover,
.result-item:hover {
  background: rgba(255, 255, 255, 0.08);
}

.history-item {
  color: var(--text-subdued);
  font-size: var(--font-size-caption);
}

.remove-history {
  margin-left: auto;
  background: none;
  border: none;
  cursor: pointer;
  color: var(--text-subdued);
  opacity: 0;
  transition: opacity 0.2s;
  padding: 4px;
}

.history-item:hover .remove-history {
  opacity: 1;
}

/* Result Cover */
.result-cover-wrapper {
  position: relative;
  width: 48px;
  height: 48px;
  flex-shrink: 0;
}

.result-cover {
  width: 100%;
  height: 100%;
  border-radius: var(--radius-md);
  object-fit: cover;
}

.result-cover-placeholder {
  width: 100%;
  height: 100%;
  border-radius: var(--radius-md);
  background: var(--bg-highlight);
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--text-subdued);
}

.result-cover-play {
  position: absolute;
  inset: 0;
  background: rgba(0, 0, 0, 0.6);
  border-radius: var(--radius-md);
  display: flex;
  align-items: center;
  justify-content: center;
  opacity: 0;
  transition: opacity 0.2s;
}

.result-item:hover .result-cover-play {
  opacity: 1;
}

.result-info {
  flex: 1;
  min-width: 0;
}

.result-title {
  font-size: var(--font-size-caption);
  color: var(--text-base);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  font-weight: var(--font-weight-regular);
}

.result-artist {
  font-size: var(--font-size-small);
  color: var(--text-subdued);
  margin-top: 2px;
}

.show-more {
  padding: 12px 16px;
  border-top: 1px solid rgba(255, 255, 255, 0.06);
  text-align: center;
}

.show-more-btn {
  background: none;
  border: none;
  color: var(--text-subdued);
  font-size: var(--font-size-small);
  font-weight: var(--font-weight-bold);
  cursor: pointer;
}

.show-more-btn:hover {
  color: var(--text-base);
}

.no-results {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 40px 20px;
  gap: 8px;
}

.no-results p {
  margin: 0;
  font-size: var(--font-size-caption);
  color: var(--text-base);
  font-weight: var(--font-weight-bold);
}

.no-results span {
  font-size: var(--font-size-small);
  color: var(--text-subdued);
}

.slide-fade-enter-active,
.slide-fade-leave-active {
  transition: all 0.2s ease;
}

.slide-fade-enter-from,
.slide-fade-leave-to {
  opacity: 0;
  transform: translateY(-8px);
}

/* Scrollbar */
.search-dropdown::-webkit-scrollbar {
  width: 6px;
}

.search-dropdown::-webkit-scrollbar-track {
  background: transparent;
}

.search-dropdown::-webkit-scrollbar-thumb {
  background: rgba(255, 255, 255, 0.15);
  border-radius: 3px;
}

.search-dropdown::-webkit-scrollbar-thumb:hover {
  background: rgba(255, 255, 255, 0.3);
}
</style>
