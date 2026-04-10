<template>
  <div class="home-view">
    <!-- Hero Section -->
    <section class="hero-section">
      <div class="hero-greeting">
        <h1 class="greeting-text">{{ greetingText }}</h1>
      </div>
    </section>

    <!-- Quick Access Grid -->
    <section class="section">
      <h2 class="section-title">快速访问</h2>
      <div class="quick-grid">
        <div
          v-for="item in quickAccess"
          :key="item.name"
          class="quick-card"
        >
          <div class="quick-card-cover">
            <img v-if="item.cover" :src="item.cover" :alt="item.name" class="quick-card-img" />
            <div v-else class="quick-card-placeholder">
              <svg viewBox="0 0 24 24" width="24" height="24" fill="currentColor"><path d="M12 3v10.55c-.59-.34-1.27-.55-2-.55C7.79 13 6 14.79 6 17s1.79 4 4 4 4-1.79 4-4V7h4V3h-6z"/></svg>
            </div>
            <button class="quick-card-play">
              <svg viewBox="0 0 24 24" width="22" height="22" fill="#000"><path d="M8 5.14v13.72a1 1 0 001.5.86l11.04-6.86a1 1 0 000-1.72L9.5 4.28a1 1 0 00-1.5.86z"/></svg>
            </button>
          </div>
          <span class="quick-card-name">{{ item.name }}</span>
        </div>
      </div>
    </section>

    <!-- Recently Played -->
    <section class="section">
      <div class="section-header">
        <h2 class="section-title">最近播放</h2>
        <button class="show-all-btn">显示全部</button>
      </div>
      <div class="album-grid">
        <div
          v-for="song in recentSongs"
          :key="song.id"
          class="album-card"
          @click="playSong(song)"
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
    </section>

    <!-- Made For You -->
    <section class="section">
      <div class="section-header">
        <h2 class="section-title">为你推荐</h2>
        <button class="show-all-btn">显示全部</button>
      </div>
      <div class="album-grid">
        <div
          v-for="item in madeForYou"
          :key="item.name"
          class="album-card"
        >
          <div class="album-cover-wrapper">
            <img v-if="item.cover" :src="item.cover" :alt="item.name" class="album-cover" />
            <div v-else class="album-cover-placeholder">
              <svg viewBox="0 0 24 24" width="32" height="32" fill="currentColor"><path d="M12 3v10.55c-.59-.34-1.27-.55-2-.55C7.79 13 6 14.79 6 17s1.79 4 4 4 4-1.79 4-4V7h4V3h-6z"/></svg>
            </div>
            <button class="album-play-btn">
              <svg viewBox="0 0 24 24" width="28" height="28" fill="#000"><path d="M8 5.14v13.72a1 1 0 001.5.86l11.04-6.86a1 1 0 000-1.72L9.5 4.28a1 1 0 00-1.5.86z"/></svg>
            </button>
          </div>
          <div class="album-title">{{ item.name }}</div>
          <div class="album-subtitle">{{ item.description }}</div>
        </div>
      </div>
    </section>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { usePlayerStore } from '../stores/player'
import type { Song } from '../types/music'

const playerStore = usePlayerStore()

// Greeting based on time
const greetingText = computed(() => {
  const hour = new Date().getHours()
  if (hour < 12) return '早上好'
  if (hour < 18) return '下午好'
  return '晚上好'
})

// Recent songs from playlist
const recentSongs = computed(() => playerStore.playlist.slice(0, 6))

// Quick access cards
const quickAccess = computed(() => {
  const list = playerStore.playlist.slice(0, 6)
  if (list.length === 0) {
    return [
      { name: '我喜欢的音乐', cover: '' },
      { name: '最近播放', cover: '' },
      { name: '每日推荐', cover: '' },
      { name: '流行热歌', cover: '' },
      { name: '古典精选', cover: '' },
      { name: '轻音乐', cover: '' },
    ]
  }
  return list.map(s => ({ name: s.title, cover: s.cover }))
})

// Made for you
const madeForYou = [
  { name: '每日推荐', description: '根据你的口味推荐', cover: '' },
  { name: '流行热歌', description: '最热的流行音乐', cover: '' },
  { name: '华语新歌', description: '最新华语单曲', cover: '' },
  { name: '欧美金曲', description: '经典欧美热门', cover: '' },
  { name: '古典精选', description: '优雅古典乐章', cover: '' },
  { name: '电子节拍', description: '电子音乐精选', cover: '' },
]

async function playSong(song: Song) {
  await playerStore.playSong(song)
}
</script>

<style scoped>
.home-view {
  padding: 24px 32px 120px;
  min-height: 100%;
}

/* Hero / Greeting */
.hero-section {
  margin-bottom: 24px;
}

.greeting-text {
  font-size: 32px;
  font-weight: var(--font-weight-bold);
  color: var(--text-base);
  margin: 0;
  letter-spacing: -0.02em;
}

/* Section */
.section {
  margin-bottom: 32px;
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
}

.section-title {
  font-size: var(--font-size-title);
  font-weight: var(--font-weight-bold);
  color: var(--text-base);
  margin: 0;
  letter-spacing: -0.02em;
}

.show-all-btn {
  background: none;
  border: none;
  color: var(--text-subdued);
  font-size: var(--font-size-small);
  font-weight: var(--font-weight-bold);
  cursor: pointer;
  text-transform: uppercase;
  letter-spacing: 1px;
}

.show-all-btn:hover {
  color: var(--text-base);
}

/* Quick Access Grid — Spotify home cards */
.quick-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(180px, 1fr));
  gap: 8px;
}

.quick-card {
  display: flex;
  align-items: center;
  gap: 12px;
  background: rgba(255, 255, 255, 0.06);
  border-radius: var(--radius-md);
  overflow: hidden;
  cursor: pointer;
  transition: background 0.2s;
  height: 64px;
}

.quick-card:hover {
  background: rgba(255, 255, 255, 0.12);
}

.quick-card-cover {
  width: 64px;
  height: 64px;
  flex-shrink: 0;
  position: relative;
}

.quick-card-img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.quick-card-placeholder {
  width: 100%;
  height: 100%;
  background: var(--bg-highlight);
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--text-subdued);
}

.quick-card-play {
  position: absolute;
  right: -4px;
  bottom: -4px;
  width: 32px;
  height: 32px;
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

.quick-card:hover .quick-card-play {
  opacity: 1;
  transform: translateY(0);
}

.quick-card-play:hover {
  transform: scale(1.06);
  background: var(--spotify-green-bright);
}

.quick-card-name {
  font-size: var(--font-size-caption);
  font-weight: var(--font-weight-bold);
  color: var(--text-base);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  padding-right: 12px;
}

/* Album Grid — Spotify card grid */
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
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
  line-height: 1.4;
}

/* Responsive */
@media (max-width: 768px) {
  .home-view {
    padding: 16px;
  }

  .greeting-text {
    font-size: 24px;
  }

  .album-grid {
    grid-template-columns: repeat(auto-fill, minmax(140px, 1fr));
    gap: 16px;
  }

  .quick-grid {
    grid-template-columns: repeat(auto-fill, minmax(150px, 1fr));
  }
}
</style>
