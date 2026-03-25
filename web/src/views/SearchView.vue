<template>
  <div class="search-view">
    <!-- Header -->
    <div class="search-header safe-area-top">
      <n-button quaternary circle @click="$router.back()" class="back-btn">
        <template #icon>
          <n-icon :component="ArrowBack" :size="22" />
        </template>
      </n-button>
      <SearchBar
        @search="onSearch"
        @select="onSongSelect"
        @show-more="onShowMore"
        class="search-input"
      />
    </div>

    <!-- Search Results -->
    <div class="search-content">
      <!-- Loading -->
      <div v-if="isLoading" class="loading-state">
        <n-spin size="large" />
        <p>搜索中...</p>
      </div>

      <!-- Results List -->
      <div v-else-if="searchResults.length > 0" class="results-section">
        <div class="results-header">
          <span>搜索结果</span>
          <span class="results-count">共 {{ totalResults }} 首</span>
        </div>
        <div class="results-list">
          <div
            v-for="(song, index) in searchResults"
            :key="song.id"
            class="song-item"
            @click="playSong(song, index)"
          >
            <div class="song-index">{{ index + 1 }}</div>
            <img
              v-if="song.cover"
              :src="song.cover"
              class="song-cover"
              @error="onImageError"
            />
            <div v-else class="song-cover-placeholder">
              <n-icon :component="MusicalNote" :size="24" />
            </div>
            <div class="song-info">
              <div class="song-title">{{ song.title }}</div>
              <div class="song-meta">{{ song.artist }} · {{ song.album || '未知专辑' }}</div>
            </div>
            <div class="song-duration">{{ formatTime(song.duration) }}</div>
            <n-button quaternary circle size="tiny" @click.stop="addToPlaylist(song)">
              <template #icon>
                <n-icon :component="AddCircleOutline" :size="20" />
              </template>
            </n-button>
          </div>
        </div>

        <!-- Load More -->
        <div v-if="hasMore" class="load-more">
          <n-button text type="primary" :loading="isLoadingMore" @click="loadMore">
            加载更多
          </n-button>
        </div>
      </div>

      <!-- Empty State -->
      <div v-else-if="hasSearched && !isLoading" class="empty-state">
        <n-icon :component="SearchOutline" :size="64" />
        <p>未找到相关歌曲</p>
        <span>换个关键词试试</span>
      </div>

      <!-- Default State -->
      <div v-else class="default-state">
        <n-icon :component="SearchOutline" :size="64" />
        <p>搜索你喜欢的音乐</p>
        <span>支持歌曲名、歌手名搜索</span>
      </div>
    </div>

    <!-- Bottom Player Bar -->
    <div v-if="playerStore.currentSong" class="mini-player safe-area-bottom">
      <img
        v-if="playerStore.currentSong.cover"
        :src="playerStore.currentSong.cover"
        class="mini-cover"
        @error="onImageError"
      />
      <div v-else class="mini-cover-placeholder">
        <n-icon :component="MusicalNote" :size="18" />
      </div>
      <div class="mini-info">
        <div class="mini-title">{{ playerStore.currentSong.title }}</div>
        <div class="mini-artist">{{ playerStore.currentSong.artist }}</div>
      </div>
      <n-button quaternary circle @click="playerStore.togglePlay()">
        <template #icon>
          <n-icon :component="playerStore.isPlaying ? PauseCircle : PlayCircle" :size="32" />
        </template>
      </n-button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { NButton, NIcon, NSpin, useMessage } from 'naive-ui'
import {
  ArrowBack,
  MusicalNote,
  AddCircleOutline,
  SearchOutline,
  PlayCircle,
  PauseCircle,
} from '@vicons/ionicons5'
import SearchBar from '@/components/SearchBar.vue'
import { searchMusic } from '@/api/music'
import { usePlayerStore } from '@/stores/player'
import type { Song } from '@/types/music'

const route = useRoute()
const playerStore = usePlayerStore()
const message = useMessage()

const searchResults = ref<Song[]>([])
const totalResults = ref(0)
const isLoading = ref(false)
const isLoadingMore = ref(false)
const hasSearched = ref(false)
const hasMore = ref(false)
const currentQuery = ref('')
const currentPage = ref(1)
const pageSize = 20

onMounted(() => {
  const q = route.query.q as string
  if (q) {
    currentQuery.value = q
    performSearch(q)
  }
})

async function performSearch(query: string, page: number = 1) {
  if (!query.trim()) return

  if (page === 1) {
    isLoading.value = true
    searchResults.value = []
  } else {
    isLoadingMore.value = true
  }

  hasSearched.value = true
  currentQuery.value = query
  currentPage.value = page

  try {
    const response = await searchMusic(query, page, pageSize)
    if (page === 1) {
      searchResults.value = response.songs
    } else {
      searchResults.value.push(...response.songs)
    }
    totalResults.value = response.total
    hasMore.value = response.songs.length >= pageSize
  } catch (error: any) {
    console.error('Search failed:', error)
    message.error('搜索失败：' + (error.message || '网络错误'))
  } finally {
    isLoading.value = false
    isLoadingMore.value = false
  }
}

function onSearch(query: string) {
  performSearch(query)
}

async function onSongSelect(song: Song) {
  playerStore.addToPlaylist(song)
  await playerStore.playSong(song)
}

function onShowMore(query: string) {
  performSearch(query)
}

async function playSong(song: Song, index: number) {
  // 把当前搜索结果设为播放列表
  playerStore.setPlaylist(searchResults.value, index)
  await playerStore.playSong(song, index)
}

function addToPlaylist(song: Song) {
  playerStore.addToPlaylist(song)
  message.success(`已添加「${song.title}」到播放列表`)
}

function loadMore() {
  performSearch(currentQuery.value, currentPage.value + 1)
}

function formatTime(seconds: number): string {
  if (!seconds || isNaN(seconds)) return '--:--'
  const mins = Math.floor(seconds / 60)
  const secs = Math.floor(seconds % 60)
  return `${mins.toString().padStart(2, '0')}:${secs.toString().padStart(2, '0')}`
}

function onImageError(e: Event) {
  const target = e.target as HTMLImageElement
  target.style.display = 'none'
}
</script>

<style scoped>
.search-view {
  display: flex;
  flex-direction: column;
  height: 100vh;
  overflow: hidden;
}

.search-header {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 12px 16px;
  background: rgba(0, 0, 0, 0.2);
  backdrop-filter: blur(10px);
  z-index: 10;
}

.back-btn {
  flex-shrink: 0;
}

.search-input {
  flex: 1;
}

.search-content {
  flex: 1;
  overflow-y: auto;
  padding: 16px;
  padding-bottom: 80px;
}

.loading-state, .empty-state, .default-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 16px;
  padding: 60px 20px;
  color: rgba(255, 255, 255, 0.5);
  text-align: center;
}

.loading-state p, .empty-state p, .default-state p {
  font-size: 16px;
  color: rgba(255, 255, 255, 0.7);
}

.empty-state span, .default-state span {
  font-size: 13px;
}

.results-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
  font-size: 14px;
  color: rgba(255, 255, 255, 0.6);
}

.results-count {
  font-size: 12px;
  color: rgba(255, 255, 255, 0.4);
}

.results-list {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.song-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px;
  border-radius: 8px;
  cursor: pointer;
  transition: background 0.2s;
}

.song-item:hover {
  background: rgba(255, 255, 255, 0.08);
}

.song-item:active {
  background: rgba(255, 255, 255, 0.12);
}

.song-index {
  width: 24px;
  text-align: center;
  font-size: 14px;
  color: rgba(255, 255, 255, 0.4);
  flex-shrink: 0;
}

.song-cover {
  width: 48px;
  height: 48px;
  border-radius: 6px;
  object-fit: cover;
  flex-shrink: 0;
}

.song-cover-placeholder {
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

.song-info {
  flex: 1;
  min-width: 0;
}

.song-title {
  font-size: 15px;
  font-weight: 500;
  color: rgba(255, 255, 255, 0.95);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.song-meta {
  font-size: 12px;
  color: rgba(255, 255, 255, 0.45);
  margin-top: 4px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.song-duration {
  font-size: 12px;
  color: rgba(255, 255, 255, 0.4);
  font-variant-numeric: tabular-nums;
  flex-shrink: 0;
}

.load-more {
  text-align: center;
  padding: 20px;
}

.mini-player {
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 10px 16px;
  background: rgba(15, 12, 41, 0.98);
  backdrop-filter: blur(20px);
  border-top: 1px solid rgba(255, 255, 255, 0.1);
  z-index: 50;
}

.mini-cover {
  width: 40px;
  height: 40px;
  border-radius: 6px;
  object-fit: cover;
}

.mini-cover-placeholder {
  width: 40px;
  height: 40px;
  border-radius: 6px;
  background: rgba(255, 255, 255, 0.1);
  display: flex;
  align-items: center;
  justify-content: center;
  color: rgba(255, 255, 255, 0.3);
}

.mini-info {
  flex: 1;
  min-width: 0;
}

.mini-title {
  font-size: 14px;
  font-weight: 500;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.mini-artist {
  font-size: 12px;
  color: rgba(255, 255, 255, 0.5);
}
</style>
