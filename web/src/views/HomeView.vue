<template>
  <div class="home-view">
    <!-- Header -->
    <div class="header safe-area-top">
      <h1 class="app-title">🎵 MusicPlayer Pro</h1>
      <SearchBar
        @search="onSearch"
        @select="onSongSelect"
        @show-more="goToSearch"
      />
    </div>

    <!-- Main Content -->
    <div class="main-content">
      <!-- Now Playing -->
      <div v-if="playerStore.currentSong" class="now-playing-section">
        <div class="section-title">正在播放</div>
        <div class="now-playing-card">
          <img
            v-if="playerStore.currentSong.cover"
            :src="playerStore.currentSong.cover"
            class="now-playing-cover"
            @error="onImageError"
          />
          <div v-else class="now-playing-cover-placeholder">
            <n-icon :component="MusicalNote" :size="48" />
          </div>
          <div class="now-playing-info">
            <div class="now-playing-title">{{ playerStore.currentSong.title }}</div>
            <div class="now-playing-artist">{{ playerStore.currentSong.artist }}</div>
          </div>
        </div>
      </div>

      <!-- Quick Actions -->
      <div class="quick-actions">
        <div class="action-card" @click="goToSearch">
          <n-icon :component="SearchIcon" :size="32" />
          <span>搜索音乐</span>
        </div>
        <div class="action-card" @click="goToPlaylist">
          <n-icon :component="ListIcon" :size="32" />
          <span>播放列表</span>
          <span v-if="playerStore.playlist.length > 0" class="badge">
            {{ playerStore.playlist.length }}
          </span>
        </div>
        <div class="action-card" @click="goToSettings">
          <n-icon :component="SettingsIcon" :size="32" />
          <span>设置</span>
        </div>
      </div>

      <!-- Lyrics Display -->
      <div v-if="playerStore.currentSong" class="lyrics-section">
        <LyricDisplay
          :lyrics="playerStore.lyrics"
          :currentTime="playerStore.currentTime"
          :autoScroll="true"
        />
      </div>

      <!-- Empty State -->
      <div v-if="!playerStore.currentSong" class="empty-state">
        <n-icon :component="MusicalNotesOutline" :size="80" />
        <h2>欢迎使用 MusicPlayer Pro</h2>
        <p>搜索并播放你喜欢的音乐</p>
        <n-button type="primary" size="large" @click="goToSearch">
          开始搜索
        </n-button>
      </div>
    </div>

    <!-- Player Controls (Bottom Bar) -->
    <div class="player-bar safe-area-bottom">
      <PlayerControls :showVolume="!isMobile" @toggle-playlist="showPlaylistDrawer = true" />
    </div>

    <!-- Playlist Drawer -->
    <n-drawer v-model:show="showPlaylistDrawer" :width="360" placement="right">
      <n-drawer-content title="播放列表" :native-scrollbar="false">
        <PlaylistComponent
          :songs="playerStore.playlist"
          :currentSong="playerStore.currentSong"
          @play="onPlayFromPlaylist"
          @remove="onRemoveFromPlaylist"
          @clear-all="onClearPlaylist"
          @add="goToSearch"
        />
      </n-drawer-content>
    </n-drawer>

    <!-- Bottom Navigation -->
    <div class="bottom-nav safe-area-bottom">
      <div class="nav-item active" @click="$router.push('/')">
        <n-icon :component="Home" :size="22" />
        <span>首页</span>
      </div>
      <div class="nav-item" @click="$router.push('/search')">
        <n-icon :component="SearchIcon" :size="22" />
        <span>搜索</span>
      </div>
      <div class="nav-item" @click="$router.push('/playlist')">
        <n-icon :component="ListIcon" :size="22" />
        <span>歌单</span>
      </div>
      <div class="nav-item" @click="$router.push('/settings')">
        <n-icon :component="SettingsIcon" :size="22" />
        <span>设置</span>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { NIcon, NButton, NDrawer, NDrawerContent } from 'naive-ui'
import {
  MusicalNote,
  MusicalNotesOutline,
  Search as SearchIcon,
  List as ListIcon,
  Settings as SettingsIcon,
  Home,
} from '@vicons/ionicons5'
import SearchBar from '@/components/SearchBar.vue'
import PlayerControls from '@/components/PlayerControls.vue'
import LyricDisplay from '@/components/LyricDisplay.vue'
import PlaylistComponent from '@/components/Playlist.vue'
import { usePlayerStore } from '@/stores/player'
import type { Song } from '@/types/music'

const router = useRouter()
const playerStore = usePlayerStore()
const showPlaylistDrawer = ref(false)

const isMobile = computed(() => window.innerWidth < 768)

function onSearch(query: string) {
  router.push({ path: '/search', query: { q: query } })
}

async function onSongSelect(song: Song) {
  playerStore.addToPlaylist(song)
  await playerStore.playSong(song)
}

function goToSearch() {
  router.push('/search')
}

function goToPlaylist() {
  router.push('/playlist')
}

function goToSettings() {
  router.push('/settings')
}

async function onPlayFromPlaylist(song: Song, index: number) {
  await playerStore.playSong(song, index)
}

function onRemoveFromPlaylist(index: number) {
  playerStore.removeFromPlaylist(index)
}

function onClearPlaylist() {
  playerStore.clearPlaylist()
}

function onImageError(e: Event) {
  const target = e.target as HTMLImageElement
  target.style.display = 'none'
}
</script>

<style scoped>
.home-view {
  display: flex;
  flex-direction: column;
  height: 100vh;
  overflow: hidden;
}

.header {
  padding: 16px 20px;
  display: flex;
  flex-direction: column;
  gap: 12px;
  background: rgba(0, 0, 0, 0.2);
  backdrop-filter: blur(10px);
  z-index: 10;
}

.app-title {
  font-size: 22px;
  font-weight: 700;
  background: linear-gradient(135deg, #6366f1, #a78bfa);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
}

.main-content {
  flex: 1;
  overflow-y: auto;
  padding: 20px;
  padding-bottom: 180px;
}

.section-title {
  font-size: 14px;
  font-weight: 600;
  color: rgba(255, 255, 255, 0.6);
  margin-bottom: 12px;
  text-transform: uppercase;
  letter-spacing: 1px;
}

.now-playing-section {
  margin-bottom: 24px;
}

.now-playing-card {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 16px;
  background: rgba(255, 255, 255, 0.08);
  border-radius: 12px;
  backdrop-filter: blur(10px);
}

.now-playing-cover {
  width: 64px;
  height: 64px;
  border-radius: 8px;
  object-fit: cover;
}

.now-playing-cover-placeholder {
  width: 64px;
  height: 64px;
  border-radius: 8px;
  background: rgba(255, 255, 255, 0.1);
  display: flex;
  align-items: center;
  justify-content: center;
  color: rgba(255, 255, 255, 0.3);
}

.now-playing-info {
  flex: 1;
  min-width: 0;
}

.now-playing-title {
  font-size: 16px;
  font-weight: 600;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.now-playing-artist {
  font-size: 13px;
  color: rgba(255, 255, 255, 0.5);
  margin-top: 4px;
}

.quick-actions {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 12px;
  margin-bottom: 24px;
}

.action-card {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
  padding: 20px 12px;
  background: rgba(255, 255, 255, 0.06);
  border-radius: 12px;
  cursor: pointer;
  transition: all 0.2s;
  position: relative;
}

.action-card:hover {
  background: rgba(255, 255, 255, 0.12);
  transform: translateY(-2px);
}

.action-card span {
  font-size: 13px;
  color: rgba(255, 255, 255, 0.8);
}

.badge {
  position: absolute;
  top: 8px;
  right: 8px;
  background: #6366f1;
  color: white;
  font-size: 11px;
  padding: 2px 6px;
  border-radius: 10px;
  min-width: 20px;
  text-align: center;
}

.lyrics-section {
  height: 300px;
  margin-bottom: 24px;
  background: rgba(255, 255, 255, 0.04);
  border-radius: 12px;
  overflow: hidden;
}

.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 16px;
  padding: 60px 20px;
  color: rgba(255, 255, 255, 0.5);
  text-align: center;
}

.empty-state h2 {
  font-size: 20px;
  font-weight: 600;
  color: rgba(255, 255, 255, 0.8);
}

.empty-state p {
  font-size: 14px;
}

.player-bar {
  position: fixed;
  bottom: 56px;
  left: 0;
  right: 0;
  z-index: 50;
  background: rgba(15, 12, 41, 0.95);
  backdrop-filter: blur(20px);
  border-top: 1px solid rgba(255, 255, 255, 0.1);
}

.bottom-nav {
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
  display: flex;
  justify-content: space-around;
  background: rgba(15, 12, 41, 0.98);
  backdrop-filter: blur(20px);
  border-top: 1px solid rgba(255, 255, 255, 0.1);
  z-index: 60;
  padding: 8px 0;
}

.nav-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 4px;
  padding: 6px 16px;
  cursor: pointer;
  color: rgba(255, 255, 255, 0.5);
  transition: color 0.2s;
}

.nav-item.active {
  color: #6366f1;
}

.nav-item:hover {
  color: rgba(255, 255, 255, 0.8);
}

.nav-item span {
  font-size: 11px;
}

@media (max-width: 768px) {
  .header {
    padding: 12px 16px;
  }

  .main-content {
    padding: 16px;
  }

  .quick-actions {
    grid-template-columns: repeat(3, 1fr);
    gap: 8px;
  }

  .action-card {
    padding: 16px 8px;
  }
}
</style>
