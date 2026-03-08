# Music Player Pro - 跨平台音乐播放器

> 🎵 现代化、跨平台、多音源音乐播放器

## 📋 项目信息

- **项目名称**: Music Player Pro
- **版本**: v1.0.0
- **平台**: Windows / macOS / Linux / Android / Web
- **技术栈**: Vue3 + TypeScript + Vite + Go + Gin + Electron

## 🎯 核心功能

### 前端功能
- 🎨 **现代化界面** - Naive UI 组件库，暗黑/明亮模式
- 🎵 **播放器控制** - 播放/暂停/上一首/下一首/进度调节
- 📝 **歌词同步** - 实时歌词滚动，逐字高亮
- 🔍 **智能搜索** - 多音源聚合搜索，历史记录
- 📋 **播放列表** - 创建/编辑/管理播放列表
- 📱 **响应式布局** - 适配桌面端和移动端
- 💻 **桌面端封装** - Electron 打包，系统托盘

### 后端功能
- 🔌 **多音源聚合** - 整合多个音乐平台API
- 🔍 **歌曲解析** - 智能解析歌曲信息
- ▶️ **播放链接** - 获取高质量播放链接
- ⬇️ **下载管理** - 多线程下载，断点续传
- 📡 **API接口** - RESTful API，文档完善

### 构建与发布
- 🔧 **自动化构建** - GitHub Actions CI/CD
- 📦 **多平台打包** - Windows/macOS/Linux/Android
- 🐳 **Docker支持** - 容器化部署
- 🚀 **自动发版** - 版本管理，Release自动化

### 团队协作
- 🌿 **Git分支管理** - Git Flow 工作流
- 📝 **Commit规范** - Conventional Commits
- 📚 **文档体系** - API文档、开发指南、贡献规范
- ✅ **代码审核** - PR Review，质量保障

## 🏗️ 项目架构

```
MusicPlayer-Pro/
├── web/                    # Vue3 前端项目
│   ├── src/
│   │   ├── components/     # 组件
│   │   ├── views/          # 页面
│   │   ├── stores/         # Pinia状态管理
│   │   ├── api/            # API接口
│   │   └── utils/          # 工具函数
│   ├── public/             # 静态资源
│   └── package.json        # 前端依赖
│
├── server/                 # Go 后端服务
│   ├── cmd/                # 入口文件
│   ├── internal/           # 内部模块
│   │   ├── api/            # API路由
│   │   ├── service/        # 业务逻辑
│   │   ├── model/          # 数据模型
│   │   ├── crawler/        # 音源解析
│   │   └── utils/          # 工具函数
│   ├── pkg/                # 公共包
│   ├── go.mod              # Go模块
│   └── Dockerfile          # Docker配置
│
├── electron/               # Electron 桌面端
│   ├── main/               # 主进程
│   ├── renderer/           # 渲染进程(引用web构建产物)
│   └── package.json        # Electron配置
│
├── docs/                   # 项目文档
├── tests/                  # 测试文件
├── scripts/                # 构建脚本
├── .github/workflows/      # CI/CD配置
└── Makefile                # 便捷命令
```

## 🚀 快速开始

### 环境要求
- Node.js 18+
- Go 1.21+
- SQLite 3

### 安装依赖
```bash
# 安装前端依赖
cd web && npm install

# 安装后端依赖
cd server && go mod download

# 安装 Electron 依赖
cd electron && npm install
```

### 开发模式
```bash
# 使用 Makefile
make dev          # 启动所有服务
make dev-web      # 仅启动前端
make dev-server   # 仅启动后端

# 或手动启动
# 终端1: 启动后端
cd server && go run cmd/server/main.go

# 终端2: 启动前端
cd web && npm run dev

# 终端3: 启动 Electron
cd electron && npm run dev
```

### 构建打包
```bash
# 构建 Web
make build-web

# 构建后端
make build-server

# 构建桌面端
make build-desktop

# 构建 Android
make build-android

# 构建所有
make build-all
```

### 本地验证
```bash
make validate     # 快速验证
make test         # 运行测试
make lint         # 代码检查
```

## 👥 开发团队

| 角色 | 职责 | 成员 |
|------|------|------|
| 技术总监 | 架构设计、API设计、代码审查 | TD-02 |
| 前端开发 ×2 | Vue3播放器、Electron封装 | FE-03, FE-04 |
| 后端开发 ×2 | Go服务、多音源解析 | BE-03, BE-04 |
| 测试工程师 | 测试用例、自动化测试 | QA-02 |
| DevOps工程师 | CI/CD、多平台构建 | DO-02 |

## 📚 文档

- [API文档](docs/API.md)
- [开发指南](docs/DEVELOPMENT.md)
- [贡献规范](docs/CONTRIBUTING.md)
- [架构设计](docs/ARCHITECTURE.md)

## 📄 许可证

MIT License
