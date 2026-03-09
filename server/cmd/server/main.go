package main

import (
	"log"
	"time"

	"github.com/xfengyin/MusicPlayer-Pro/server/internal/api"
	"github.com/xfengyin/MusicPlayer-Pro/server/internal/crawler"
	"github.com/xfengyin/MusicPlayer-Pro/server/internal/model"
	"github.com/xfengyin/MusicPlayer-Pro/server/internal/service"
	"github.com/xfengyin/MusicPlayer-Pro/server/pkg/config"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize database
	db, err := config.InitializeDB(cfg.Database.Path)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Initialize music service with cache
	musicService := service.NewMusicService(db, &service.MusicServiceConfig{
		CacheTTL: 5 * time.Minute,
	})

	// Initialize crawler
	crawlerConfig := crawler.DefaultConfig()
	musicCrawler := crawler.NewMusicCrawler(crawlerConfig)

	// Register music sources
	neteaseCrawler := crawler.NewNetEaseCrawler(nil)
	musicCrawler.RegisterSource(neteaseCrawler)

	// Optional: register Kugou crawler
	// kugouCrawler := crawler.NewKuGouCrawler(nil)
	// musicCrawler.RegisterSource(kugouCrawler)

	// Set crawler to service
	musicService.SetCrawler(musicCrawler)

	// Setup service container
	services := &api.MusicServices{
		Music: musicService,
	}

	// Setup Gin router
	router := gin.Default()

	// Setup API routes (includes CORS middleware)
	api.SetupRoutes(router, services)

	// Auto migrate database tables
	_ = model.Song{} // 确保 model 包被使用

	// Start server
	log.Printf("Server starting on port %s", cfg.Server.Port)
	log.Printf("Database path: %s", cfg.Database.Path)
	log.Printf("Mode: %s", cfg.Server.Mode)
	log.Printf("Registered music sources: netease")
	
	if err := router.Run(":" + cfg.Server.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
