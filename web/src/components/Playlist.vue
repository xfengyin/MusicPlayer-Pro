<template>
  <div class="playlist">
    <div class="playlist-header">
      <h3 class="playlist-title">
        <svg viewBox="0 0 24 24" width="20" height="20" fill="currentColor" style="color: var(--spotify-green);"><path d="M3 4h18v2H3V4zm0 7h12v2H3v-2zm0 7h18v2H3v-2zm14-3l6-4v8l-6-4z"/></svg>
        播放列表
        <span class="playlist-count">({{ songs.length }} 首)</span>
      </h3>
      <div class="playlist-actions">
        <n-button quaternary circle size="small" @click="onAddSong" title="添加歌曲">
          <template #icon>
            <n-icon :component="Add" />
          </template>
        </n-button>
        <n-button
          quaternary
          circle
          size="small"
          :disabled="songs.length === 0"
          @click="onClearAll"
          title="清空列表"
        >
          <template #icon>
            <n-icon :component="Trash" />
          </template>
        </n-button>
        <n-dropdown
          :options="modeOptions"
          :value="loopMode"
          @select="onModeSelect"
          trigger="click"
          placement="bottom-end"
        >
          <n-button quaternary circle size="small" :title="loopModeText">
            <template #icon>
              <n-icon :component="loopModeIcon" :size="18" />
            </template>
          </n-button>
        </n-dropdown>
      </div>
    </div>

    <div class="playlist-content">
      <transition-group name="list" tag="div" class="song-list">
        <div
          v-for="(song, index) in songs"
          :key="song.id"
          :ref="(el) => setSongRef(el, index)"
          :class="[
            'song-item',
            {
              'is-playing': currentSong?.id === song.id,
              'is-dragging': isDragging && dragIndex === index,
            },
          ]"
          draggable
          @dragstart="onDragStart(index, $event)"
          @dragover.prevent="onDragOver(index)"
          @dragleave="onDragLeave"
          @drop="onDrop(index)"
          @dragend="onDragEnd"
          @click="onPlaySong(song, index)"
          @contextmenu="onContextMenu($event, song, index)"
        >
          <div class="song-drag-handle">
            <n-icon :component="ReorderFour" :size="16" />
          </div>

          <div class="song-index">
            <span v-if="currentSong?.id !== song.id" class="index-num">{{ index + 1 }}</span>
            <div v-else class="playing-indicator">
              <span class="bar"></span>
              <span class="bar"></span>
              <span class="bar"></span>
            </div>
          </div>

          <div class="song-cover-wrapper">
            <img
              v-if="song.cover"
              :src="song.cover"
              :alt="song.title"
              class="song-cover"
              @error="onImageError"
            />
            <div v-else class="song-cover-placeholder">
              <svg viewBox="0 0 24 24" width="20" height="20" fill="currentColor"><path d="M12 3v10.55c-.59-.34-1.27-.55-2-.55C7.79 13 6 14.79 6 17s1.79 4 4 4 4-1.79 4-4V7h4V3h-6z"/></svg>
            </div>
            <div class="song-cover-play">
              <svg viewBox="0 0 24 24" width="14" height="14" fill="#000"><path d="M8 5.14v13.72a1 1 0 001.5.86l11.04-6.86a1 1 0 000-1.72L9.5 4.28a1 1 0 00-1.5.86z"/></svg>
            </div>
          </div>

          <div class="song-info">
            <div class="song-title" :class="{ 'active-title': currentSong?.id === song.id }">{{ song.title }}</div>
            <div class="song-artist">{{ song.artist }}</div>
          </div>

          <div class="song-album">{{ song.album }}</div>

          <div class="song-duration">{{ formatTime(song.duration) }}</div>

          <button
            class="remove-btn"
            @click.stop="onRemoveSong(index)"
          >
            <svg viewBox="0 0 24 24" width="14" height="14" fill="currentColor"><path d="M18 6L6 18M6 6l12 12" stroke="currentColor" stroke-width="2" stroke-linecap="round"/></svg>
          </button>
        </div>
      </transition-group>

      <div v-if="songs.length === 0" class="empty-state">
        <svg viewBox="0 0 24 24" width="64" height="64" fill="var(--text-subdued)" style="opacity: 0.3;"><path d="M3 4h18v2H3V4zm0 7h12v2H3v-2zm0 7h18v2H3v-2zm14-3l6-4v8l-6-4z"/></svg>
        <p>播放列表为空</p>
        <button class="add-first-btn" @click="onAddSong">添加第一首歌曲</button>
      </div>
    </div>

    <!-- Right-click menu -->
    <n-dropdown
      :show="showContextMenu"
      :options="contextMenuOptions"
      :x="contextMenuX"
      :y="contextMenuY"
      @select="onContextMenuSelect"
      @clickoutside="hideContextMenu"
      placement="bottom-start"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch, h } from 'vue'
import { NButton, NIcon, NDropdown, type DropdownOption } from 'naive-ui'
import {
  MusicalNote,
  Add,
  Trash,
  ReorderFour,
  Repeat,
  RepeatOne,
  Shuffle,
} from '@vicons/ionicons5'
import type { Song } from '../types/music'
import { usePlayerStore } from '../stores/player'

interface Props {
  songs: Song[]
  currentSong?: Song | null
}

const props = withDefaults(defineProps<Props>(), {
  songs: () => [],
  currentSong: null,
})

const emit = defineEmits<{
  play: [song: Song, index: number]
  remove: [index: number]
  'clear-all': []
  add: []
  reorder: [fromIndex: number, toIndex: number]
  'context-menu': [event: MouseEvent, song: Song, index: number]
}>()

const playerStore = usePlayerStore()

const isDragging = ref(false)
const dragIndex = ref(-1)
const dragOverIndex = ref(-1)
const songRefs = ref<(HTMLElement | null)[]>([])

const showContextMenu = ref(false)
const contextMenuX = ref(0)
const contextMenuY = ref(0)
const contextMenuSong = ref<Song | null>(null)
const contextMenuIndex = ref(-1)

const modeOptions = computed(() => [
  {
    label: '列表循环',
    key: 'list',
    icon: () => h(NIcon, null, { default: () => h(Repeat) }),
  },
  {
    label: '单曲循环',
    key: 'single',
    icon: () => h(NIcon, null, { default: () => h(RepeatOne) }),
  },
  {
    label: '随机播放',
    key: 'random',
    icon: () => h(NIcon, null, { default: () => h(Shuffle) }),
  },
])

const loopMode = computed(() => playerStore.loopMode)
const loopModeIcon = computed(() => {
  switch (playerStore.loopMode.value) {
    case 'single': return RepeatOne
    case 'random': return Shuffle
    default: return Repeat
  }
})
const loopModeText = computed(() => {
  switch (playerStore.loopMode.value) {
    case 'single': return '单曲循环'
    case 'random': return '随机播放'
    default: return '列表循环'
  }
})

const contextMenuOptions = computed(() => [
  {
    label: '播放这首',
    key: 'play',
    icon: () => h(NIcon, null, { default: () => h(MusicalNote) }),
  },
  {
    label: '从列表中移除',
    key: 'remove',
    icon: () => h(NIcon, null, { default: () => h(Trash) }),
  },
  {
    label: '移到顶部',
    key: 'move-top',
    disabled: contextMenuIndex.value === 0,
  },
  {
    label: '移到底部',
    key: 'move-bottom',
    disabled: contextMenuIndex.value === props.songs.length - 1,
  },
])

function setSongRef(el: HTMLElement | null, index: number) {
  songRefs.value[index] = el
}

function formatTime(seconds: number): string {
  if (!seconds || isNaN(seconds)) return '-:--'
  const mins = Math.floor(seconds / 60)
  const secs = Math.floor(seconds % 60)
  return `${mins}:${secs.toString().padStart(2, '0')}`
}

function onImageError(e: Event) {
  const target = e.target as HTMLImageElement
  target.style.display = 'none'
}

function onDragStart(index: number, event: DragEvent) {
  isDragging.value = true
  dragIndex.value = index
  if (event.dataTransfer) {
    event.dataTransfer.effectAllowed = 'move'
    event.dataTransfer.setData('text/plain', index.toString())
  }
}

function onDragOver(index: number) {
  if (dragIndex.value === -1 || dragIndex.value === index) return
  dragOverIndex.value = index
}

function onDragLeave() {
  dragOverIndex.value = -1
}

function onDrop(toIndex: number) {
  if (dragIndex.value !== -1 && dragIndex.value !== toIndex) {
    emit('reorder', dragIndex.value, toIndex)
    playerStore.reorderPlaylist(dragIndex.value, toIndex)
  }
  onDragEnd()
}

function onDragEnd() {
  isDragging.value = false
  dragIndex.value = -1
  dragOverIndex.value = -1
}

async function onPlaySong(song: Song, index: number) {
  await playerStore.playSong(song, index)
  emit('play', song, index)
}

function onRemoveSong(index: number) {
  playerStore.removeFromPlaylist(index)
  emit('remove', index)
}

function onClearAll() {
  playerStore.clearPlaylist()
  emit('clear-all')
}

function onAddSong() {
  emit('add')
}

function onModeSelect(mode: string) {
  playerStore.setLoopMode(mode as 'list' | 'single' | 'random')
}

function onContextMenu(event: MouseEvent, song: Song, index: number) {
  event.preventDefault()
  contextMenuSong.value = song
  contextMenuIndex.value = index
  contextMenuX.value = event.clientX
  contextMenuY.value = event.clientY
  showContextMenu.value = true
}

function hideContextMenu() {
  showContextMenu.value = false
  contextMenuSong.value = null
  contextMenuIndex.value = -1
}

function onContextMenuSelect(key: string) {
  if (contextMenuSong.value === null || contextMenuIndex.value === -1) return

  switch (key) {
    case 'play':
      onPlaySong(contextMenuSong.value, contextMenuIndex.value)
      break
    case 'remove':
      onRemoveSong(contextMenuIndex.value)
      break
    case 'move-top':
      playerStore.reorderPlaylist(contextMenuIndex.value, 0)
      emit('reorder', contextMenuIndex.value, 0)
      break
    case 'move-bottom':
      playerStore.reorderPlaylist(contextMenuIndex.value, props.songs.length - 1)
      emit('reorder', contextMenuIndex.value, props.songs.length - 1)
      break
  }

  hideContextMenu()
}

watch(
  () => props.songs,
  () => {},
  { deep: true }
)
</script>

<style scoped>
.playlist {
  display: flex;
  flex-direction: column;
  height: 100%;
  overflow: hidden;
}

.playlist-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 20px;
  border-bottom: 1px solid rgba(255, 255, 255, 0.06);
}

.playlist-title {
  display: flex;
  align-items: center;
  gap: 8px;
  margin: 0;
  font-size: var(--font-size-body);
  font-weight: var(--font-weight-bold);
  color: var(--text-base);
}

.playlist-count {
  font-size: var(--font-size-small);
  font-weight: var(--font-weight-regular);
  color: var(--text-subdued);
}

.playlist-actions {
  display: flex;
  gap: 4px;
}

.playlist-content {
  flex: 1;
  overflow-y: auto;
  padding: 0 8px;
}

.song-list {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

/* Song Item — Spotify Table Row Style */
.song-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 8px 12px;
  border-radius: var(--radius-sm);
  cursor: pointer;
  transition: background 0.2s;
  position: relative;
  user-select: none;
}

.song-item:hover {
  background: rgba(255, 255, 255, 0.08);
}

.song-item:hover .index-num {
  display: none;
}

.song-item:hover .song-drag-handle {
  opacity: 0;
}

.song-item:hover .remove-btn {
  opacity: 1;
}

.song-item:hover .song-cover-play {
  opacity: 1;
}

.song-item:hover .song-album {
  color: var(--text-base);
}

.song-item.is-playing {
  background: rgba(29, 185, 84, 0.1);
}

.song-item.is-dragging {
  opacity: 0.5;
  background: rgba(255, 255, 255, 0.1);
}

/* Drag Handle */
.song-drag-handle {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 24px;
  color: rgba(255, 255, 255, 0.3);
  cursor: grab;
  transition: opacity 0.2s;
}

.song-drag-handle:active {
  cursor: grabbing;
}

/* Song Index */
.song-index {
  width: 24px;
  text-align: center;
  font-size: var(--font-size-caption);
  color: var(--text-subdued);
  display: flex;
  align-items: center;
  justify-content: center;
}

.index-num {
  font-variant-numeric: tabular-nums;
}

/* Playing indicator — animated bars (Spotify style) */
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

.playing-indicator .bar:nth-child(1) {
  height: 6px;
  animation-delay: 0s;
}

.playing-indicator .bar:nth-child(2) {
  height: 10px;
  animation-delay: 0.2s;
}

.playing-indicator .bar:nth-child(3) {
  height: 4px;
  animation-delay: 0.4s;
}

@keyframes bar-bounce {
  0%, 100% { transform: scaleY(1); }
  50% { transform: scaleY(0.4); }
}

/* Song Cover */
.song-cover-wrapper {
  position: relative;
  width: 40px;
  height: 40px;
  flex-shrink: 0;
}

.song-cover {
  width: 100%;
  height: 100%;
  border-radius: var(--radius-sm);
  object-fit: cover;
}

.song-cover-placeholder {
  width: 100%;
  height: 100%;
  border-radius: var(--radius-sm);
  background: var(--bg-highlight);
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--text-subdued);
}

.song-cover-play {
  position: absolute;
  inset: 0;
  background: rgba(0, 0, 0, 0.6);
  border-radius: var(--radius-sm);
  display: flex;
  align-items: center;
  justify-content: center;
  opacity: 0;
  transition: opacity 0.2s;
}

/* Song Info */
.song-info {
  flex: 1;
  min-width: 0;
}

.song-title {
  font-size: var(--font-size-caption);
  color: var(--text-base);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  font-weight: var(--font-weight-regular);
}

.song-title.active-title {
  color: var(--spotify-green);
}

.song-artist {
  font-size: var(--font-size-small);
  color: var(--text-subdued);
  margin-top: 1px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.song-artist:hover {
  text-decoration: underline;
}

/* Song Album */
.song-album {
  font-size: var(--font-size-small);
  color: var(--text-subdued);
  width: 160px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  transition: color 0.2s;
}

/* Song Duration */
.song-duration {
  font-size: var(--font-size-small);
  color: var(--text-subdued);
  min-width: 44px;
  text-align: right;
  font-variant-numeric: tabular-nums;
}

/* Remove Button */
.remove-btn {
  background: none;
  border: none;
  cursor: pointer;
  color: var(--text-subdued);
  padding: 4px;
  border-radius: var(--radius-circle);
  opacity: 0;
  transition: all 0.2s;
  display: flex;
  align-items: center;
}

.remove-btn:hover {
  color: var(--text-base);
}

/* Empty State */
.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 300px;
  gap: 16px;
}

.empty-state p {
  margin: 0;
  font-size: var(--font-size-caption);
  color: var(--text-subdued);
}

.add-first-btn {
  background: none;
  border: none;
  cursor: pointer;
  color: var(--text-base);
  font-size: var(--font-size-caption);
  font-weight: var(--font-weight-bold);
  text-decoration: underline;
  text-underline-offset: 3px;
}

.add-first-btn:hover {
  color: var(--spotify-green);
}

/* List animations */
.list-enter-active,
.list-leave-active {
  transition: all 0.3s ease;
  position: absolute;
  width: calc(100% - 16px);
}

.list-enter-from {
  opacity: 0;
  transform: translateX(-30px);
}

.list-leave-to {
  opacity: 0;
  transform: translateX(30px);
}

.list-move {
  transition: transform 0.3s ease;
}

/* Scrollbar */
.playlist-content::-webkit-scrollbar {
  width: 6px;
}

.playlist-content::-webkit-scrollbar-track {
  background: transparent;
}

.playlist-content::-webkit-scrollbar-thumb {
  background: rgba(255, 255, 255, 0.15);
  border-radius: 3px;
}

.playlist-content::-webkit-scrollbar-thumb:hover {
  background: rgba(255, 255, 255, 0.3);
}
</style>
