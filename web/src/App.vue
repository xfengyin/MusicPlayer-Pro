<template>
  <n-config-provider :theme="darkTheme" :theme-overrides="themeOverrides">
    <n-message-provider>
      <n-dialog-provider>
        <div class="app-layout">
          <!-- Sidebar -->
          <aside class="sidebar">
            <div class="sidebar-logo">
              <svg viewBox="0 0 24 24" width="28" height="28" fill="currentColor">
                <path d="M12 0C5.4 0 0 5.4 0 12s5.4 12 12 12 12-5.4 12-12S18.66 0 12 0zm5.521 17.34c-.24.359-.66.48-1.021.24-2.82-1.74-6.36-2.101-10.561-1.141-.418.122-.779-.179-.899-.539-.12-.421.18-.78.54-.9 4.56-1.021 8.52-.6 11.64 1.32.42.18.479.659.301 1.02zm1.44-3.3c-.301.42-.841.6-1.262.3-3.239-1.98-8.159-2.58-11.939-1.38-.479.12-1.02-.12-1.14-.6-.12-.48.12-1.021.6-1.141C9.6 9.9 15 10.561 18.72 12.84c.361.181.54.78.241 1.2zm.12-3.36C15.24 8.4 8.82 8.16 5.16 9.301c-.6.179-1.2-.181-1.38-.721-.18-.601.18-1.2.72-1.381 4.26-1.26 11.28-1.02 15.721 1.621.539.3.719 1.02.419 1.56-.299.421-1.02.599-1.559.3z"/>
              </svg>
              <span class="logo-text">Music Player Pro</span>
            </div>

            <nav class="sidebar-nav">
              <router-link to="/" class="nav-item" :class="{ active: $route.path === '/' }">
                <svg viewBox="0 0 24 24" width="24" height="24" fill="currentColor"><path d="M12.5 3.247a1 1 0 00-1 0L4 7.577V20h4.5v-6a1 1 0 011-1h5a1 1 0 011 1v6H20V7.577l-7.5-4.33zm-2-1.732a3 3 0 013 0l7.5 4.33a2 2 0 011 1.732V21a1 1 0 01-1 1h-6.5a1 1 0 01-1-1v-6h-3v6a1 1 0 01-1 1H3a1 1 0 01-1-1V7.577a2 2 0 011-1.732l7.5-4.33z"/></svg>
                <span>首页</span>
              </router-link>
              <router-link to="/search" class="nav-item" :class="{ active: $route.path === '/search' }">
                <svg viewBox="0 0 24 24" width="24" height="24" fill="currentColor"><path d="M10.533 1.279c-5.18 0-9.407 4.14-9.407 9.279s4.226 9.279 9.407 9.279c2.234 0 4.29-.77 5.907-2.058l4.353 4.353a1 1 0 101.414-1.414l-4.344-4.344a9.157 9.157 0 002.077-5.816c0-5.14-4.226-9.28-9.407-9.28zm-7.407 9.279c0-4.006 3.302-7.28 7.407-7.28s7.407 3.274 7.407 7.28-3.302 7.279-7.407 7.279-7.407-3.273-7.407-7.28z"/></svg>
                <span>搜索</span>
              </router-link>
              <router-link to="/playlist" class="nav-item" :class="{ active: $route.path === '/playlist' }">
                <svg viewBox="0 0 24 24" width="24" height="24" fill="currentColor"><path d="M3 22a1 1 0 01-1-1V3a1 1 0 012 0v18a1 1 0 01-1 1zM15.5 2.134A1 1 0 0014 3v18a1 1 0 001.5.866l11-9a1 1 0 000-1.732l-11-9zM16 4.732L23.744 12 16 19.268V4.732zM7 22a1 1 0 01-1-1V3a1 1 0 012 0v18a1 1 0 01-1 1z"/></svg>
                <span>播放列表</span>
              </router-link>
              <router-link to="/settings" class="nav-item" :class="{ active: $route.path === '/settings' }">
                <svg viewBox="0 0 24 24" width="24" height="24" fill="currentColor"><path d="M12 1a3 3 0 00-3 3v1.07A7.005 7.005 0 005.07 9H4a3 3 0 000 6h1.07A7.005 7.005 0 009 18.93V20a3 3 0 006 0v-1.07A7.005 7.005 0 0018.93 16H20a3 3 0 000-6h-1.07A7.005 7.005 0 0015 5.07V4a3 3 0 00-3-3zm0 2a1 1 0 011 1v2.423l.71.265A5.002 5.002 0 0117.312 10.29l.265.71H20a1 1 0 010 2h-2.423l-.265.71a5.002 5.002 0 01-3.602 3.602l-.71.265V20a1 1 0 01-2 0v-2.423l-.71-.265A5.002 5.002 0 016.688 13.71l-.265-.71H4a1 1 0 010-2h2.423l.265-.71A5.002 5.002 0 019.29 6.688l.71-.265V4a1 1 0 011-1zm0 6a3 3 0 100 6 3 3 0 000-6zm0 2a1 1 0 110 2 1 1 0 010-2z"/></svg>
                <span>设置</span>
              </router-link>
            </nav>

            <div class="sidebar-divider"></div>

            <!-- Playlist shortcuts -->
            <div class="sidebar-section">
              <div class="sidebar-section-title">我的歌单</div>
              <div class="sidebar-playlists">
                <div class="playlist-link" v-for="item in quickPlaylists" :key="item.name">
                  <span>{{ item.name }}</span>
                </div>
              </div>
            </div>
          </aside>

          <!-- Main Content -->
          <main class="main-content">
            <div class="content-scroll">
              <router-view />
            </div>
          </main>

          <!-- Fixed Bottom Player -->
          <footer class="player-bar">
            <!-- Now Playing Info -->
            <div class="player-song-info">
              <div v-if="currentSong" class="now-playing-cover">
                <img
                  v-if="currentSong.cover"
                  :src="currentSong.cover"
                  :alt="currentSong.title"
                  class="cover-img"
                />
                <div v-else class="cover-placeholder">
                  <svg viewBox="0 0 24 24" width="20" height="20" fill="currentColor"><path d="M12 3v10.55c-.59-.34-1.27-.55-2-.55C7.79 13 6 14.79 6 17s1.79 4 4 4 4-1.79 4-4V7h4V3h-6z"/></svg>
                </div>
              </div>
              <div v-if="currentSong" class="now-playing-text">
                <div class="now-playing-title">{{ currentSong.title }}</div>
                <div class="now-playing-artist">{{ currentSong.artist }}</div>
              </div>
              <div v-if="!currentSong" class="now-playing-empty">
                <span style="color: var(--text-subdued); font-size: var(--font-size-caption);">未在播放</span>
              </div>
            </div>

            <!-- Player Controls -->
            <div class="player-center">
              <PlayerControls
                :showVolume="false"
                @toggle-playlist="togglePlaylistPanel"
              />
            </div>

            <!-- Volume & Extra -->
            <div class="player-right">
              <div class="volume-control">
                <n-icon :component="volumeIcon" :size="18" @click="toggleMute" style="cursor: pointer; color: var(--text-subdued);" />
                <n-slider
                  v-model:value="volumeValue"
                  :min="0"
                  :max="1"
                  :step="0.01"
                  :tooltip="false"
                  @update:value="onVolumeChange"
                  class="volume-slider"
                  style="width: 100px;"
                />
              </div>
            </div>
          </footer>

          <!-- Playlist Side Panel -->
          <transition name="slide-right">
            <aside v-if="showPlaylistPanel" class="playlist-panel">
              <div class="playlist-panel-header">
                <h3>播放列表</h3>
                <n-button quaternary circle size="small" @click="showPlaylistPanel = false">
                  <template #icon>
                    <n-icon :component="CloseIcon" />
                  </template>
                </n-button>
              </div>
              <Playlist :songs="playlist" :currentSong="currentSong" />
            </aside>
          </transition>
        </div>
      </n-dialog-provider>
    </n-message-provider>
  </n-config-provider>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { NConfigProvider, NMessageProvider, NDialogProvider, darkTheme, type GlobalThemeOverrides } from 'naive-ui'
import { VolumeHigh, VolumeLow, VolumeMute, Close as CloseIcon } from '@vicons/ionicons5'
import PlayerControls from './components/PlayerControls.vue'
import Playlist from './components/Playlist.vue'
import { usePlayerStore } from './stores/player'

const playerStore = usePlayerStore()

// Theme overrides for Naive UI
const themeOverrides: GlobalThemeOverrides = {
  common: {
    primaryColor: '#1DB954',
    primaryColorHover: '#1ed760',
    primaryColorPressed: '#1aa34a',
    borderRadius: '8px',
    fontFamily: "'Inter', -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif",
  },
  Button: {
    borderRadiusMedium: '9999px',
    borderRadiusSmall: '9999px',
    borderRadiusTiny: '9999px',
  },
  Input: {
    borderRadius: '500px',
  },
  Slider: {
    fillColor: '#1DB954',
    fillColorHover: '#1ed760',
    dotBorderColor: '#1DB954',
  },
}

// Player state
const currentSong = computed(() => playerStore.currentSong)
const playlist = computed(() => playerStore.playlist)
const volumeValue = ref(playerStore.volume)

const volumeIcon = computed(() => {
  if (playerStore.isMuted || volumeValue.value === 0) return VolumeMute
  if (volumeValue.value < 0.3) return VolumeLow
  return VolumeHigh
})

watch(() => playerStore.volume, (v) => { volumeValue.value = v })

function onVolumeChange(value: number) {
  playerStore.setVolume(value)
}

function toggleMute() {
  playerStore.toggleMute()
}

// Playlist panel
const showPlaylistPanel = ref(false)
function togglePlaylistPanel() {
  showPlaylistPanel.value = !showPlaylistPanel.value
}

// Quick playlists (placeholder)
const quickPlaylists = [
  { name: '我喜欢的音乐' },
  { name: '最近播放' },
  { name: '每日推荐' },
]
</script>

<style scoped>
/* App Layout */
.app-layout {
  display: flex;
  width: 100%;
  height: 100vh;
  overflow: hidden;
  background: var(--bg-base);
  position: relative;
}

/* Sidebar */
.sidebar {
  width: var(--sidebar-width);
  min-width: var(--sidebar-width);
  height: calc(100vh - var(--player-height));
  background: #000000;
  display: flex;
  flex-direction: column;
  overflow-y: auto;
  overflow-x: hidden;
}

.sidebar-logo {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 20px 24px 10px;
  color: var(--text-base);
}

.logo-text {
  font-size: 18px;
  font-weight: var(--font-weight-bold);
  letter-spacing: -0.02em;
}

.sidebar-logo svg {
  color: var(--spotify-green);
}

/* Navigation */
.sidebar-nav {
  display: flex;
  flex-direction: column;
  padding: 8px 12px;
  gap: 4px;
}

.nav-item {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 8px 12px;
  border-radius: var(--radius-sm);
  color: var(--text-subdued);
  font-size: var(--font-size-caption);
  font-weight: var(--font-weight-bold);
  transition: all 0.2s;
  position: relative;
  text-decoration: none;
}

.nav-item:hover {
  color: var(--text-base);
}

.nav-item.active {
  color: var(--text-base);
}

.nav-item.active::before {
  content: '';
  position: absolute;
  left: 0;
  top: 50%;
  transform: translateY(-50%);
  width: 3px;
  height: 24px;
  background: var(--spotify-green);
  border-radius: 0 2px 2px 0;
}

.nav-item svg {
  flex-shrink: 0;
}

.sidebar-divider {
  height: 1px;
  background: rgba(255, 255, 255, 0.08);
  margin: 8px 12px;
}

/* Sidebar section */
.sidebar-section {
  padding: 8px 12px;
  flex: 1;
  overflow-y: auto;
}

.sidebar-section-title {
  font-size: 11px;
  font-weight: var(--font-weight-bold);
  color: var(--text-subdued);
  text-transform: uppercase;
  letter-spacing: 1.5px;
  padding: 8px 12px;
}

.sidebar-playlists {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.playlist-link {
  padding: 8px 12px;
  border-radius: var(--radius-sm);
  color: var(--text-subdued);
  font-size: var(--font-size-caption);
  cursor: pointer;
  transition: color 0.2s;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.playlist-link:hover {
  color: var(--text-base);
}

/* Main Content */
.main-content {
  flex: 1;
  height: calc(100vh - var(--player-height));
  overflow: hidden;
  border-radius: 8px 0 0 0;
  margin-top: 0;
}

.content-scroll {
  height: 100%;
  overflow-y: auto;
  overflow-x: hidden;
  background: linear-gradient(180deg, #1a1a2e 0%, var(--bg-base) 30%);
}

/* Player Bar — Fixed Bottom */
.player-bar {
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
  height: var(--player-height);
  background: #181818;
  border-top: 1px solid rgba(255, 255, 255, 0.06);
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 16px;
  z-index: 1000;
}

/* Now Playing Info */
.player-song-info {
  display: flex;
  align-items: center;
  gap: 14px;
  min-width: 200px;
  width: 30%;
}

.now-playing-cover {
  width: 56px;
  height: 56px;
  flex-shrink: 0;
  border-radius: var(--radius-md);
  overflow: hidden;
}

.cover-img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.cover-placeholder {
  width: 100%;
  height: 100%;
  background: var(--bg-highlight);
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--text-subdued);
}

.now-playing-text {
  min-width: 0;
}

.now-playing-title {
  font-size: var(--font-size-caption);
  font-weight: var(--font-weight-bold);
  color: var(--text-base);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  max-width: 180px;
}

.now-playing-title:hover {
  text-decoration: underline;
  cursor: pointer;
}

.now-playing-artist {
  font-size: 12px;
  font-weight: var(--font-weight-regular);
  color: var(--text-subdued);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  max-width: 180px;
  margin-top: 2px;
}

.now-playing-artist:hover {
  color: var(--text-base);
  text-decoration: underline;
  cursor: pointer;
}

.now-playing-empty {
  display: flex;
  align-items: center;
}

/* Player Center Controls */
.player-center {
  flex: 1;
  display: flex;
  justify-content: center;
  max-width: 722px;
}

/* Player Right */
.player-right {
  display: flex;
  align-items: center;
  justify-content: flex-end;
  min-width: 200px;
  width: 30%;
}

.volume-control {
  display: flex;
  align-items: center;
  gap: 8px;
}

.volume-control:hover .volume-slider :deep(.n-slider-handle) {
  opacity: 1;
}

/* Playlist Side Panel */
.playlist-panel {
  position: fixed;
  right: 0;
  top: 0;
  bottom: var(--player-height);
  width: 360px;
  background: #121212;
  border-left: 1px solid rgba(255, 255, 255, 0.06);
  z-index: 900;
  display: flex;
  flex-direction: column;
  box-shadow: var(--shadow-dialog);
}

.playlist-panel-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 20px;
  border-bottom: 1px solid rgba(255, 255, 255, 0.06);
}

.playlist-panel-header h3 {
  font-size: var(--font-size-body);
  font-weight: var(--font-weight-bold);
  color: var(--text-base);
  margin: 0;
}

/* Slide Right Transition */
.slide-right-enter-active,
.slide-right-leave-active {
  transition: transform 0.3s ease;
}

.slide-right-enter-from,
.slide-right-leave-to {
  transform: translateX(100%);
}
</style>
