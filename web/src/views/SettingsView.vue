<template>
  <div class="settings-view">
    <!-- Header -->
    <div class="settings-header safe-area-top">
      <n-button quaternary circle @click="$router.back()">
        <template #icon>
          <n-icon :component="ArrowBack" :size="22" />
        </template>
      </n-button>
      <h1 class="page-title">设置</h1>
    </div>

    <!-- Settings Content -->
    <div class="settings-content">
      <!-- API 配置 -->
      <div class="settings-group">
        <div class="group-title">🔌 服务器配置</div>
        <div class="setting-item">
          <div class="setting-label">
            <span class="label-text">API 服务器地址</span>
            <span class="label-desc">APK 需要连接远程后端服务</span>
          </div>
          <n-input
            v-model:value="apiUrl"
            placeholder="例如: https://your-server.com/api"
            @blur="saveApiUrl"
            clearable
          />
        </div>
        <div class="setting-item">
          <n-button type="primary" size="small" @click="testConnection" :loading="isTesting">
            测试连接
          </n-button>
          <span v-if="connectionStatus === 'ok'" class="status-ok">✅ 连接正常</span>
          <span v-else-if="connectionStatus === 'fail'" class="status-fail">❌ 连接失败</span>
        </div>
      </div>

      <!-- 播放设置 -->
      <div class="settings-group">
        <div class="group-title">🎵 播放设置</div>
        <div class="setting-item">
          <span class="label-text">默认音量</span>
          <n-slider
            v-model:value="volume"
            :min="0"
            :max="1"
            :step="0.05"
            :tooltip="true"
            :format-tooltip="(v: number) => Math.round(v * 100) + '%'"
            @update:value="onVolumeChange"
            style="width: 200px"
          />
        </div>
        <div class="setting-item">
          <span class="label-text">播放模式</span>
          <n-radio-group v-model:value="loopMode" @update:value="onLoopModeChange">
            <n-radio value="list">列表循环</n-radio>
            <n-radio value="single">单曲循环</n-radio>
            <n-radio value="random">随机播放</n-radio>
          </n-radio-group>
        </div>
      </div>

      <!-- 关于 -->
      <div class="settings-group">
        <div class="group-title">ℹ️ 关于</div>
        <div class="setting-item">
          <span class="label-text">应用名称</span>
          <span class="label-value">MusicPlayer Pro</span>
        </div>
        <div class="setting-item">
          <span class="label-text">版本</span>
          <span class="label-value">v2.1.0</span>
        </div>
        <div class="setting-item">
          <span class="label-text">运行环境</span>
          <span class="label-value">{{ platformInfo }}</span>
        </div>
        <div class="setting-item">
          <span class="label-text">GitHub</span>
          <a href="https://github.com/xfengyin/MusicPlayer-Pro" target="_blank" class="link">
            xfengyin/MusicPlayer-Pro
          </a>
        </div>
      </div>

      <!-- 数据管理 -->
      <div class="settings-group">
        <div class="group-title">🗃️ 数据管理</div>
        <div class="setting-item">
          <n-button type="warning" size="small" @click="clearSearchHistory">
            清除搜索历史
          </n-button>
          <n-button type="error" size="small" @click="clearAllData">
            清除所有数据
          </n-button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { NButton, NIcon, NInput, NSlider, NRadioGroup, NRadio, useMessage, useDialog } from 'naive-ui'
import { ArrowBack } from '@vicons/ionicons5'
import { usePlayerStore } from '@/stores/player'
import { getPlatform, isNative } from '@/utils/platform'
import { refreshApiClient } from '@/api/music'
import axios from 'axios'

const playerStore = usePlayerStore()
const message = useMessage()
const dialog = useDialog()

const apiUrl = ref('')
const volume = ref(playerStore.volume)
const loopMode = ref(playerStore.loopMode)
const isTesting = ref(false)
const connectionStatus = ref<'none' | 'ok' | 'fail'>('none')

const platformInfo = computed(() => {
  const platform = getPlatform()
  return isNative() ? `原生 App (${platform})` : `Web 浏览器`
})

onMounted(() => {
  apiUrl.value = localStorage.getItem('api_base_url') || ''
})

function saveApiUrl() {
  const url = apiUrl.value.trim()
  if (url) {
    localStorage.setItem('api_base_url', url)
    refreshApiClient()
    message.success('API 地址已保存')
  } else {
    localStorage.removeItem('api_base_url')
    refreshApiClient()
  }
  connectionStatus.value = 'none'
}

async function testConnection() {
  const url = apiUrl.value.trim() || import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080/api'
  isTesting.value = true
  connectionStatus.value = 'none'

  try {
    const baseUrl = url.replace(/\/api\/?$/, '')
    await axios.get(`${baseUrl}/health`, { timeout: 5000 })
    connectionStatus.value = 'ok'
    message.success('连接成功！')
  } catch (error) {
    connectionStatus.value = 'fail'
    message.error('连接失败，请检查地址是否正确')
  } finally {
    isTesting.value = false
  }
}

function onVolumeChange(v: number) {
  playerStore.setVolume(v)
}

function onLoopModeChange(mode: string) {
  playerStore.setLoopMode(mode as 'list' | 'single' | 'random')
  const labels: Record<string, string> = { list: '列表循环', single: '单曲循环', random: '随机播放' }
  message.success(`已切换到${labels[mode]}`)
}

function clearSearchHistory() {
  localStorage.removeItem('searchHistory')
  message.success('搜索历史已清除')
}

function clearAllData() {
  dialog.warning({
    title: '清除所有数据',
    content: '这将清除播放列表、搜索历史和所有本地设置。确定吗？',
    positiveText: '确定清除',
    negativeText: '取消',
    onPositiveClick: () => {
      playerStore.clearPlaylist()
      localStorage.clear()
      message.success('所有数据已清除')
    },
  })
}
</script>

<style scoped>
.settings-view {
  display: flex;
  flex-direction: column;
  height: 100vh;
  overflow: hidden;
}

.settings-header {
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

.settings-content {
  flex: 1;
  overflow-y: auto;
  padding: 16px;
}

.settings-group {
  margin-bottom: 24px;
  background: rgba(255, 255, 255, 0.05);
  border-radius: 12px;
  padding: 16px;
}

.group-title {
  font-size: 15px;
  font-weight: 600;
  margin-bottom: 16px;
  color: rgba(255, 255, 255, 0.9);
}

.setting-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 10px 0;
  border-bottom: 1px solid rgba(255, 255, 255, 0.06);
  flex-wrap: wrap;
}

.setting-item:last-child {
  border-bottom: none;
}

.setting-label {
  flex: 1;
  min-width: 0;
}

.label-text {
  font-size: 14px;
  color: rgba(255, 255, 255, 0.85);
}

.label-desc {
  display: block;
  font-size: 12px;
  color: rgba(255, 255, 255, 0.4);
  margin-top: 4px;
}

.label-value {
  font-size: 14px;
  color: rgba(255, 255, 255, 0.5);
}

.link {
  font-size: 14px;
  color: #6366f1;
  text-decoration: none;
}

.link:hover {
  text-decoration: underline;
}

.status-ok {
  font-size: 13px;
  color: #22c55e;
}

.status-fail {
  font-size: 13px;
  color: #ef4444;
}
</style>
