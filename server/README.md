# MusicPlayer-Pro Server

Go 后端服务器，提供音乐搜索、播放和播放列表管理功能。

## 功能特性

### 1. 音源爬虫 (`internal/crawler/`)
- **crawler.go**: 爬虫接口定义和管理器
- **netease.go**: 网易云音乐爬虫
- **kugou.go**: 酷狗音乐爬虫（可选）

### 2. 服务层 (`internal/service/`)
- 真实音乐搜索逻辑
- 热门搜索结果缓存（TTL: 5 分钟）
- 完整的错误处理和日志
- 数据库本地搜索优先，爬虫搜索补充

### 3. 数据库操作 (`internal/model/`)
- GORM 自动迁移
- 完整的 CRUD 操作
- SQLite 数据库

### 4. API 路由 (`internal/api/`)
- RESTful API 设计
- 全局 CORS 跨域支持
- 请求验证和错误处理

## API 端点

### 音乐相关
- `GET /api/search?q=关键词` - 搜索歌曲
- `GET /api/song/:id` - 获取歌曲详情
- `GET /api/song/:id/url` - 获取播放链接
- `GET /api/hot-searches` - 获取热门搜索

### 播放列表
- `GET /api/playlist` - 获取所有播放列表
- `GET /api/playlist/:id` - 获取播放列表详情
- `POST /api/playlist` - 创建播放列表
- `DELETE /api/playlist/:id` - 删除播放列表
- `POST /api/playlist/:id/songs/:songId` - 添加歌曲到播放列表
- `DELETE /api/playlist/:id/songs/:songId` - 从播放列表移除歌曲

### 系统
- `GET /health` - 健康检查
- `POST /api/cache/clear` - 清除缓存

## 配置

通过环境变量配置：

```bash
# 服务器端口（默认：8080）
SERVER_PORT=8080

# 运行模式：debug/release/test（默认：debug）
SERVER_MODE=debug

# SQLite 数据库路径（默认：./data/music.db）
DB_PATH=./data/music.db
```

## 运行

```bash
cd server/cmd/server
go run main.go
```

## API 测试示例

### 搜索歌曲
```bash
curl "http://localhost:8080/api/search?q=周杰伦"
```

### 获取歌曲详情
```bash
curl http://localhost:8080/api/song/1
```

### 获取播放链接
```bash
curl http://localhost:8080/api/song/1/url
```

### 创建播放列表
```bash
curl -X POST http://localhost:8080/api/playlist \
  -H "Content-Type: application/json" \
  -d '{
    "name": "我的歌单",
    "description": "测试歌单",
    "user_id": 1,
    "song_ids": [1, 2, 3]
  }'
```

### 健康检查
```bash
curl http://localhost:8080/health
```

## 项目结构

```
server/
├── cmd/
│   └── server/
│       └── main.go              # 程序入口
├── internal/
│   ├── api/
│   │   └── routes.go            # API 路由和处理器
│   ├── crawler/
│   │   ├── crawler.go           # 爬虫接口
│   │   ├── netease.go           # 网易云爬虫
│   │   └── kugou.go             # 酷狗爬虫
│   ├── model/
│   │   └── song.go              # 数据模型
│   └── service/
│       └── music.go             # 业务逻辑
└── pkg/
    └── config/
        └── config.go            # 配置管理
```

## 注意事项

1. **爬虫限制**: 当前爬虫使用公开 API，可能受到反爬虫限制
2. **缓存策略**: 搜索结果缓存 5 分钟，可通过 `/api/cache/clear` 清除
3. **CORS**: 默认允许所有来源，生产环境应限制具体域名
4. **数据库**: 使用 SQLite，适合个人使用，生产环境可切换到 PostgreSQL/MySQL
