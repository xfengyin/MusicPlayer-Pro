package config

import (
	"fmt"
	"os"

	"github.com/xfengyin/MusicPlayer-Pro/server/internal/model"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Config 服务器配置结构体
type Config struct {
	Server   ServerConfig   // 服务器配置
	Database DatabaseConfig // 数据库配置
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Port string // 服务器端口
	Mode string // 运行模式 (debug/release/test)
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Path string // SQLite 数据库文件路径
}

// Load 从环境变量加载配置
// 配置项说明:
// - SERVER_PORT: 服务器端口 (默认：8080)
// - SERVER_MODE: 运行模式 (默认：debug)
// - DB_PATH: SQLite 数据库路径 (默认：./data/music.db)
func Load() (*Config, error) {
	config := &Config{}

	// 服务器配置
	config.Server.Port = getEnv("SERVER_PORT", "8080")
	config.Server.Mode = getEnv("SERVER_MODE", "debug")

	// 数据库配置
	config.Database.Path = getEnv("DB_PATH", "./data/music.db")

	return config, nil
}

// getEnv 获取环境变量，如果不存在则返回默认值
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// InitializeDB 初始化数据库连接
// 自动创建数据库文件和表结构
func InitializeDB(dbPath string) (*gorm.DB, error) {
	// 确保数据库目录存在
	dbDir := "./data"
	if err := os.MkdirAll(dbDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create database directory: %w", err)
	}

	// 打开 SQLite 数据库连接
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// 自动迁移数据库表结构
	if err := autoMigrate(db); err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	return db, nil
}

// autoMigrate 自动迁移数据库表
func autoMigrate(db *gorm.DB) error {
	// 执行自动迁移
	err := db.AutoMigrate(
		&model.Song{},
		&model.Playlist{},
		&model.PlaylistSong{},
	)
	if err != nil {
		return err
	}

	return nil
}

// SeedData 初始化示例数据（可选）
func SeedData(db *gorm.DB) error {
	// 检查是否已有数据
	var count int64
	db.Model(&model.Song{}).Count(&count)
	if count > 0 {
		return nil // 已有数据，跳过 seeding
	}

	// 示例歌曲数据
	sampleSongs := []model.Song{
		{
			Title:    "示例歌曲 1",
			Artist:   "艺术家 A",
			Album:    "专辑 X",
			Duration: 240,
			URL:      "https://example.com/song1.mp3",
			CoverURL: "https://example.com/cover1.jpg",
		},
		{
			Title:    "示例歌曲 2",
			Artist:   "艺术家 B",
			Album:    "专辑 Y",
			Duration: 180,
			URL:      "https://example.com/song2.mp3",
			CoverURL: "https://example.com/cover2.jpg",
		},
	}

	for _, song := range sampleSongs {
		if err := db.Create(&song).Error; err != nil {
			return fmt.Errorf("failed to seed song: %w", err)
		}
	}

	return nil
}
