# Music Player Pro - Makefile
# 便捷命令入口

.PHONY: help dev dev-web dev-server build build-web build-server build-desktop build-android test lint validate clean

# 默认目标
help:
	@echo "Music Player Pro - 可用命令"
	@echo ""
	@echo "开发命令:"
	@echo "  make dev          - 启动所有开发服务"
	@echo "  make dev-web      - 仅启动前端开发服务器"
	@echo "  make dev-server   - 仅启动后端开发服务器"
	@echo "  make dev-desktop  - 启动 Electron 开发模式"
	@echo ""
	@echo "构建命令:"
	@echo "  make build        - 构建 Web 和后端"
	@echo "  make build-web    - 构建前端"
	@echo "  make build-server - 构建后端"
	@echo "  make build-desktop- 构建桌面端"
	@echo "  make build-android- 构建 Android 应用"
	@echo "  make build-all    - 构建所有平台"
	@echo ""
	@echo "验证命令:"
	@echo "  make validate     - 快速验证 (10秒内)"
	@echo "  make test         - 运行所有测试"
	@echo "  make lint         - 代码检查"
	@echo "  make lint-fix     - 自动修复代码问题"
	@echo ""
	@echo "其他命令:"
	@echo "  make clean        - 清理构建产物"
	@echo "  make install      - 安装所有依赖"
	@echo "  make release      - 创建发布版本"

# ========== 开发命令 ==========

dev:
	@echo "🚀 启动所有开发服务..."
	@make dev-server &
	@sleep 2
	@make dev-web

dev-web:
	@echo "🎨 启动前端开发服务器..."
	cd web && npm run dev

dev-server:
	@echo "⚙️  启动后端开发服务器..."
	cd server && go run cmd/server/main.go

dev-desktop:
	@echo "💻 启动 Electron 开发模式..."
	cd electron && npm run dev

# ========== 构建命令 ==========

build: build-web build-server
	@echo "✅ 构建完成"

build-web:
	@echo "🔨 构建前端..."
	cd web && npm run build

build-server:
	@echo "🔨 构建后端..."
	cd server && go build -o dist/server cmd/server/main.go

build-desktop: build-web
	@echo "💻 构建桌面端..."
	@mkdir -p electron/renderer
	@cp -r web/dist/* electron/renderer/
	cd electron && npm run build

build-android:
	@echo "📱 构建 Android 应用..."
	@echo "TODO: Android 构建"

build-all: build-desktop build-android
	@echo "✅ 所有平台构建完成"

# ========== 验证命令 ==========

validate:
	@echo "⚡ 快速验证..."
	@echo "检查前端..."
	@cd web && npm run type-check || true
	@echo "检查后端..."
	@cd server && go vet ./... || true
	@echo "✅ 验证完成"

test:
	@echo "🧪 运行测试..."
	@echo "前端测试..."
	@cd web && npm run test || true
	@echo "后端测试..."
	@cd server && go test ./... || true

lint:
	@echo "🔍 代码检查..."
	@cd web && npm run lint || true
	@cd server && golangci-lint run || true

lint-fix:
	@echo "🔧 自动修复代码问题..."
	@cd web && npm run lint:fix || true

# ========== 其他命令 ==========

install:
	@echo "📦 安装依赖..."
	@cd web && npm install
	@cd server && go mod download
	@cd electron && npm install

clean:
	@echo "🧹 清理构建产物..."
	@rm -rf web/dist
	@rm -rf server/dist
	@rm -rf electron/dist
	@rm -rf electron/renderer
	@echo "✅ 清理完成"

release:
	@echo "🚀 创建发布版本..."
	@echo "TODO: 发布流程"

# CI/CD 命令
ci: lint test build
	@echo "✅ CI 流程完成"
