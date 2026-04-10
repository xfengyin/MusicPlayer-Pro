<template>
  <div class="playlist-view">
    <!-- Playlist Header -->
    <div class="playlist-hero">
      <div class="playlist-hero-cover">
        <svg viewBox="0 0 24 24" width="64" height="64" fill="currentColor" style="color: var(--text-subdued);"><path d="M3 4h18v2H3V4zm0 7h12v2H3v-2zm0 7h18v2H3v-2zm14-3l6-4v8l-6-4z"/></svg>
      </div>
      <div class="playlist-hero-info">
        <span class="playlist-type">播放列表</span>
        <h1 class="playlist-hero-title">我的播放列表</h1>
        <p class="playlist-hero-meta">{{ songs.length }} 首歌曲</p>
      </div>
    </div>

    <!-- Action Bar -->
    <div class="action-bar">
      <button class="play-all-btn" @click="playAll">
        <svg viewBox="0 0 24 24" width="24" height="24" fill="#000"><path d="M8 5.14v13.72a1 1 0 001.5.86l11.04-6.86a1 1 0 000-1.72L9.5 4.28a1 1 0 00-1.5.86z"/></svg>
      </button>
    </div>

    <!-- Song List Table -->
    <div class="song-table">
      <div class="song-table-header">
        <span class="col-num">#</span>
        <span class="col-title">标题</span>
        <span class="col-album">专辑</span>
        <span class="col-duration">
          <svg viewBox="0 0 24 24" width="16" height="16" fill="var(--text-subdued)"><path d="M12 2C6.477 2 2 6.477 2 12s4.477 10 10 10 10-4.477 10-10S17.523 2 12 2zm0 18a8 8 0 110-16 8 8 0 010 16zm-.5-13a.5.5 0 01.5.5v4.5h3.5a.5.5 0 010 1H11.5a.5.5 0 01-.5-.5V7.5a.5.5 0 01.5-.5z"/></svg>
        </span>
      </div>

      <div class="song-table-body">
        <div
          v-for="(song, index) in songs"
          :key="song.id"
          :class="['song-row', { 'is-playing': currentSong?.id === song.id }]"
          @click="playSong(song, index)"
        >
          <span class="col-num">
            <span v-if="currentSong?.id !== song.id" class="row-index">{{ index + 1 }}</span>
            <div v-else class="playing-indicator">
              <span class="bar"></span>
              <span class="bar"></span>
              <span class="bar"></span>
            </div>
          </span>
          <div class="col-title">
            <div class="song-cover-mini">
              <img v-if="song.cover" :src="song.cover" :alt="song.title" class="cover-img" />
              <div v-else class="cover-placeholder-mini">
                <svg viewBox="0 0 24 24" width="14" height="14" fill="currentColor"><path d="M12 3v10.55c-.59-.34-1.27-.55-2-.55C7.79 13 6 14.79 6 17s1.79 4 4 4 4-1.79 4-4V7h4V3h-6z"/></svg>
              </div>
            </div>
            <div class="song-title-info">
              <div class="song-name" :class="{ 'active-name': currentSong?.id === song.id }">{{ song.title }}</div>
              <div class="song-artist-name">{{ song.artist }}</div>
            </div>
          </div>
          <span class="col-album">{{ song.album }}</span>
          <span class="col-duration">{{ formatTime(song.duration) }}</span>
        </div>
      </div>

      <div v-if="songs.length === 0" class="empty-state">
        <p>播放列表为空</p>
        <span>去搜索添加你喜欢的歌曲吧</span>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { usePlayerStore } from '../stores/player'
import type { Song } from '../types/music'

const playerStore = usePlayerStore()
const songs = computed(() => playerStore.playlist)
const currentSong = computed(() => playerStore.currentSong)

function formatTime(seconds: number): string {
  if (!seconds || isNaN(seconds)) return '-:--'
  const mins = Math.floor(seconds / 60)
  const secs = Math.floor(seconds % 60)
  return `${mins}:${secs.toString().padStart(2, '0')}`
}

async function playSong(song: Song, index: number) {
  await playerStore.playSong(song, index)
}

async function playAll() {
  if (songs.value.length > 0) {
    await playerStore.playSong(songs.value[0], 0)
  }
}
</script>

<style scoped>
.playlist-view {
  padding: 0 0 120px;
  min-height: 100%;
}

/* Hero Section */
.playlist-hero {
  display: flex;
  align-items: flex-end;
  gap: 24px;
  padding: 40px 32px 24px;
  background: linear-gradient(180deg, #3d2a54 0%, var(--bg-base) 100%);
}

.playlist-hero-cover {
  width: 232px;
  height: 232px;
  border-radius: var(--radius-md);
  background: var(--bg-highlight);
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  box-shadow: 0 4px 60px rgba(0, 0, 0, 0.5);
}

.playlist-hero-info {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.playlist-type {
  font-size: var(--font-size-small);
  font-weight: var(--font-weight-bold);
  color: var(--text-base);
  text-transform: uppercase;
  letter-spacing: 1px;
}

.playlist-hero-title {
  font-size: 48px;
  font-weight: var(--font-weight-bold);
  color: var(--text-base);
  margin: 0;
  letter-spacing: -0.04em;
  line-height: 1.1;
}

.playlist-hero-meta {
  font-size: var(--font-size-small);
  color: var(--text-subdued);
  margin: 0;
}

/* Action Bar */
.action-bar {
  padding: 16px 32px;
  display: flex;
  align-items: center;
  gap: 24px;
  background: linear-gradient(180deg, rgba(61, 42, 84, 0.4) 0%, transparent 100%);
}

.play-all-btn {
  width: 56px;
  height: 56px;
  border-radius: var(--radius-circle);
  background: var(--spotify-green);
  border: none;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.1s;
  box-shadow: 0 8px 16px rgba(0, 0, 0, 0.3);
}

.play-all-btn:hover {
  transform: scale(1.06);
  background: var(--spotify-green-bright);
}

.play-all-btn:active {
  transform: scale(0.98);
}

/* Song Table */
.song-table {
  padding: 0 32px;
}

.song-table-header {
  display: flex;
  align-items: center;
  padding: 8px 12px;
  border-bottom: 1px solid rgba(255, 255, 255, 0.06);
  color: var(--text-subdued);
  font-size: var(--font-size-small);
  font-weight: var(--font-weight-bold);
  text-transform: uppercase;
  letter-spacing: 1px;
}

.col-num {
  width: 32px;
  text-align: center;
  flex-shrink: 0;
}

.col-title {
  flex: 1;
  min-width: 0;
}

.col-album {
  width: 200px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  flex-shrink: 0;
}

.col-duration {
  width: 60px;
  text-align: right;
  flex-shrink: 0;
}

/* Song Row */
.song-table-body {
  display: flex;
  flex-direction: column;
}

.song-row {
  display: flex;
  align-items: center;
  padding: 8px 12px;
  border-radius: var(--radius-sm);
  cursor: pointer;
  transition: background 0.2s;
}

.song-row:hover {
  background: rgba(255, 255, 255, 0.08);
}

.song-row:hover .col-album {
  color: var(--text-base);
}

.song-row.is-playing .song-name {
  color: var(--spotify-green);
}

/* Playing indicator */
.playing-indicator {
  display: flex;
  align-items: flex-end;
  gap: 2px;
  height: 14px;
}

.playing-indicator .bar {
  width: 3px;
  background: var(--spotify-green);
  border-radius: 1px;
  animation: bar-bounce 0.8s ease-in-out infinite;
}

.playing-indicator .bar:nth-child(1) { height: 6px; animation-delay: 0s; }
.playing-indicator .bar:nth-child(2) { height: 10px; animation-delay: 0.2s; }
.playing-indicator .bar:nth-child(3) { height: 4px; animation-delay: 0.4s; }

@keyframes bar-bounce {
  0%, 100% { transform: scaleY(1); }
  50% { transform: scaleY(0.4); }
}

/* Song cover mini */
.song-cover-mini {
  width: 40px;
  height: 40px;
  border-radius: var(--radius-sm);
  overflow: hidden;
  flex-shrink: 0;
  margin-right: 12px;
}

.cover-img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.cover-placeholder-mini {
  width: 100%;
  height: 100%;
  background: var(--bg-highlight);
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--text-subdued);
}

.col-title {
  display: flex;
  align-items: center;
}

.song-title-info {
  min-width: 0;
}

.song-name {
  font-size: var(--font-size-caption);
  color: var(--text-base);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  font-weight: var(--font-weight-regular);
}

.song-name.active-name {
  color: var(--spotify-green);
}

.song-artist-name {
  font-size: var(--font-size-small);
  color: var(--text-subdued);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.song-row .col-album {
  color: var(--text-subdued);
  font-size: var(--font-size-caption);
  transition: color 0.2s;
}

.song-row .col-duration {
  font-size: var(--font-size-small);
  color: var(--text-subdued);
  font-variant-numeric: tabular-nums;
}

/* Empty */
.empty-state {
  padding: 60px 0;
  text-align: center;
}

.empty-state p {
  font-size: var(--font-size-body);
  font-weight: var(--font-weight-bold);
  color: var(--text-base);
  margin-bottom: 8px;
}

.empty-state span {
  font-size: var(--font-size-caption);
  color: var(--text-subdued);
}

@media (max-width: 768px) {
  .playlist-hero {
    flex-direction: column;
    align-items: center;
    padding: 24px;
  }

  .playlist-hero-cover {
    width: 160px;
    height: 160px;
  }

  .playlist-hero-title {
    font-size: 28px;
    text-align: center;
  }

  .col-album {
    display: none;
  }

  .song-table {
    padding: 0 16px;
  }
}
</style>
