<template>
  <div class="player-controls">
    <!-- 进度条 -->
    <div class="progress-section">
      <span class="time-label">{{ formatTime(currentTime) }}</span>
      <n-slider
        v-model:value="progressValue"
        :min="0"
        :max="duration"
        :step="0.1"
        :tooltip="false"
        :disabled="!currentSong"
        @update:value="onProgressChange"
        @mousedown="onProgressStart"
        @mouseup="onProgressEnd"
        class="progress-slider"
      />
      <span class="time-label">{{ formatTime(duration) }}</span>
    </div>

    <!-- 控制按钮 -->
    <div class="control-buttons">
      <!-- 播放模式 -->
      <n-button
        quaternary
        circle
        @click="toggleLoopMode"
        :title="loopModeText"
      >
        <template #icon>
          <n-icon :component="loopModeIcon" :size="20" />
        </template>
      </n-button>

      <n-button
        quaternary
        circle
        :disabled="!hasPrevSong"
        @click="onPrevSong"
        title="上一首"
      >
        <template #icon>
          <n-icon :component="SkipBack" :size="24" />
        </template>
      </n-button>

      <n-button
        tertiary
        circle
        :loading="isLoading"
        :disabled="!currentSong"
        @click="onTogglePlay"
        class="play-button"
        title="播放/暂停"
      >
        <template #icon>
          <n-icon :component="isPlaying ? Pause : Play" :size="32" />
        </template>
      </n-button>

      <n-button
        quaternary
        circle
        :disabled="!hasNextSong"
        @click="onNextSong"
        title="下一首"
      >
        <template #icon>
          <n-icon :component="SkipForward" :size="24" />
        </template>
      </n-button>

      <!-- 播放列表 -->
      <n-button
        quaternary
        circle
        :disabled="playlist.length === 0"
        @click="onTogglePlaylist"
        title="播放列表"
      >
        <template #icon>
          <n-icon :component="List" :size="20" />
        </template>
      </n-button>
    </div>

    <!-- 音量控制 -->
    <div v-if="showVolume" class="volume-section">
      <n-icon :component="volumeIcon" :size="20" @click="toggleMute" style="cursor: pointer" />
      <n-slider
        v-model:value="volumeValue"
        :min="0"
        :max="1"
        :step="0.01"
        :tooltip="false"
        @update:value="onVolumeChange"
        class="volume-slider"
        style="width: 100px"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch, onMounted, onBeforeUnmount } from 'vue'
import { NButton, NIcon, NSlider } from 'naive-ui'
import {
  Play,
  Pause,
  SkipBack,
  SkipForward,
  VolumeHigh,
  VolumeLow,
  VolumeMute,
  Repeat,
  RepeatOne,
  Shuffle,
  List,
} from '@vicons/ionicons5'
import { usePlayerStore } from '../stores/player'
import { message } from 'naive-ui'

// Props
interface Props {
  showVolume?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  showVolume: true,
})

// Emits
const emit = defineEmits<{
  'toggle-playlist': []
}>()

// 使用 Pinia store
const playerStore = usePlayerStore()

// 本地状态
const progressValue = ref(0)
const volumeValue = ref(playerStore.volume)
const isDraggingProgress = ref(false)

// 从 store 获取状态
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

// 音量图标
const volumeIcon = computed(() => {
  if (playerStore.isMuted || volumeValue.value === 0) return VolumeMute
  if (volumeValue.value < 0.3) return VolumeLow
  return VolumeHigh
})

// 播放模式图标和文本
const loopModeIcon = computed(() => {
  switch (loopMode.value) {
    case 'single':
      return RepeatOne
    case 'random':
      return Shuffle
    default:
      return Repeat
  }
})

const loopModeText = computed(() => {
  switch (loopMode.value) {
    case 'single':
      return '单曲循环'
    case 'random':
      return '随机播放'
    default:
      return '列表循环'
  }
})

// 监听 store 时间变化（仅在非拖动状态时更新）
watch(
  () => playerStore.currentTime,
  (newTime) => {
    if (!isDraggingProgress.value) {
      progressValue.value = newTime
    }
  }
)

// 监听音量变化
watch(
  () => playerStore.volume,
  (newVolume) => {
    volumeValue.value = newVolume
  }
)

// 格式化时间
function formatTime(seconds: number): string {
  if (!seconds || isNaN(seconds)) return '00:00'
  const mins = Math.floor(seconds / 60)
  const secs = Math.floor(seconds % 60)
  return `${mins.toString().padStart(2, '0')}:${secs.toString().padStart(2, '0')}`
}

// 事件处理
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

function onVolumeChange(value: number) {
  playerStore.setVolume(value)
}

function toggleMute() {
  playerStore.toggleMute()
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

// 键盘快捷键
function handleKeyDown(event: KeyboardEvent) {
  // 空格键：播放/暂停
  if (event.code === 'Space' && event.target === document.body) {
    event.preventDefault()
    onTogglePlay()
  }
  // 左箭头：上一首
  if (event.code === 'ArrowLeft' && event.target === document.body) {
    event.preventDefault()
    onPrevSong()
  }
  // 右箭头：下一首
  if (event.code === 'ArrowRight' && event.target === document.body) {
    event.preventDefault()
    onNextSong()
  }
  // 上箭头：音量 +
  if (event.code === 'ArrowUp' && event.target === document.body) {
    event.preventDefault()
    playerStore.setVolume(Math.min(1, volumeValue.value + 0.1))
  }
  // 下箭头：音量 -
  if (event.code === 'ArrowDown' && event.target === document.body) {
    event.preventDefault()
    playerStore.setVolume(Math.max(0, volumeValue.value - 0.1))
  }
}

// 生命周期
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
  gap: 16px;
  padding: 16px;
  background: rgba(255, 255, 255, 0.05);
  border-radius: 12px;
  backdrop-filter: blur(10px);
}

.progress-section {
  display: flex;
  align-items: center;
  gap: 12px;
  width: 100%;
  max-width: 600px;
}

.time-label {
  font-size: 12px;
  color: rgba(255, 255, 255, 0.6);
  min-width: 45px;
  text-align: center;
  font-variant-numeric: tabular-nums;
}

.progress-slider {
  flex: 1;
  cursor: pointer;
}

.control-buttons {
  display: flex;
  align-items: center;
  gap: 12px;
}

.play-button {
  width: 56px;
  height: 56px;
}

.volume-section {
  display: flex;
  align-items: center;
  gap: 8px;
}

.volume-slider {
  cursor: pointer;
}

/* 禁用状态 */
.player-controls:deep(.n-slider--disabled) {
  opacity: 0.5;
  cursor: not-allowed;
}

.player-controls:deep(.n-button--disabled) {
  opacity: 0.5;
}
</style>
