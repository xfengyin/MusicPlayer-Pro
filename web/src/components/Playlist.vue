<template>
  <div class="playlist">
    <div class="playlist-header">
      <h3 class="playlist-title">
        <n-icon :component="MusicalNotes" />
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
            <span v-if="currentSong?.id !== song.id">{{ index + 1 }}</span>
            <n-icon v-else :component="VolumeHigh" :size="18" class="playing-icon" />
          </div>

          <img
            v-if="song.cover"
            :src="song.cover"
            :alt="song.title"
            class="song-cover"
            @error="onImageError"
          />
          <div v-else class="song-cover-placeholder">
            <n-icon :component="MusicalNote" :size="24" />
          </div>

          <div class="song-info">
            <div class="song-title">{{ song.title }}</div>
            <div class="song-artist">{{ song.artist }}</div>
          </div>

          <div class="song-duration">{{ formatTime(song.duration) }}</div>

          <n-button
            quaternary
            circle
            size="tiny"
            class="remove-btn"
            @click.stop="onRemoveSong(index)"
          >
            <template #icon>
              <n-icon :component="Close" />
            </template>
          </n-button>
        </div>
      </transition-group>

      <div v-if="songs.length === 0" class="empty-state">
        <n-icon :component="MusicalNotesOutline" :size="64" />
        <p>播放列表为空</p>
        <n-button text type="primary" @click="onAddSong">
          添加第一首歌曲
        </n-button>
      </div>
    </div>

    <!-- 右键菜单 -->
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
  MusicalNotes,
  MusicalNote,
  MusicalNotesOutline,
  VolumeHigh,
  Add,
  Trash,
  Close,
  ReorderFour,
  Repeat,
  RepeatOne,
  Shuffle,
} from '@vicons/ionicons5'
import type { Song } from '../types/music'
import { usePlayerStore } from '../stores/player'

// Props
interface Props {
  songs: Song[]
  currentSong?: Song | null
}

const props = withDefaults(defineProps<Props>(), {
  songs: () => [],
  currentSong: null,
})

// Emits
const emit = defineEmits<{
  play: [song: Song, index: number]
  remove: [index: number]
  'clear-all': []
  add: []
  reorder: [fromIndex: number, toIndex: number]
  'context-menu': [event: MouseEvent, song: Song, index: number]
}>()

// Store
const playerStore = usePlayerStore()

// 拖拽状态
const isDragging = ref(false)
const dragIndex = ref(-1)
const dragOverIndex = ref(-1)
const songRefs = ref<(HTMLElement | null)[]>([])

// 右键菜单状态
const showContextMenu = ref(false)
const contextMenuX = ref(0)
const contextMenuY = ref(0)
const contextMenuSong = ref<Song | null>(null)
const contextMenuIndex = ref(-1)

// 播放模式选项
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
    case 'single':
      return RepeatOne
    case 'random':
      return Shuffle
    default:
      return Repeat
  }
})
const loopModeText = computed(() => {
  switch (playerStore.loopMode.value) {
    case 'single':
      return '单曲循环'
    case 'random':
      return '随机播放'
    default:
      return '列表循环'
  }
})

// 右键菜单选项
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

// 设置歌曲引用
function setSongRef(el: HTMLElement | null, index: number) {
  songRefs.value[index] = el
}

// 格式化时间
function formatTime(seconds: number): string {
  if (!seconds || isNaN(seconds)) return '00:00'
  const mins = Math.floor(seconds / 60)
  const secs = Math.floor(seconds % 60)
  return `${mins.toString().padStart(2, '0')}:${secs.toString().padStart(2, '0')}`
}

// 图片加载失败处理
function onImageError(e: Event) {
  const target = e.target as HTMLImageElement
  target.style.display = 'none'
}

// 拖拽事件处理
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

// 播放控制
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

// 右键菜单
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
      playerStore.reorderPlaylist(
        contextMenuIndex.value,
        props.songs.length - 1
      )
      emit('reorder', contextMenuIndex.value, props.songs.length - 1)
      break
  }

  hideContextMenu()
}

// 监听播放列表变化并持久化
watch(
  () => props.songs,
  () => {
    // 播放列表变化时自动保存到 localStorage（在 store 中已处理）
  },
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
  padding: 16px;
  border-bottom: 1px solid rgba(255, 255, 255, 0.1);
}

.playlist-title {
  display: flex;
  align-items: center;
  gap: 8px;
  margin: 0;
  font-size: 16px;
  font-weight: 600;
  color: rgba(255, 255, 255, 0.9);
}

.playlist-count {
  font-size: 12px;
  font-weight: 400;
  color: rgba(255, 255, 255, 0.5);
}

.playlist-actions {
  display: flex;
  gap: 8px;
}

.playlist-content {
  flex: 1;
  overflow-y: auto;
  padding: 8px;
}

.song-list {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.song-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 10px 12px;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.2s;
  position: relative;
  user-select: none;
}

.song-item:hover {
  background: rgba(255, 255, 255, 0.08);
}

.song-item.is-playing {
  background: rgba(99, 102, 241, 0.25);
}

.song-item.is-dragging {
  opacity: 0.5;
  background: rgba(255, 255, 255, 0.1);
  transform: scale(1.02);
}

.song-drag-handle {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 24px;
  color: rgba(255, 255, 255, 0.3);
  cursor: grab;
}

.song-drag-handle:hover {
  color: rgba(255, 255, 255, 0.6);
}

.song-drag-handle:active {
  cursor: grabbing;
}

.song-index {
  width: 24px;
  text-align: center;
  font-size: 14px;
  color: rgba(255, 255, 255, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
}

.playing-icon {
  color: #6366f1;
  animation: pulse 1.5s ease-in-out infinite;
}

@keyframes pulse {
  0%, 100% {
    opacity: 1;
    transform: scale(1);
  }
  50% {
    opacity: 0.7;
    transform: scale(1.1);
  }
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
  font-size: 14px;
  color: rgba(255, 255, 255, 0.95);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  font-weight: 500;
}

.song-artist {
  font-size: 12px;
  color: rgba(255, 255, 255, 0.5);
  margin-top: 2px;
}

.song-duration {
  font-size: 12px;
  color: rgba(255, 255, 255, 0.4);
  padding: 0 12px;
  font-variant-numeric: tabular-nums;
}

.remove-btn {
  opacity: 0;
  transition: opacity 0.2s;
}

.song-item:hover .remove-btn {
  opacity: 1;
}

.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 300px;
  color: rgba(255, 255, 255, 0.4);
  gap: 16px;
}

.empty-state p {
  margin: 0;
  font-size: 14px;
}

/* 列表动画 */
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

/* 滚动条样式 */
.playlist-content::-webkit-scrollbar {
  width: 6px;
}

.playlist-content::-webkit-scrollbar-track {
  background: rgba(255, 255, 255, 0.05);
  border-radius: 3px;
}

.playlist-content::-webkit-scrollbar-thumb {
  background: rgba(255, 255, 255, 0.2);
  border-radius: 3px;
}

.playlist-content::-webkit-scrollbar-thumb:hover {
  background: rgba(255, 255, 255, 0.3);
}
</style>
