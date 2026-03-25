<template>
  <div class="playlist-view">
    <!-- Header -->
    <div class="playlist-header safe-area-top">
      <n-button quaternary circle @click="$router.back()">
        <template #icon>
          <n-icon :component="ArrowBack" :size="22" />
        </template>
      </n-button>
      <h1 class="page-title">播放列表</h1>
      <n-button
        v-if="playerStore.playlist.length > 0"
        quaternary
        circle
        @click="confirmClear"
      >
        <template #icon>
          <n-icon :component="TrashOutline" :size="20" />
        </template>
      </n-button>
    </div>

    <!-- Playlist Content -->
    <div class="playlist-content">
      <PlaylistComponent
        :songs="playerStore.playlist"
        :currentSong="playerStore.currentSong"
        @play="onPlay"
        @remove="onRemove"
        @clear-all="confirmClear"
        @add="goToSearch"
        @reorder="onReorder"
      />
    </div>

    <!-- Bottom Player -->
    <div v-if="playerStore.currentSong" class="bottom-player safe-area-bottom">
      <PlayerControls :showVolume="false" @toggle-playlist="() => {}" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router'
import { NButton, NIcon, useDialog, useMessage } from 'naive-ui'
import { ArrowBack, TrashOutline } from '@vicons/ionicons5'
import PlaylistComponent from '@/components/Playlist.vue'
import PlayerControls from '@/components/PlayerControls.vue'
import { usePlayerStore } from '@/stores/player'
import type { Song } from '@/types/music'

const router = useRouter()
const playerStore = usePlayerStore()
const dialog = useDialog()
const message = useMessage()

async function onPlay(song: Song, index: number) {
  await playerStore.playSong(song, index)
}

function onRemove(index: number) {
  const song = playerStore.playlist[index]
  playerStore.removeFromPlaylist(index)
  message.success(`已移除「${song?.title || '歌曲'}」`)
}

function onReorder(from: number, to: number) {
  playerStore.reorderPlaylist(from, to)
}

function confirmClear() {
  dialog.warning({
    title: '清空播放列表',
    content: `确定要清空全部 ${playerStore.playlist.length} 首歌曲吗？`,
    positiveText: '清空',
    negativeText: '取消',
    onPositiveClick: () => {
      playerStore.clearPlaylist()
      message.success('播放列表已清空')
    },
  })
}

function goToSearch() {
  router.push('/search')
}
</script>

<style scoped>
.playlist-view {
  display: flex;
  flex-direction: column;
  height: 100vh;
  overflow: hidden;
}

.playlist-header {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px 16px;
  background: rgba(0, 0, 0, 0.2);
  backdrop-filter: blur(10px);
  z-index: 10;
}

.page-title {
  flex: 1;
  font-size: 18px;
  font-weight: 600;
  margin: 0;
}

.playlist-content {
  flex: 1;
  overflow-y: auto;
  padding-bottom: 120px;
}

.bottom-player {
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
  background: rgba(15, 12, 41, 0.98);
  backdrop-filter: blur(20px);
  border-top: 1px solid rgba(255, 255, 255, 0.1);
  z-index: 50;
}
</style>
