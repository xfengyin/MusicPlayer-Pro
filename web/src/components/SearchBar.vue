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
    >
      <template #prefix>
        <n-icon :component="Search" />
      </template>
    </n-input>

    <!-- 加载状态 -->
    <div v-if="isLoading" class="loading-indicator">
      <n-spin size="small" />
      <span>搜索中...</span>
    </div>

    <!-- 搜索历史下拉 -->
    <transition name="slide-fade">
      <div
        v-if="showHistory && searchHistory.length > 0 && isFocused && !hasQuery"
        class="search-dropdown"
      >
        <div class="dropdown-header">
          <span>搜索历史</span>
          <n-button text type="primary" size="tiny" @click="clearHistory">
            清空
          </n-button>
        </div>
        <div class="history-list">
          <div
            v-for="(item, index) in searchHistory"
            :key="index"
            class="history-item"
            @click="selectHistory(item)"
          >
            <n-icon :component="Time" />
            <span>{{ item }}</span>
            <n-icon
              :component="Close"
              class="remove-btn"
              @click.stop="removeHistory(index)"
            />
          </div>
        </div>
      </div>
    </transition>

    <!-- 搜索结果预览 -->
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
            <img
              v-if="song.cover"
              :src="song.cover"
              :alt="song.title"
              class="result-cover"
              @error="onImageError"
            />
            <div v-else class="result-cover-placeholder">
              <n-icon :component="MusicalNote" :size="24" />
            </div>
            <div class="result-info">
              <div class="result-title">{{ song.title }}</div>
              <div class="result-artist">{{ song.artist }}</div>
            </div>
            <n-icon :component="Play" class="play-icon" />
          </div>
        </div>
        <div v-if="hasMoreResults" class="show-more">
          <n-button text type="primary" @click="onShowMore">查看全部结果</n-button>
        </div>
      </div>
    </transition>

    <!-- 无结果提示 -->
    <transition name="slide-fade">
      <div
        v-if="showResults && searchResults.length === 0 && !isLoading && hasQuery"
        class="search-dropdown"
      >
        <div class="no-results">
          <n-icon :component="SearchOff" :size="48" />
          <p>未找到相关歌曲</p>
          <span>试试其他关键词</span>
        </div>
      </div>
    </transition>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch, onMounted } from 'vue'
import { NInput, NIcon, NButton, NSpin } from 'naive-ui'
import { Search, Time, Close, Play, MusicalNote, SearchOff } from '@vicons/ionicons5'
import { searchMusic } from '../api/music'
import type { Song } from '../types/music'
import { usePlayerStore } from '../stores/player'

// Props
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

// Emits
const emit = defineEmits<{
  search: [query: string]
  select: [song: Song]
  'history-clear': []
  'show-more': [query: string]
}>()

// 状态
const inputRef = ref()
const searchQuery = ref('')
const isFocused = ref(false)
const showHistory = ref(false)
const showResults = ref(false)
const isLoading = ref(false)
const searchHistory = ref<string[]>([])
const searchResults = ref<Song[]>([])
const hasMoreResults = ref(false)

// 播放器 store
const playerStore = usePlayerStore()

// 防抖定时器
let debounceTimer: ReturnType<typeof setTimeout> | null = null

// 计算属性
const hasQuery = computed(() => searchQuery.value.trim().length > 0)

// 监听搜索输入（带防抖）
watch(searchQuery, (newQuery) => {
  if (debounceTimer) {
    clearTimeout(debounceTimer)
  }

  if (newQuery.trim().length > 0) {
    // 延迟搜索
    debounceTimer = setTimeout(() => {
      performSearch(newQuery)
    }, props.debounceMs)
  } else {
    // 清空时显示历史记录
    searchResults.value = []
    showResults.value = false
    if (searchHistory.value.length > 0 && isFocused.value) {
      showHistory.value = true
    }
  }
})

// 方法
function onFocus() {
  isFocused.value = true
  if (searchHistory.value.length > 0 && !hasQuery.value) {
    showHistory.value = true
  }
}

function onBlur() {
  // 延迟关闭，允许点击下拉内容
  setTimeout(() => {
    isFocused.value = false
    showHistory.value = false
    // showResults 保持显示，直到下次搜索或失焦
  }, 200)
}

function onSearch() {
  const query = searchQuery.value.trim()
  if (query) {
    addToHistory(query)
    emit('search', query)
    showHistory.value = false
    // 执行搜索
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
  // 添加到播放列表
  playerStore.addToPlaylist(song)
  // 播放歌曲
  await playerStore.playSong(song)
  // 触发事件
  emit('select', song)
  // 清空搜索框
  searchQuery.value = ''
  showResults.value = false
}

function onShowMore() {
  const query = searchQuery.value.trim()
  if (query) {
    emit('show-more', query)
  }
}

function loadHistory() {
  try {
    const saved = localStorage.getItem('searchHistory')
    if (saved) {
      searchHistory.value = JSON.parse(saved)
    }
  } catch (e) {
    console.error('Failed to load search history:', e)
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
  // 移除重复项
  const index = searchHistory.value.indexOf(query)
  if (index > -1) {
    searchHistory.value.splice(index, 1)
  }
  // 添加到开头
  searchHistory.value.unshift(query)
  // 限制数量
  if (searchHistory.value.length > props.maxHistory) {
    searchHistory.value.pop()
  }
  saveHistory()
}

function onImageError(e: Event) {
  const target = e.target as HTMLImageElement
  target.style.display = 'none'
}

// 生命周期
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
  color: rgba(255, 255, 255, 0.6);
  font-size: 13px;
}

.search-dropdown {
  position: absolute;
  top: calc(100% + 8px);
  left: 0;
  right: 0;
  background: rgba(30, 30, 40, 0.98);
  border-radius: 8px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.4);
  backdrop-filter: blur(10px);
  z-index: 100;
  max-height: 400px;
  overflow-y: auto;
  border: 1px solid rgba(255, 255, 255, 0.1);
}

.dropdown-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 16px;
  border-bottom: 1px solid rgba(255, 255, 255, 0.1);
  font-size: 12px;
  color: rgba(255, 255, 255, 0.5);
}

.history-list,
.results-list {
  padding: 8px 0;
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
  color: rgba(255, 255, 255, 0.8);
}

.history-item .n-icon:first-child {
  color: rgba(255, 255, 255, 0.4);
}

.remove-btn {
  margin-left: auto;
  opacity: 0;
  transition: opacity 0.2s;
}

.history-item:hover .remove-btn {
  opacity: 1;
}

.result-cover {
  width: 48px;
  height: 48px;
  border-radius: 6px;
  object-fit: cover;
  flex-shrink: 0;
}

.result-cover-placeholder {
  width: 48px;
  height: 48px;
  border-radius: 6px;
  background: rgba(255, 255, 255, 0.1);
  display: flex;
  align-items: center;
  justify-content: center;
  color: rgba(255, 255, 255, 0.3);
  flex-shrink: 0;
}

.result-info {
  flex: 1;
  min-width: 0;
}

.result-title {
  font-size: 14px;
  color: rgba(255, 255, 255, 0.95);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  font-weight: 500;
}

.result-artist {
  font-size: 12px;
  color: rgba(255, 255, 255, 0.5);
  margin-top: 2px;
}

.play-icon {
  opacity: 0;
  transition: opacity 0.2s;
  color: rgba(255, 255, 255, 0.6);
}

.result-item:hover .play-icon {
  opacity: 1;
}

.show-more {
  padding: 12px 16px;
  border-top: 1px solid rgba(255, 255, 255, 0.1);
  text-align: center;
}

.no-results {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 40px 20px;
  color: rgba(255, 255, 255, 0.4);
  gap: 12px;
}

.no-results p {
  margin: 0;
  font-size: 14px;
  color: rgba(255, 255, 255, 0.6);
}

.no-results span {
  font-size: 12px;
}

.slide-fade-enter-active,
.slide-fade-leave-active {
  transition: all 0.2s ease;
}

.slide-fade-enter-from,
.slide-fade-leave-to {
  opacity: 0;
  transform: translateY(-10px);
}

/* 滚动条样式 */
.search-dropdown::-webkit-scrollbar {
  width: 6px;
}

.search-dropdown::-webkit-scrollbar-track {
  background: rgba(255, 255, 255, 0.05);
  border-radius: 3px;
}

.search-dropdown::-webkit-scrollbar-thumb {
  background: rgba(255, 255, 255, 0.2);
  border-radius: 3px;
}

.search-dropdown::-webkit-scrollbar-thumb:hover {
  background: rgba(255, 255, 255, 0.3);
}
</style>
