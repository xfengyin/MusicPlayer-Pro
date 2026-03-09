# MusicPlayer-Pro Server 快速开始

## 前置要求

- Go 1.21+
- Git

## 安装步骤

### 1. 安装依赖

```bash
cd server
go mod tidy
```

### 2. 配置环境变量（可选）

```bash
# 创建 .env 文件或导出环境变量
export SERVER_PORT=8080
export SERVER_MODE=debug
export DB_PATH=./data/music.db
```

### 3. 启动服务器

```bash
go run cmd/server/main.go
```

启动成功后会看到：
```
Server starting on port 8080
Database path: ./data/music.db
Mode: debug
Registered music sources: netease
```

## 快速测试

### 使用 curl 测试

```bash
# 健康检查
curl http://localhost:8080/health

# 搜索歌曲
curl "http://localhost:8080/api/search?q=周杰伦"

# 获取热门搜索
curl http://localhost:8080/api/hot-searches
```

### 使用测试脚本

```bash
./test-api.sh
```

## API 文档

### 搜索歌曲

**请求**
```http
GET /api/search?q=关键词
```

**响应**
```json
{
  "data": {
    "songs": [
      {
        "id": 1,
        "title": "歌曲名",
        "artist": "艺术家",
        "album": "专辑",
        "duration": 240,
        "url": "https://...",
        "cover_url": "https://..."
      }
    ],
    "total": 10,
    "query": "周杰伦",
    "duration_ms": 150
  }
}
```

### 创建播放列表

**请求**
```http
POST /api/playlist
Content-Type: application/json

{
  "name": "我的歌单",
  "description": "喜欢的歌曲",
  "user_id": 1,
  "song_ids": [1, 2, 3]
}
```

**响应**
```json
{
  "data": {
    "id": 1,
    "name": "我的歌单",
    "description": "喜欢的歌曲",
    "user_id": 1,
    "songs": [...],
    "created_at": "2024-01-01T00:00:00Z"
  }
}
```

## 前端集成示例

### Vue3 + Axios

```javascript
// api/music.js
import axios from 'axios'

const api = axios.create({
  baseURL: 'http://localhost:8080/api',
  timeout: 10000
})

export const musicAPI = {
  // 搜索歌曲
  search(query) {
    return api.get('/search', { params: { q: query } })
  },
  
  // 获取歌曲详情
  getSong(id) {
    return api.get(`/song/${id}`)
  },
  
  // 获取播放链接
  getPlayUrl(id) {
    return api.get(`/song/${id}/url`)
  },
  
  // 获取播放列表
  getPlaylists() {
    return api.get('/playlist')
  },
  
  // 创建播放列表
  createPlaylist(data) {
    return api.post('/playlist', data)
  }
}
```

### 使用示例

```vue
<script setup>
import { ref } from 'vue'
import { musicAPI } from './api/music'

const searchQuery = ref('')
const searchResults = ref([])

const handleSearch = async () => {
  try {
    const response = await musicAPI.search(searchQuery.value)
    searchResults.value = response.data.data.songs
  } catch (error) {
    console.error('搜索失败:', error)
  }
}
</script>

<template>
  <div>
    <input v-model="searchQuery" @keyup.enter="handleSearch" />
    <button @click="handleSearch">搜索</button>
    
    <div v-for="song in searchResults" :key="song.id">
      <h3>{{ song.title }}</h3>
      <p>{{ song.artist }}</p>
    </div>
  </div>
</template>
```

## 常见问题

### 1. 端口被占用

**错误**: `bind: address already in use`

**解决**: 修改端口
```bash
export SERVER_PORT=8081
```

### 2. CORS 错误

确保后端已启动，并且前端请求的 URL 正确。后端默认允许所有来源的 CORS 请求。

### 3. 数据库文件不存在

数据库文件会在首次运行时自动创建。确保 `./data/` 目录有写入权限。

### 4. 爬虫搜索失败

- 检查网络连接
- 查看日志输出
- 爬虫可能会失败，系统会自动降级到本地数据库搜索

## 开发模式

### 热重载

使用 `air` 实现热重载：

```bash
go install github.com/air-verse/air@latest
air
```

### 调试模式

设置 `SERVER_MODE=debug` 启用详细日志。

## 生产部署

### 编译

```bash
go build -o music-server cmd/server/main.go
```

### 运行

```bash
./music-server
```

### Docker (可选)

```dockerfile
FROM golang:1.21-alpine
WORKDIR /app
COPY . .
RUN go build -o server cmd/server/main.go
EXPOSE 8080
CMD ["./server"]
```

## 下一步

1. 阅读 [README.md](README.md) 了解完整功能
2. 查看 [IMPLEMENTATION_REPORT.md](../IMPLEMENTATION_REPORT.md) 了解实现细节
3. 开始集成前端
