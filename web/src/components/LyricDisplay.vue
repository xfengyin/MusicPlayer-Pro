<template>
  <div class="lyric-display" ref="containerRef">
    <div v-if="lyrics.length === 0" class="no-lyrics">
      <n-icon :component="MusicalNotesOutline" :size="64" />
      <p>{{ loading ? '加载歌词中...' : '暂无歌词' }}</p>
      <span v-if="!loading && currentSong">纯音乐或歌词不可用</span>
    </div>

    <div
      v-else
      class="lyrics-container"
      :style="{ transform: `translateY(${offsetY}px)` }"
    >
      <!-- 顶部占位，使第一行歌词能滚动到中间 -->
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

      <!-- 底部占位，使最后一行歌词能滚动到中间 -->
      <div class="lyric-spacer" :style="{ height: `${containerHeight / 2}px` }"></div>
    </div>

    <!-- 滚动控制按钮 -->
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

    <!-- 歌词设置 -->
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
import { MusicalNotesOutline, ChevronUp, ChevronDown, Text } from '@vicons/ionicons5'
import type { LyricLine } from '../types/music'

// 歌词行接口（扩展）
export interface LyricLineExtended extends LyricLine {
  translation?: string
}

// Props
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

// Emits
const emit = defineEmits<{
  'line-click': [index: number, time: number]
  'scroll-end': []
}>()

// 状态
const containerRef = ref<HTMLElement | null>(null)
const activeIndex = ref(0)
const containerHeight = ref(300)
const lineHeight = ref(80) // 每行高度（包括间距）
const loading = ref(false)
const fontSize = ref(18)
const lineRefs = ref<(HTMLElement | null)[]>([])

// 当前歌曲（用于显示空状态）
const currentSong = ref(true)

// 计算偏移量
const offsetY = computed(() => {
  if (props.lyrics.length === 0) return 0
  // 计算目标位置，使活跃行居中
  const targetY = -(activeIndex.value * lineHeight.value)
  return targetY
})

// 设置歌词行引用
function setLineRef(el: HTMLElement | null, index: number) {
  lineRefs.value[index] = el
}

// 查找当前活跃的歌词行
function findActiveIndex(time: number): number {
  if (props.lyrics.length === 0) return 0

  for (let i = props.lyrics.length - 1; i >= 0; i--) {
    if (time >= props.lyrics[i].time) {
      return i
    }
  }
  return 0
}

// 平滑滚动到活跃行
function scrollToActiveLine() {
  if (!props.autoScroll || props.lyrics.length === 0) return

  // 使用 CSS transform 已经实现平滑滚动
  // 这里可以添加额外的逻辑，如边界检测
  emit('scroll-end')
}

// 监听当前时间变化
watch(
  () => props.currentTime,
  (newTime) => {
    const newIndex = findActiveIndex(newTime)
    if (newIndex !== activeIndex.value) {
      activeIndex.value = newIndex
      scrollToActiveLine()
    }
  }
)

// 监听歌词变化
watch(
  () => props.lyrics,
  (newLyrics) => {
    loading.value = true
    activeIndex.value = findActiveIndex(props.currentTime)
    nextTick(() => {
      updateContainerHeight()
      if (props.autoScroll && newLyrics.length > 0) {
        scrollToActiveLine()
      }
      loading.value = false
    })
  },
  { immediate: true }
)

// 更新容器高度
function updateContainerHeight() {
  if (containerRef.value) {
    containerHeight.value = containerRef.value.clientHeight
  }
}

// 点击歌词行
function onLineClick(index: number, time: number) {
  activeIndex.value = index
  emit('line-click', index, time)
}

// 手动滚动控制
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

// 切换字体大小
function toggleFontSize() {
  const sizes = [16, 18, 20, 22, 24]
  const currentIndex = sizes.indexOf(fontSize.value)
  fontSize.value = sizes[(currentIndex + 1) % sizes.length]
}

// 暴露方法给父组件
defineExpose({
  scrollToLine: (index: number) => {
    activeIndex.value = index
    scrollToActiveLine()
  },
  reset: () => {
    activeIndex.value = 0
  },
})

// 生命周期
onMounted(() => {
  updateContainerHeight()
  // 监听窗口大小变化
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
  border-radius: 8px;
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
  color: rgba(255, 255, 255, 0.6);
  font-weight: 500;
  transition: all 0.3s ease;
  text-shadow: 0 2px 4px rgba(0, 0, 0, 0.3);
}

.lyric-line.active .lyric-text {
  color: rgba(255, 255, 255, 0.98);
  font-size: v-bind('(fontSize + 4) + "px"');
  font-weight: 600;
  text-shadow: 0 4px 12px rgba(255, 255, 255, 0.3), 0 2px 4px rgba(0, 0, 0, 0.4);
}

.lyric-translation {
  font-size: v-bind('(fontSize - 4) + "px"');
  line-height: 1.5;
  color: rgba(255, 255, 255, 0.4);
  transition: all 0.3s ease;
}

.lyric-line.active .lyric-translation {
  color: rgba(255, 255, 255, 0.75);
  font-size: v-bind('(fontSize - 2) + "px"');
}

.no-lyrics {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 16px;
  color: rgba(255, 255, 255, 0.4);
  padding: 40px 20px;
}

.no-lyrics p {
  margin: 0;
  font-size: 16px;
  font-weight: 500;
  color: rgba(255, 255, 255, 0.6);
}

.no-lyrics span {
  font-size: 13px;
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

/* 响应式 */
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

/* 滚动条隐藏 */
.lyric-display::-webkit-scrollbar {
  display: none;
}

/* 禁用状态 */
.lyric-controls:deep(.n-button--disabled),
.lyric-settings:deep(.n-button--disabled) {
  opacity: 0.3;
}
</style>
