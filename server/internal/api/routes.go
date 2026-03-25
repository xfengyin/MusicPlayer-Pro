package api

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/xfengyin/MusicPlayer-Pro/server/internal/service"
)

// MusicServices 包含所有音乐相关的服务
type MusicServices struct {
	Music *service.MusicService
}

// SetupRoutes 设置所有 API 路由
func SetupRoutes(router *gin.Engine, services *MusicServices) {
	// 全局 CORS 中间件
	router.Use(CORSMiddleware())

	// 创建 API 路由组 - 保持与前端一致的路径
	api := router.Group("/api")
	{
		// 搜索歌曲
		api.GET("/search", SearchSongsHandler(services.Music))
		// 为前端提供/music/路径兼容
		router.GET("/music/search", SearchSongsHandler(services.Music))

		// 获取歌曲详情
		api.GET("/song/:id", GetSongDetailHandler(services.Music))
		// 为前端提供/music/路径兼容
		router.GET("/music/song/:id", GetSongDetailHandler(services.Music))

		// 获取播放链接
		api.GET("/song/:id/url", GetPlayUrlHandler(services.Music))
		// 为前端提供/music/路径兼容
		router.GET("/music/song/:id/url", GetPlayUrlHandler(services.Music))

		// 获取歌词
		api.GET("/song/:id/lyric", GetLyricsHandler(services.Music))
		// 为前端提供/music/路径兼容
		router.GET("/music/song/:id/lyric", GetLyricsHandler(services.Music))

		// 获取播放列表
		api.GET("/playlist", GetPlaylistsHandler(services.Music))
		api.GET("/playlist/:id", GetPlaylistByIDHandler(services.Music))

		// 创建播放列表
		api.POST("/playlist", CreatePlaylistHandler(services.Music))

		// 删除播放列表
		api.DELETE("/playlist/:id", DeletePlaylistHandler(services.Music))

		// 添加歌曲到播放列表
		api.POST("/playlist/:id/songs/:songId", AddSongToPlaylistHandler(services.Music))

		// 从播放列表移除歌曲
		api.DELETE("/playlist/:id/songs/:songId", RemoveSongFromPlaylistHandler(services.Music))

		// 热门搜索
		api.GET("/hot-searches", GetHotSearchesHandler(services.Music))

		// 清除缓存
		api.POST("/cache/clear", ClearCacheHandler(services.Music))
	}

	// 健康检查
	router.GET("/health", HealthHandler())
}

// CORSMiddleware CORS 跨域中间件
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 生产环境应限制为具体的可信域名
		origin := c.Request.Header.Get("Origin")
		allowedOrigins := []string{
			"http://localhost:3000",
			"http://localhost:5173", 
			"http://127.0.0.1:5173",
			"http://localhost:8080",
			"https://localhost",       // Capacitor Android/iOS
			"capacitor://localhost",   // Capacitor iOS
			"http://localhost",        // Capacitor fallback
			// 在生产环境中，添加您的实际域名
			// "https://yourdomain.com",
			// "https://www.yourdomain.com",
		}
		
		// 检查来源是否在允许列表中
		isAllowed := false
		if origin == "" {
			isAllowed = true // 没有来源头，通常表示同源请求
		} else {
			for _, allowed := range allowedOrigins {
				if origin == allowed {
					isAllowed = true
					break
				}
			}
		}
		
		if isAllowed {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
		}
		
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		// 处理预检请求
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

// SearchSongsHandler 搜索歌曲处理函数
func SearchSongsHandler(musicService *service.MusicService) gin.HandlerFunc {
	return func(c *gin.Context) {
		query := c.Query("q")
		if query == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "query parameter 'q' is required",
			})
			return
		}

		startTime := time.Now()
		songs, err := musicService.SearchSongs(query)
		duration := time.Since(startTime)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"data": gin.H{
				"songs":      songs,
				"total":      len(songs),
				"query":      query,
				"duration_ms": duration.Milliseconds(),
			},
		})
	}
}

// GetSongDetailHandler 获取歌曲详情处理函数
func GetSongDetailHandler(musicService *service.MusicService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if id == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "song id is required",
			})
			return
		}

		song, err := musicService.GetSongDetail(id)
		if err != nil {
			// 使用更准确的错误检查
			if err.Error() == "record not found" || 
			   err.Error() == "song not found" || 
			   err.Error()[:13] == "song not found" {
				c.JSON(http.StatusNotFound, gin.H{
					"error": "song not found",
				})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"data": song,
		})
	}
}

// GetPlayUrlHandler 获取播放链接处理函数
func GetPlayUrlHandler(musicService *service.MusicService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if id == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "song id is required",
			})
			return
		}

		url, err := musicService.GetPlayUrl(id)
		if err != nil {
			statusCode := http.StatusInternalServerError
			if err.Error()[:13] == "song not found" || err.Error()[:24] == "no play url available for" {
				statusCode = http.StatusNotFound
			}
			c.JSON(statusCode, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"data": gin.H{
				"url": url,
			},
		})
	}
}

// GetPlaylistsHandler 获取播放列表处理函数
func GetPlaylistsHandler(musicService *service.MusicService) gin.HandlerFunc {
	return func(c *gin.Context) {
		playlists, err := musicService.GetPlaylists()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"data": playlists,
		})
	}
}

// CreatePlaylistRequest 创建播放列表请求体
type CreatePlaylistRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	CoverURL    string `json:"cover_url"`
	UserID      uint   `json:"user_id" binding:"required"`
	SongIDs     []uint `json:"song_ids"`
}

// CreatePlaylistHandler 创建播放列表处理函数
func CreatePlaylistHandler(musicService *service.MusicService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req CreatePlaylistRequest

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid request body: " + err.Error(),
			})
			return
		}

		playlist, err := musicService.CreatePlaylist(
			req.Name,
			req.Description,
			req.UserID,
			req.SongIDs,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"data": playlist,
		})
	}
}

// GetPlaylistByIDHandler 根据 ID 获取播放列表处理函数
func GetPlaylistByIDHandler(musicService *service.MusicService) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.ParseUint(idStr, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid playlist id",
			})
			return
		}

		playlist, err := musicService.GetPlaylistByID(uint(id))
		if err != nil {
			statusCode := http.StatusInternalServerError
			if err.Error()[:18] == "playlist not found" {
				statusCode = http.StatusNotFound
			}
			c.JSON(statusCode, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"data": playlist,
		})
	}
}

// DeletePlaylistHandler 删除播放列表处理函数
func DeletePlaylistHandler(musicService *service.MusicService) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.ParseUint(idStr, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid playlist id",
			})
			return
		}

		if err := musicService.DeletePlaylist(uint(id)); err != nil {
			statusCode := http.StatusInternalServerError
			if err.Error()[:18] == "playlist not found" {
				statusCode = http.StatusNotFound
			}
			c.JSON(statusCode, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "playlist deleted successfully",
		})
	}
}

// AddSongToPlaylistHandler 添加歌曲到播放列表处理函数
func AddSongToPlaylistHandler(musicService *service.MusicService) gin.HandlerFunc {
	return func(c *gin.Context) {
		playlistIDStr := c.Param("id")
		songIDStr := c.Param("songId")

		playlistID, err := strconv.ParseUint(playlistIDStr, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid playlist id",
			})
			return
		}

		songID, err := strconv.ParseUint(songIDStr, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid song id",
			})
			return
		}

		if err := musicService.AddSongToPlaylist(uint(playlistID), uint(songID)); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "song added to playlist successfully",
		})
	}
}

// RemoveSongFromPlaylistHandler 从播放列表移除歌曲处理函数
func RemoveSongFromPlaylistHandler(musicService *service.MusicService) gin.HandlerFunc {
	return func(c *gin.Context) {
		playlistIDStr := c.Param("id")
		songIDStr := c.Param("songId")

		playlistID, err := strconv.ParseUint(playlistIDStr, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid playlist id",
			})
			return
		}

		songID, err := strconv.ParseUint(songIDStr, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid song id",
			})
			return
		}

		if err := musicService.RemoveSongFromPlaylist(uint(playlistID), uint(songID)); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "song removed from playlist successfully",
		})
	}
}

// GetLyricsHandler 获取歌词处理函数
func GetLyricsHandler(musicService *service.MusicService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if id == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "song id is required",
			})
			return
		}

		// 目前歌词功能暂未完全实现，返回空歌词
		// TODO: 实现歌词获取逻辑
		c.JSON(http.StatusOK, gin.H{
			"data": gin.H{
				"lyric": "",
				"tlyric": "", // 翻译歌词
			},
		})
	}
}

// GetHotSearchesHandler 获取热门搜索处理函数
func GetHotSearchesHandler(musicService *service.MusicService) gin.HandlerFunc {
	return func(c *gin.Context) {
		hotSearches := musicService.GetHotSearches()
		c.JSON(http.StatusOK, gin.H{
			"data": gin.H{
				"hot_searches": hotSearches,
			},
		})
	}
}

// ClearCacheHandler 清除缓存处理函数
func ClearCacheHandler(musicService *service.MusicService) gin.HandlerFunc {
	return func(c *gin.Context) {
		musicService.ClearCache()
		c.JSON(http.StatusOK, gin.H{
			"message": "cache cleared successfully",
		})
	}
}

// HealthHandler 健康检查处理函数
func HealthHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":    "healthy",
			"timestamp": time.Now().Format(time.RFC3339),
		})
	}
}
