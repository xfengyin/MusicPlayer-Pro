<template>
  <div class="settings-view">
    <h1 class="page-title">设置</h1>

    <!-- Playback Section -->
    <section class="settings-section">
      <h2 class="section-label">播放</h2>
      <div class="setting-item">
        <div class="setting-info">
          <div class="setting-name">音频质量</div>
          <div class="setting-desc">选择流媒体音质</div>
        </div>
        <select class="setting-select">
          <option>自动</option>
          <option>低</option>
          <option>标准</option>
          <option>高</option>
        </select>
      </div>
      <div class="setting-item">
        <div class="setting-info">
          <div class="setting-name">淡入淡出</div>
          <div class="setting-desc">歌曲之间平滑过渡</div>
        </div>
        <div class="setting-toggle" :class="{ active: crossfade }" @click="crossfade = !crossfade">
          <div class="toggle-dot"></div>
        </div>
      </div>
      <div class="setting-item">
        <div class="setting-info">
          <div class="setting-name">无缝播放</div>
          <div class="setting-desc">消除歌曲间的间隙</div>
        </div>
        <div class="setting-toggle" :class="{ active: gapless }" @click="gapless = !gapless">
          <div class="toggle-dot"></div>
        </div>
      </div>
    </section>

    <!-- Display Section -->
    <section class="settings-section">
      <h2 class="section-label">显示</h2>
      <div class="setting-item">
        <div class="setting-info">
          <div class="setting-name">显示歌词</div>
          <div class="setting-desc">播放时显示歌词</div>
        </div>
        <div class="setting-toggle" :class="{ active: showLyrics }" @click="showLyrics = !showLyrics">
          <div class="toggle-dot"></div>
        </div>
      </div>
      <div class="setting-item">
        <div class="setting-info">
          <div class="setting-name">翻译歌词</div>
          <div class="setting-desc">显示歌词翻译</div>
        </div>
        <div class="setting-toggle" :class="{ active: showTranslation }" @click="showTranslation = !showTranslation">
          <div class="toggle-dot"></div>
        </div>
      </div>
    </section>

    <!-- About Section -->
    <section class="settings-section">
      <h2 class="section-label">关于</h2>
      <div class="setting-item">
        <div class="setting-info">
          <div class="setting-name">版本</div>
          <div class="setting-desc">Music Player Pro v1.0.0</div>
        </div>
      </div>
      <div class="setting-item">
        <div class="setting-info">
          <div class="setting-name">开源许可</div>
          <div class="setting-desc">查看第三方开源组件</div>
        </div>
        <svg viewBox="0 0 24 24" width="16" height="16" fill="var(--text-subdued)" style="flex-shrink: 0;"><path d="M8.59 16.59L13.17 12 8.59 7.41 10 6l6 6-6 6-1.41-1.41z"/></svg>
      </div>
    </section>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'

const crossfade = ref(false)
const gapless = ref(true)
const showLyrics = ref(true)
const showTranslation = ref(true)
</script>

<style scoped>
.settings-view {
  padding: 24px 32px 120px;
  min-height: 100%;
}

.page-title {
  font-size: var(--font-size-title);
  font-weight: var(--font-weight-bold);
  color: var(--text-base);
  margin: 0 0 32px;
  letter-spacing: -0.02em;
}

/* Settings Section */
.settings-section {
  margin-bottom: 32px;
}

.section-label {
  font-size: var(--font-size-small);
  font-weight: var(--font-weight-bold);
  color: var(--text-subdued);
  text-transform: uppercase;
  letter-spacing: 1.5px;
  margin: 0 0 16px;
  padding-bottom: 8px;
  border-bottom: 1px solid rgba(255, 255, 255, 0.06);
}

/* Setting Item */
.setting-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 16px 0;
  border-bottom: 1px solid rgba(255, 255, 255, 0.04);
}

.setting-info {
  flex: 1;
  min-width: 0;
}

.setting-name {
  font-size: var(--font-size-body);
  font-weight: var(--font-weight-regular);
  color: var(--text-base);
  margin-bottom: 2px;
}

.setting-desc {
  font-size: var(--font-size-small);
  color: var(--text-subdued);
}

/* Select Dropdown */
.setting-select {
  background: var(--bg-elevated);
  color: var(--text-base);
  border: 1px solid var(--border-subtle);
  border-radius: var(--radius-sm);
  padding: 8px 32px 8px 12px;
  font-size: var(--font-size-caption);
  font-family: inherit;
  cursor: pointer;
  appearance: none;
  background-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' width='12' height='12' fill='%23b3b3b3' viewBox='0 0 24 24'%3E%3Cpath d='M7 10l5 5 5-5z'/%3E%3C/svg%3E");
  background-repeat: no-repeat;
  background-position: right 8px center;
}

.setting-select:hover {
  border-color: var(--border-bright);
}

.setting-select:focus {
  outline: none;
  border-color: var(--text-base);
}

/* Toggle Switch — Spotify Green */
.setting-toggle {
  width: 44px;
  height: 24px;
  border-radius: 12px;
  background: rgba(255, 255, 255, 0.15);
  cursor: pointer;
  position: relative;
  transition: background 0.2s;
  flex-shrink: 0;
}

.setting-toggle.active {
  background: var(--spotify-green);
}

.toggle-dot {
  position: absolute;
  top: 2px;
  left: 2px;
  width: 20px;
  height: 20px;
  border-radius: var(--radius-circle);
  background: var(--text-base);
  transition: transform 0.2s;
}

.setting-toggle.active .toggle-dot {
  transform: translateX(20px);
}

@media (max-width: 768px) {
  .settings-view {
    padding: 16px;
  }
}
</style>
