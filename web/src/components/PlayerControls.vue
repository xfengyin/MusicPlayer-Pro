<template>
  <div class="player-controls">
    <!-- Progress Bar -->
    <div class="progress-section">
      <span class="time-label">{{ formatTime(currentTime) }}</span>
      <div class="progress-bar-wrapper">
        <n-slider
          v-model:value="progressValue"
          :min="0"
          :max="duration || 1"
          :step="0.1"
          :tooltip="false"
          :disabled="!currentSong"
          @update:value="onProgressChange"
          @mousedown="onProgressStart"
          @mouseup="onProgressEnd"
          class="progress-slider"
        />
      </div>
      <span class="time-label">{{ formatTime(duration) }}</span>
    </div>

    <!-- Control Buttons -->
    <div class="control-buttons">
      <!-- Shuffle / Loop Mode -->
      <button class="ctrl-btn" :class="{ active: loopMode !== 'list' }" @click="toggleLoopMode" :title="loopModeText">
        <n-icon :component="loopModeIcon" :size="18" />
      </button>

      <!-- Previous -->
      <button class="ctrl-btn" :disabled="!hasPrevSong" @click="onPrevSong" title="上一首">
        <svg viewBox="0 0 24 24" width="20" height="20" fill="currentColor"><path d="M6 3a1 1 0 00-1 1v16a1 1 0 102 0V4a1 1 0 00-1-1zm13.707 1.293a1 1 0 00-1.414 0L10 12.586l8.293 8.293a1 1 0 001.414-1.414L12.414 12l7.293-7.293a1 1 0 000-1.414z"/></svg>
      </button>

      <!-- Play / Pause — Main Circle -->
      <button class="play-btn" :disabled="!currentSong" @click="onTogglePlay" title="播放/暂停">
        <svg v-if="isPlaying && !isLoading" viewBox="0 0 24 24" width="24" height="24" fill="currentColor"><path d="M5.7 3a.7.7 0 00-.7.7v16.6a.7.7 0 00.7.7h2.6a.7.7 0 00.7-.7V3.7a.7.7 0 00-.7-.7H5.7zm10 0a.7.7 0 00-.7.7v16.6a.7.7 0 00.7.7h2.6a.7.7 0 00.7-.7V3.7a.7.7 0 00-.7-.7h-2.6z"/></svg>
        <svg v-else-if="!isLoading" viewBox="0 0 24 24" width="24" height="24" fill="currentColor"><path d="M8 5.14v13.72a1 1 0 001.5.86l11.04-6.86a1 1 0 000-1.72L9.5 4.28a1 1 0 00-1.5.86z"/></svg>
      </button>

      <!-- Next -->
      <button class="ctrl-btn" :disabled="!hasNextSong" @click="onNextSong" title="下一首">
        <svg viewBox="0 0 24 24" width="20" height="20" fill="currentColor"><path d="M18 3a1 1 0 011 1v16a1 1 0 11-2 0V4a1 1 0 011-1zM4.293 4.293a1 1 0 011.414 0L14 12.586l-8.293 8.293a1 1 0 01-1.414-1.414L11.586 12 4.293 4.707a1 1 0 010-1.414z"/></svg>
      </button>

      <!-- Playlist Toggle -->
      <button class="ctrl-btn" @click="onTogglePlaylist" title="播放列表">
        <svg viewBox="0 0 24 24" width="20" height="20" fill="currentColor"><path d="M3 4h18v2H3V4zm0 7h12v2H3v-2zm0 7h18v2H3v-2zm14-3l6-4v8l-6-4z"/></svg>
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch, onMounted, onBeforeUnmount } from 'vue'
import { NIcon, NSlider } from 'naive-ui'
import {
  Repeat,
  RepeatOne,
  Shuffle,
} from '@vicons/ionicons5'
import { usePlayerStore } from '../stores/player'
import { message } from 'naive-ui'

interface Props {
  showVolume?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  showVolume: true,
})

const emit = defineEmits<{
  'toggle-playlist': []
}>()

const playerStore = usePlayerStore()

const progressValue = ref(0)
const isDraggingProgress = ref(false)

const {
  isPlaying,
  currentTime,
  duration,
  hasNextSong,
  hasPrevSong,
  currentSong,
  playlist,
  isLoading,
  loopMode,
} = playerStore

const loopModeIcon = computed(() => {
  switch (loopMode.value) {
    case 'single': return RepeatOne
    case 'random': return Shuffle
    default: return Repeat
  }
})

const loopModeText = computed(() => {
  switch (loopMode.value) {
    case 'single': return '单曲循环'
    case 'random': return '随机播放'
    default: return '列表循环'
  }
})

watch(() => playerStore.currentTime, (newTime) => {
  if (!isDraggingProgress.value) {
    progressValue.value = newTime
  }
})

function formatTime(seconds: number): string {
  if (!seconds || isNaN(seconds)) return '-:--'
  const mins = Math.floor(seconds / 60)
  const secs = Math.floor(seconds % 60)
  return `${mins}:${secs.toString().padStart(2, '0')}`
}

function onTogglePlay() {
  playerStore.togglePlay()
}

async function onNextSong() {
  await playerStore.nextSong()
}

async function onPrevSong() {
  await playerStore.prevSong()
}

function onProgressChange(value: number) {
  if (isDraggingProgress.value) {
    progressValue.value = value
  }
}

function onProgressStart() {
  isDraggingProgress.value = true
}

function onProgressEnd() {
  isDraggingProgress.value = false
  playerStore.seekTo(progressValue.value)
}

function toggleLoopMode() {
  const modes: Array<'list' | 'single' | 'random'> = ['list', 'single', 'random']
  const currentIndex = modes.indexOf(loopMode.value)
  const nextMode = modes[(currentIndex + 1) % modes.length]
  playerStore.setLoopMode(nextMode)
  message.success(`已切换到${loopModeText.value}`, { duration: 1000 })
}

function onTogglePlaylist() {
  emit('toggle-playlist')
}

// Keyboard shortcuts
function handleKeyDown(event: KeyboardEvent) {
  if (event.code === 'Space' && event.target === document.body) {
    event.preventDefault()
    onTogglePlay()
  }
  if (event.code === 'ArrowLeft' && event.target === document.body) {
    event.preventDefault()
    onPrevSong()
  }
  if (event.code === 'ArrowRight' && event.target === document.body) {
    event.preventDefault()
    onNextSong()
  }
  if (event.code === 'ArrowUp' && event.target === document.body) {
    event.preventDefault()
    playerStore.setVolume(Math.min(1, playerStore.volume + 0.1))
  }
  if (event.code === 'ArrowDown' && event.target === document.body) {
    event.preventDefault()
    playerStore.setVolume(Math.max(0, playerStore.volume - 0.1))
  }
}

onMounted(() => {
  window.addEventListener('keydown', handleKeyDown)
})

onBeforeUnmount(() => {
  window.removeEventListener('keydown', handleKeyDown)
})
</script>

<style scoped>
.player-controls {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
  width: 100%;
}

/* Progress Section */
.progress-section {
  display: flex;
  align-items: center;
  gap: 8px;
  width: 100%;
}

.time-label {
  font-size: 11px;
  color: var(--text-subdued);
  min-width: 40px;
  text-align: center;
  font-variant-numeric: tabular-nums;
  font-weight: var(--font-weight-regular);
}

.progress-bar-wrapper {
  flex: 1;
  cursor: pointer;
}

.progress-slider {
  width: 100%;
}

.progress-slider :deep(.n-slider-rail) {
  height: 4px;
  border-radius: 2px;
  background: rgba(255, 255, 255, 0.15);
}

.progress-slider :deep(.n-slider-fill) {
  height: 4px;
  border-radius: 2px;
  background: var(--text-base);
  transition: background 0.2s;
}

.progress-bar-wrapper:hover .progress-slider :deep(.n-slider-fill) {
  background: var(--spotify-green);
}

.progress-slider :deep(.n-slider-handle) {
  width: 12px;
  height: 12px;
  border: 0;
  background: var(--text-base);
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.5);
  opacity: 0;
  transition: opacity 0.2s;
}

.progress-bar-wrapper:hover .progress-slider :deep(.n-slider-handle) {
  opacity: 1;
}

/* Control Buttons */
.control-buttons {
  display: flex;
  align-items: center;
  gap: 16px;
}

.ctrl-btn {
  background: none;
  border: none;
  cursor: pointer;
  color: var(--text-subdued);
  padding: 8px;
  border-radius: var(--radius-circle);
  transition: all 0.2s;
  display: flex;
  align-items: center;
  justify-content: center;
  line-height: 1;
}

.ctrl-btn:hover {
  color: var(--text-base);
  transform: scale(1.05);
}

.ctrl-btn.active {
  color: var(--spotify-green);
}

.ctrl-btn:disabled {
  opacity: 0.4;
  cursor: not-allowed;
}

.ctrl-btn:disabled:hover {
  color: var(--text-subdued);
  transform: none;
}

/* Main Play Button — White Circle */
.play-btn {
  background: var(--text-base);
  border: none;
  cursor: pointer;
  width: 36px;
  height: 36px;
  border-radius: var(--radius-circle);
  display: flex;
  align-items: center;
  justify-content: center;
  color: #000;
  transition: transform 0.1s;
  padding: 0;
}

.play-btn:hover {
  transform: scale(1.06);
}

.play-btn:active {
  transform: scale(0.98);
}

.play-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.play-btn svg {
  display: block;
}

/* Disabled slider */
.player-controls :deep(.n-slider--disabled) {
  opacity: 0.4;
  cursor: not-allowed;
}
</style>
