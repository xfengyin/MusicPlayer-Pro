<template>
  <div class="lyric-display" ref="containerRef">
    <div v-if="lyrics.length === 0" class="no-lyrics">
      <svg viewBox="0 0 24 24" width="64" height="64" fill="currentColor" style="color: var(--text-subdued); opacity: 0.3;"><path d="M12 3v10.55c-.59-.34-1.27-.55-2-.55C7.79 13 6 14.79 6 17s1.79 4 4 4 4-1.79 4-4V7h4V3h-6z"/></svg>
      <p>{{ loading ? '加载歌词中...' : '暂无歌词' }}</p>
      <span v-if="!loading && currentSong">纯音乐或歌词不可用</span>
    </div>

    <div
      v-else
      class="lyrics-container"
      :style="{ transform: `translateY(${offsetY}px)` }"
    >
      <div class="lyric-spacer" :style="{ height: `${containerHeight / 2}px` }"></div>

      <div
        v-for="(line, index) in lyrics"
        :key="index"
        :ref="(el) => setLineRef(el, index)"
        :class="['lyric-line', { active: index === activeIndex }]"
        @click="onLineClick(index, line.time)"
      >
        <div class="lyric-text">{{ line.text }}</div>
        <div v-if="line.translation" class="lyric-translation">
          {{ line.translation }}
        </div>
      </div>

      <div class="lyric-spacer" :style="{ height: `${containerHeight / 2}px` }"></div>
    </div>

    <!-- Scroll controls -->
    <div v-if="showControls" class="lyric-controls">
      <n-button
        quaternary
        circle
        size="small"
        @click="scrollUp"
        title="向上滚动"
        :disabled="activeIndex === 0"
      >
        <template #icon>
          <n-icon :component="ChevronUp" />
        </template>
      </n-button>
      <n-button
        quaternary
        circle
        size="small"
        @click="scrollDown"
        title="向下滚动"
        :disabled="activeIndex === lyrics.length - 1"
      >
        <template #icon>
          <n-icon :component="ChevronDown" />
        </template>
      </n-button>
    </div>

    <!-- Settings -->
    <div v-if="showSettings" class="lyric-settings">
      <n-button
        quaternary
        circle
        size="small"
        @click="toggleFontSize"
        :title="`字体大小：${fontSize}`"
      >
        <template #icon>
          <n-icon :component="Text" :size="18" />
        </template>
      </n-button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch, nextTick, onMounted, onBeforeUnmount } from 'vue'
import { NButton, NIcon } from 'naive-ui'
import { ChevronUp, ChevronDown, Text } from '@vicons/ionicons5'
import type { LyricLine } from '../types/music'

export interface LyricLineExtended extends LyricLine {
  translation?: string
}

interface Props {
  lyrics: LyricLineExtended[]
  currentTime: number
  showControls?: boolean
  showSettings?: boolean
  autoScroll?: boolean
  highlightOffset?: number
}

const props = withDefaults(defineProps<Props>(), {
  lyrics: () => [],
  currentTime: 0,
  showControls: false,
  showSettings: false,
  autoScroll: true,
  highlightOffset: 0,
})

const emit = defineEmits<{
  'line-click': [index: number, time: number]
  'scroll-end': []
}>()

const containerRef = ref<HTMLElement | null>(null)
const activeIndex = ref(0)
const containerHeight = ref(300)
const lineHeight = ref(80)
const loading = ref(false)
const fontSize = ref(18)
const lineRefs = ref<(HTMLElement | null)[]>([])
const currentSong = ref(true)

const offsetY = computed(() => {
  if (props.lyrics.length === 0) return 0
  return -(activeIndex.value * lineHeight.value)
})

function setLineRef(el: HTMLElement | null, index: number) {
  lineRefs.value[index] = el
}

function findActiveIndex(time: number): number {
  if (props.lyrics.length === 0) return 0
  for (let i = props.lyrics.length - 1; i >= 0; i--) {
    if (time >= props.lyrics[i].time) return i
  }
  return 0
}

function scrollToActiveLine() {
  if (!props.autoScroll || props.lyrics.length === 0) return
  emit('scroll-end')
}

watch(() => props.currentTime, (newTime) => {
  const newIndex = findActiveIndex(newTime)
  if (newIndex !== activeIndex.value) {
    activeIndex.value = newIndex
    scrollToActiveLine()
  }
})

watch(() => props.lyrics, (newLyrics) => {
  loading.value = true
  activeIndex.value = findActiveIndex(props.currentTime)
  nextTick(() => {
    updateContainerHeight()
    if (props.autoScroll && newLyrics.length > 0) scrollToActiveLine()
    loading.value = false
  })
}, { immediate: true })

function updateContainerHeight() {
  if (containerRef.value) containerHeight.value = containerRef.value.clientHeight
}

function onLineClick(index: number, time: number) {
  activeIndex.value = index
  emit('line-click', index, time)
}

function scrollUp() {
  if (activeIndex.value > 0) {
    activeIndex.value--
    scrollToActiveLine()
  }
}

function scrollDown() {
  if (activeIndex.value < props.lyrics.length - 1) {
    activeIndex.value++
    scrollToActiveLine()
  }
}

function toggleFontSize() {
  const sizes = [16, 18, 20, 22, 24]
  const currentIndex = sizes.indexOf(fontSize.value)
  fontSize.value = sizes[(currentIndex + 1) % sizes.length]
}

defineExpose({
  scrollToLine: (index: number) => {
    activeIndex.value = index
    scrollToActiveLine()
  },
  reset: () => {
    activeIndex.value = 0
  },
})

onMounted(() => {
  updateContainerHeight()
  window.addEventListener('resize', updateContainerHeight)
})

onBeforeUnmount(() => {
  window.removeEventListener('resize', updateContainerHeight)
})
</script>

<style scoped>
.lyric-display {
  position: relative;
  height: 100%;
  overflow: hidden;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 40px 20px;
}

.lyrics-container {
  transition: transform 0.5s cubic-bezier(0.4, 0, 0.2, 1);
  will-change: transform;
  width: 100%;
}

.lyric-spacer {
  width: 100%;
  flex-shrink: 0;
}

.lyric-line {
  padding: 16px 20px;
  text-align: center;
  cursor: pointer;
  transition: all 0.3s ease;
  min-height: v-bind('lineHeight + "px"');
  display: flex;
  flex-direction: column;
  justify-content: center;
  gap: 8px;
  border-radius: var(--radius-lg);
  margin: 0 20px;
}

.lyric-line:hover {
  background: rgba(255, 255, 255, 0.05);
}

.lyric-line.active {
  transform: scale(1.05);
}

.lyric-text {
  font-size: v-bind('fontSize + "px"');
  line-height: 1.6;
  color: var(--text-subdued);
  font-weight: var(--font-weight-regular);
  transition: all 0.3s ease;
}

.lyric-line.active .lyric-text {
  color: var(--text-base);
  font-size: v-bind('(fontSize + 4) + "px"');
  font-weight: var(--font-weight-bold);
}

.lyric-translation {
  font-size: v-bind('(fontSize - 4) + "px"');
  line-height: 1.5;
  color: rgba(255, 255, 255, 0.35);
  transition: all 0.3s ease;
}

.lyric-line.active .lyric-translation {
  color: var(--text-subdued);
  font-size: v-bind('(fontSize - 2) + "px"');
}

.no-lyrics {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 16px;
  color: var(--text-subdued);
  padding: 40px 20px;
}

.no-lyrics p {
  margin: 0;
  font-size: var(--font-size-body);
  font-weight: var(--font-weight-bold);
  color: var(--text-base);
}

.no-lyrics span {
  font-size: var(--font-size-small);
  color: var(--text-subdued);
}

.lyric-controls {
  position: absolute;
  right: 20px;
  top: 50%;
  transform: translateY(-50%);
  display: flex;
  flex-direction: column;
  gap: 8px;
  opacity: 0;
  transition: opacity 0.2s;
}

.lyric-settings {
  position: absolute;
  right: 20px;
  bottom: 20px;
  opacity: 0;
  transition: opacity 0.2s;
}

.lyric-display:hover .lyric-controls,
.lyric-display:hover .lyric-settings {
  opacity: 1;
}

@media (max-width: 768px) {
  .lyric-text {
    font-size: v-bind('(fontSize - 2) + "px"');
  }
  .lyric-line.active .lyric-text {
    font-size: v-bind('(fontSize + 2) + "px"');
  }
  .lyric-translation {
    font-size: v-bind('(fontSize - 6) + "px"');
  }
  .lyric-line.active .lyric-translation {
    font-size: v-bind('(fontSize - 4) + "px"');
  }
  .lyric-line {
    margin: 0 10px;
  }
}

.lyric-display::-webkit-scrollbar {
  display: none;
}
</style>
