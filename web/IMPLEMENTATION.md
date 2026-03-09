# MusicPlayer-Pro 前端实现文档

## 概述

本项目实现了 MusicPlayer-Pro 的完整前端功能，使用 Vue 3 + TypeScript + Naive UI 构建。

## 已实现功能

### 1. 播放器控制 (PlayerControls.vue)

**功能特性：**
- ✅ 集成 HTML5 Audio API 实现真实播放控制
- ✅ 播放/暂停/上一首/下一首控制
- ✅ 进度条拖动跳转功能
- ✅ 音量控制（支持鼠标滚轮和拖动）
- ✅ 音量持久化（localStorage）
- ✅ 播放模式切换（列表循环/单曲循环/随机播放）
- ✅ 键盘快捷键支持（空格、方向键）
- ✅ 加载状态显示

**快捷键：**
- `空格`: 播放/暂停
- `←`: 上一首
- `→`: 下一首
- `↑`: 音量 +
- `↓`: 音量 -

### 2. 搜索栏 (SearchBar.vue)

**功能特性：**
- ✅ 搜索防抖（可配置，默认 300ms）
- ✅ 搜索历史本地存储
- ✅ 搜索结果实时展示
- ✅ 点击搜索结果直接播放
- ✅ 搜索历史管理（单个删除/清空）
- ✅ 加载更多结果
- ✅ 无结果提示
- ✅ 加载状态指示器

**本地存储：**
- 搜索历史保存在 `localStorage.searchHistory`
- 最多保存 10 条历史记录

### 3. 播放列表 (Playlist.vue)

**功能特性：**
- ✅ 拖拽排序（原生 HTML5 Drag & Drop）
- ✅ 播放列表持久化（localStorage）
- ✅ 右键菜单功能
  - 播放这首
  - 从列表中移除
  - 移到顶部
  - 移到底部
- ✅ 播放模式切换下拉菜单
- ✅ 当前播放高亮显示
- ✅ 动画过渡效果
- ✅ 空状态提示

**持久化数据：**
- 播放列表：`localStorage.player_playlist`
- 当前索引：`localStorage.player_current_index`

### 4. 歌词显示 (LyricDisplay.vue)

**功能特性：**
- ✅ LRC 歌词解析
- ✅ 歌词同步滚动
- ✅ 当前歌词高亮
- ✅ 点击歌词跳转到对应时间
- ✅ 手动滚动控制
- ✅ 字体大小调节
- ✅ 支持双语歌词
- ✅ 平滑滚动动画
- ✅ 响应式设计

**歌词格式支持：**
```
[00:00.00] 歌词文本
[00:05.50] 下一句歌词
```

### 5. API 服务 (src/api/music.ts)

**API 封装：**
- ✅ `searchMusic(query, page, limit)` - 音乐搜索
- ✅ `getSongDetail(songId)` - 歌曲详情
- ✅ `getSongUrl(songId)` - 获取播放 URL
- ✅ `getLyrics(songId)` - 获取歌词
- ✅ `getRecommendPlaylists(limit)` - 推荐歌单
- ✅ `getHotSongs(limit)` - 热门歌曲
- ✅ `parseLyrics(lrcText)` - LRC 歌词解析

**特性：**
- Axios 请求/响应拦截器
- 统一错误处理
- 自动添加认证 token
- 请求超时处理

### 6. 播放器 Store (src/stores/player.ts)

**状态管理：**
- ✅ 当前歌曲
- ✅ 播放列表
- ✅ 播放状态（播放/暂停）
- ✅ 播放进度
- ✅ 音量控制
- ✅ 播放模式
- ✅ 歌词数据
- ✅ 加载状态

**核心方法：**
- `playSong(song, index)` - 播放歌曲
- `togglePlay()` - 切换播放/暂停
- `nextSong()` / `prevSong()` - 上下首
- `setVolume(value)` - 设置音量
- `seekTo(time)` - 跳转进度
- `setLoopMode(mode)` - 设置播放模式
- `addToPlaylist(song)` - 添加到播放列表
- `removeFromPlaylist(index)` - 从播放列表移除
- `reorderPlaylist(from, to)` - 重排播放列表
- `clearPlaylist()` - 清空播放列表

### 7. 工具函数

**storage.ts - 本地存储工具：**
- `getFromStorage(key, default)` - 获取数据
- `saveToStorage(key, value)` - 保存数据
- `removeFromStorage(key)` - 删除数据
- `clearStorage()` - 清空存储
- `isStorageAvailable()` - 检查可用性
- `getStorageUsage()` - 获取使用量

**request.ts - HTTP 请求工具：**
- 统一的 axios 实例配置
- 请求/响应拦截器
- 错误处理
- 自动 token 注入

## 项目结构

```
web/
├── src/
│   ├── api/
│   │   └── music.ts          # 音乐 API 服务
│   ├── components/
│   │   ├── PlayerControls.vue  # 播放器控制
│   │   ├── SearchBar.vue       # 搜索栏
│   │   ├── Playlist.vue        # 播放列表
│   │   └── LyricDisplay.vue    # 歌词显示
│   ├── stores/
│   │   └── player.ts          # 播放器状态管理
│   ├── types/
│   │   └── music.ts           # 类型定义
│   ├── utils/
│   │   ├── storage.ts         # 本地存储工具
│   │   └── request.ts         # HTTP 请求工具
│   ├── router/
│   │   └── index.ts           # 路由配置
│   ├── App.vue                # 根组件
│   └── main.ts                # 入口文件
├── .env.example               # 环境变量示例
├── package.json
└── vite.config.ts
```

## 类型定义

### Song
```typescript
interface Song {
  id: string
  title: string
  artist: string
  album: string
  cover: string
  url: string
  duration: number
}
```

### LyricLine
```typescript
interface LyricLine {
  time: number      // 时间（秒）
  text: string      // 歌词文本
  translation?: string  // 翻译（可选）
}
```

## 本地存储键名

```typescript
const StorageKeys = {
  PLAYER_VOLUME: 'player_volume',
  PLAYER_PLAYLIST: 'player_playlist',
  PLAYER_CURRENT_INDEX: 'player_current_index',
  PLAYER_LOOP_MODE: 'player_loop_mode',
  SEARCH_HISTORY: 'search_history',
}
```

## 环境变量

创建 `.env` 文件（参考 `.env.example`）：

```bash
VITE_API_BASE_URL=http://localhost:8080/api
VITE_API_TIMEOUT=10000
VITE_APP_TITLE=MusicPlayer-Pro
```

## 使用方法

### 1. 安装依赖

```bash
cd web
npm install
```

### 2. 配置环境变量

```bash
cp .env.example .env
# 编辑 .env 文件配置 API 地址
```

### 3. 启动开发服务器

```bash
npm run dev
```

### 4. 构建生产版本

```bash
npm run build
```

## API 接口规范

后端 API 应遵循以下规范：

### 搜索接口
```
GET /api/music/search?query=xxx&page=1&limit=20
```

响应：
```json
{
  "songs": [...],
  "total": 100,
  "page": 1,
  "limit": 20
}
```

### 歌曲详情
```
GET /api/music/song/:id
```

### 播放 URL
```
GET /api/music/song/:id/url
```

响应：
```json
{
  "url": "https://..."
}
```

### 歌词
```
GET /api/music/song/:id/lyric
```

响应：
```json
{
  "songId": "xxx",
  "lrc": "[00:00.00] 歌词...",
  "translation": "..."
}
```

## 注意事项

1. **CORS 配置**: 确保后端 API 允许跨域请求
2. **音频格式**: 支持 HTML5 Audio 支持的格式（MP3, AAC, WAV 等）
3. **HTTPS**: 生产环境建议使用 HTTPS
4. **存储限制**: localStorage 有 5-10MB 限制，注意清理
5. **移动端适配**: 组件已做响应式设计，支持移动端

## 后续优化建议

1. 添加 Web Audio API 高级功能（均衡器、音效）
2. 实现歌词缓存机制
3. 添加离线播放功能（Service Worker）
4. 实现歌单分享功能
5. 添加用户登录和收藏功能
6. 优化移动端手势操作
7. 添加桌面通知（播放状态）
8. 实现音频可视化效果

## 技术栈

- **框架**: Vue 3.4+
- **语言**: TypeScript 5+
- **UI 组件**: Naive UI 2.36+
- **状态管理**: Pinia 2+
- **HTTP 客户端**: Axios 1+
- **构建工具**: Vite 5+
- **图标**: @vicons/ionicons5

## 许可证

MIT
