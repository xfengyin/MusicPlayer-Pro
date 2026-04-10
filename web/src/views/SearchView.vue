<template>
  <div class="search-view">
    <!-- Search Input -->
    <div class="search-header">
      <SearchBar
        @search="onSearch"
        @select="onSelect"
        @show-more="onShowMore"
      />
    </div>

    <!-- Browse Categories (when no search) -->
    <div v-if="!hasSearched" class="browse-section">
      <h2 class="section-title">浏览全部</h2>
      <div class="category-grid">
        <div
          v-for="cat in categories"
          :key="cat.name"
          class="category-card"
          :style="{ backgroundColor: cat.color }"
        >
          <span class="category-name">{{ cat.name }}</span>
          <div class="category-icon">
            <svg viewBox="0 0 24 24" width="24" height="24" fill="currentColor" style="opacity: 0.6;"><path d="M12 3v10.55c-.59-.34-1.27-.55-2-.55C7.79 13 6 14.79 6 17s1.79 4 4 4 4-1.79 4-4V7h4V3h-6z"/></svg>
          </div>
        </div>
      </div>
    </div>

    <!-- Search Results -->
    <div v-else class="results-section">
      <h2 class="section-title">搜索结果</h2>
      <div v-if="results.length > 0" class="album-grid">
        <div
          v-for="song in results"
          :key="song.id"
          class="album-card"
          @click="playResult(song)"
        >
          <div class="album-cover-wrapper">
            <img v-if="song.cover" :src="song.cover" :alt="song.title" class="album-cover" />
            <div v-else class="album-cover-placeholder">
              <svg viewBox="0 0 24 24" width="32" height="32" fill="currentColor"><path d="M12 3v10.55c-.59-.34-1.27-.55-2-.55C7.79 13 6 14.79 6 17s1.79 4 4 4 4-1.79 4-4V7h4V3h-6z"/></svg>
            </div>
            <button class="album-play-btn">
              <svg viewBox="0 0 24 24" width="28" height="28" fill="#000"><path d="M8 5.14v13.72a1 1 0 001.5.86l11.04-6.86a1 1 0 000-1.72L9.5 4.28a1 1 0 00-1.5.86z"/></svg>
            </button>
          </div>
          <div class="album-title">{{ song.title }}</div>
          <div class="album-subtitle">{{ song.artist }}</div>
        </div>
      </div>
      <div v-else class="no-results">
        <p>未找到相关结果</p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import SearchBar from '../components/SearchBar.vue'
import { searchMusic } from '../api/music'
import { usePlayerStore } from '../stores/player'
import type { Song } from '../types/music'

const playerStore = usePlayerStore()

const hasSearched = ref(false)
const results = ref<Song[]>([])

const categories = [
  { name: '流行', color: '#8c67ab' },
  { name: '摇滚', color: '#e8115b' },
  { name: '电子', color: '#1e3264' },
  { name: '古典', color: '#7d4b32' },
  { name: '爵士', color: '#477d95' },
  { name: '民谣', color: '#ba5d07' },
  { name: '嘻哈', color: '#148a08' },
  { name: 'R&B', color: '#503750' },
  { name: '动漫', color: '#e61e32' },
  { name: '影视', color: '#2d46b9' },
  { name: '轻音乐', color: '#537aa5' },
  { name: '中国风', color: '#c39687' },
]

async function onSearch(query: string) {
  hasSearched.value = true
  try {
    const response = await searchMusic(query, 1, 20)
    results.value = response.songs
  } catch (e) {
    results.value = []
  }
}

async function onSelect(song: Song) {
  await playerStore.playSong(song)
}

function onShowMore(query: string) {
  onSearch(query)
}

async function playResult(song: Song) {
  playerStore.addToPlaylist(song)
  await playerStore.playSong(song)
}
</script>

<style scoped>
.search-view {
  padding: 24px 32px 120px;
  min-height: 100%;
}

.search-header {
  margin-bottom: 24px;
}

/* Section Title */
.section-title {
  font-size: var(--font-size-title);
  font-weight: var(--font-weight-bold);
  color: var(--text-base);
  margin: 0 0 16px;
  letter-spacing: -0.02em;
}

/* Category Grid */
.category-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(180px, 1fr));
  gap: 16px;
}

.category-card {
  position: relative;
  border-radius: var(--radius-lg);
  padding: 16px;
  height: 180px;
  overflow: hidden;
  cursor: pointer;
}

.category-name {
  font-size: var(--font-size-heading);
  font-weight: var(--font-weight-bold);
  color: var(--text-base);
}

.category-icon {
  position: absolute;
  right: -8px;
  bottom: -8px;
  transform: rotate(25deg);
  color: var(--text-base);
}

/* Album Grid */
.album-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(180px, 1fr));
  gap: 24px;
}

.album-card {
  cursor: pointer;
  padding: 16px;
  border-radius: var(--radius-lg);
  background: var(--bg-card);
  transition: background 0.3s;
}

.album-card:hover {
  background: var(--bg-card-hover);
}

.album-cover-wrapper {
  position: relative;
  width: 100%;
  aspect-ratio: 1;
  margin-bottom: 16px;
  border-radius: var(--radius-md);
  overflow: hidden;
}

.album-cover {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.album-cover-placeholder {
  width: 100%;
  height: 100%;
  background: var(--bg-highlight);
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--text-subdued);
}

.album-play-btn {
  position: absolute;
  right: 8px;
  bottom: 8px;
  width: 48px;
  height: 48px;
  border-radius: var(--radius-circle);
  background: var(--spotify-green);
  border: none;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  opacity: 0;
  transform: translateY(8px);
  transition: all 0.3s;
  box-shadow: 0 8px 16px rgba(0, 0, 0, 0.3);
}

.album-card:hover .album-play-btn {
  opacity: 1;
  transform: translateY(0);
}

.album-play-btn:hover {
  transform: scale(1.06);
  background: var(--spotify-green-bright);
}

.album-title {
  font-size: var(--font-size-caption);
  font-weight: var(--font-weight-bold);
  color: var(--text-base);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  margin-bottom: 4px;
}

.album-subtitle {
  font-size: var(--font-size-small);
  color: var(--text-subdued);
}

.no-results {
  padding: 60px 0;
  text-align: center;
  color: var(--text-subdued);
}

.no-results p {
  font-size: var(--font-size-body);
  font-weight: var(--font-weight-bold);
  color: var(--text-base);
}

@media (max-width: 768px) {
  .search-view {
    padding: 16px;
  }

  .category-grid,
  .album-grid {
    grid-template-columns: repeat(auto-fill, minmax(140px, 1fr));
  }
}
</style>
