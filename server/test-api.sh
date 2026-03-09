#!/bin/bash

# MusicPlayer-Pro API 测试脚本
# 使用前请确保服务器已启动：go run cmd/server/main.go

BASE_URL="http://localhost:8080"

echo "=========================================="
echo "MusicPlayer-Pro API 测试"
echo "=========================================="
echo ""

# 1. 健康检查
echo "1. 健康检查..."
curl -s "$BASE_URL/health" | jq .
echo ""

# 2. 搜索歌曲
echo "2. 搜索歌曲：周杰伦..."
curl -s "$BASE_URL/api/search?q=周杰伦" | jq .
echo ""

# 3. 搜索歌曲
echo "3. 搜索歌曲：Taylor Swift..."
curl -s "$BASE_URL/api/search?q=Taylor%20Swift" | jq .
echo ""

# 4. 获取歌曲详情（假设有 ID 为 1 的歌曲）
echo "4. 获取歌曲详情（ID: 1）..."
curl -s "$BASE_URL/api/song/1" | jq .
echo ""

# 5. 获取播放链接
echo "5. 获取播放链接（ID: 1）..."
curl -s "$BASE_URL/api/song/1/url" | jq .
echo ""

# 6. 获取所有播放列表
echo "6. 获取所有播放列表..."
curl -s "$BASE_URL/api/playlist" | jq .
echo ""

# 7. 创建播放列表
echo "7. 创建播放列表..."
curl -s -X POST "$BASE_URL/api/playlist" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "测试歌单",
    "description": "这是一个测试播放列表",
    "user_id": 1,
    "song_ids": []
  }' | jq .
echo ""

# 8. 获取热门搜索
echo "8. 获取热门搜索..."
curl -s "$BASE_URL/api/hot-searches" | jq .
echo ""

# 9. 清除缓存
echo "9. 清除缓存..."
curl -s -X POST "$BASE_URL/api/cache/clear" | jq .
echo ""

# 10. 测试 CORS
echo "10. 测试 CORS 预检请求..."
curl -s -X OPTIONS "$BASE_URL/api/search" \
  -H "Origin: http://localhost:3000" \
  -H "Access-Control-Request-Method: GET" \
  -i | grep -E "HTTP|Access-Control"
echo ""

echo "=========================================="
echo "测试完成!"
echo "=========================================="
