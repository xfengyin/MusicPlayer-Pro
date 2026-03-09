package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/xfengyin/MusicPlayer-Pro/server/internal/crawler"
	"github.com/xfengyin/MusicPlayer-Pro/server/internal/model"
	"gorm.io/gorm"
)

// 定义错误
var (
	ErrSongNotFound     = errors.New("song not found")
	ErrPlaylistNotFound = errors.New("playlist not found")
	ErrInvalidID        = errors.New("invalid ID")
	ErrSearchFailed     = errors.New("search failed")
)

// CacheItem 缓存项
type CacheItem struct {
	Data      interface{}
	ExpiresAt time.Time
}

// MusicService 音乐服务结构体
type MusicService struct {
	db       *gorm.DB
	crawler  *crawler.MusicCrawler
	cache    map[string]*CacheItem
	cacheMu  sync.RWMutex
	cacheTTL time.Duration
}

// MusicServiceConfig 服务配置
type MusicServiceConfig struct {
	CacheTTL time.Duration // 缓存过期时间
}

// DefaultServiceConfig 默认服务配置
func DefaultServiceConfig() *MusicServiceConfig {
	return &MusicServiceConfig{
		CacheTTL: 5 * time.Minute,
	}
}

// NewMusicService 创建新的音乐服务实例
func NewMusicService(db *gorm.DB, cfg *MusicServiceConfig) *MusicService {
	if cfg == nil {
		cfg = DefaultServiceConfig()
	}
	return &MusicService{
		db:       db,
		crawler:  nil, // 稍后通过 SetCrawler 设置
		cache:    make(map[string]*CacheItem),
		cacheTTL: cfg.CacheTTL,
	}
}

// SetCrawler 设置爬虫管理器
func (s *MusicService) SetCrawler(c *crawler.MusicCrawler) {
	s.crawler = c
}

// getCache 从缓存获取数据
func (s *MusicService) getCache(key string) (interface{}, bool) {
	s.cacheMu.RLock()
	defer s.cacheMu.RUnlock()

	item, exists := s.cache[key]
	if !exists {
		return nil, false
	}

	if time.Now().After(item.ExpiresAt) {
		return nil, false
	}

	return item.Data, true
}

// setCache 设置缓存
func (s *MusicService) setCache(key string, data interface{}) {
	s.cacheMu.Lock()
	defer s.cacheMu.Unlock()

	s.cache[key] = &CacheItem{
		Data:      data,
		ExpiresAt: time.Now().Add(s.cacheTTL),
	}
}

// SearchSongs 搜索歌曲
// 优先从缓存获取，如果没有则调用爬虫搜索并保存到数据库
func (s *MusicService) SearchSongs(query string) ([]model.Song, error) {
	// 检查缓存
	cacheKey := fmt.Sprintf("search:%s", query)
	if cached, ok := s.getCache(cacheKey); ok {
		log.Printf("[CACHE HIT] search: %s", query)
		if songs, ok := cached.([]model.Song); ok {
			return songs, nil
		}
	}

	// 先从本地数据库搜索
	localSongs, err := s.searchLocalSongs(query)
	if err == nil && len(localSongs) > 0 {
		log.Printf("[LOCAL DB] found %d songs for: %s", len(localSongs), query)
		s.setCache(cacheKey, localSongs)
		return localSongs, nil
	}

	// 本地没有，使用爬虫搜索
	if s.crawler == nil {
		log.Printf("[WARN] crawler not initialized, returning local results only")
		return localSongs, nil
	}

	log.Printf("[CRAWLER] searching: %s", query)
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	result, err := s.crawler.Search(ctx, query, 50)
	if err != nil {
		log.Printf("[ERROR] crawler search failed: %v", err)
		return localSongs, nil // 返回本地结果（可能为空）
	}

	if len(result.Songs) > 0 {
		// 保存到数据库
		s.saveSongsToDB(result.Songs)
		log.Printf("[CRAWLER] found %d songs, saved to DB", len(result.Songs))
		
		// 缓存结果
		s.setCache(cacheKey, result.Songs)
		return result.Songs, nil
	}

	return localSongs, nil
}

// searchLocalSongs 从本地数据库搜索歌曲
func (s *MusicService) searchLocalSongs(query string) ([]model.Song, error) {
	var songs []model.Song

	searchPattern := "%" + query + "%"
	err := s.db.Where("title LIKE ? OR artist LIKE ?", searchPattern, searchPattern).
		Limit(50).
		Find(&songs).Error

	if err != nil {
		return nil, fmt.Errorf("failed to search local songs: %w", err)
	}

	return songs, nil
}

// saveSongsToDB 保存歌曲到数据库
func (s *MusicService) saveSongsToDB(songs []model.Song) {
	if len(songs) == 0 {
		return
	}

	// 批量插入，忽略重复
	for _, song := range songs {
		// 检查是否已存在
		var existing model.Song
		result := s.db.Where("title = ? AND artist = ?", song.Title, song.Artist).First(&existing)
		if result.Error == nil {
			// 已存在，更新
			s.db.Model(&existing).Updates(song)
		} else if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			// 不存在，创建
			s.db.Create(&song)
		}
	}
}

// GetSongDetail 获取歌曲详情
func (s *MusicService) GetSongDetail(id string) (*model.Song, error) {
	var song model.Song

	err := s.db.First(&song, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("%w: %s", ErrSongNotFound, id)
		}
		return nil, fmt.Errorf("failed to get song detail: %w", err)
	}

	return &song, nil
}

// GetPlayUrl 获取播放链接
func (s *MusicService) GetPlayUrl(id string) (string, error) {
	song, err := s.GetSongDetail(id)
	if err != nil {
		return "", err
	}

	if song.URL == "" {
		return "", fmt.Errorf("no play url available for song: %s", id)
	}

	return song.URL, nil
}

// GetPlaylists 获取所有播放列表
func (s *MusicService) GetPlaylists() ([]model.Playlist, error) {
	var playlists []model.Playlist

	err := s.db.Preload("Songs").Find(&playlists).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get playlists: %w", err)
	}

	return playlists, nil
}

// CreatePlaylist 创建播放列表
func (s *MusicService) CreatePlaylist(name, description string, userID uint, songIDs []uint) (*model.Playlist, error) {
	playlist := model.Playlist{
		Name:        name,
		Description: description,
		UserID:      userID,
	}

	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Create(&playlist).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to create playlist: %w", err)
	}

	if len(songIDs) > 0 {
		var songs []model.Song
		if err := tx.Find(&songs, songIDs).Error; err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("failed to find songs: %w", err)
		}
		playlist.Songs = songs

		if err := tx.Save(&playlist).Error; err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("failed to associate songs: %w", err)
		}
	}

	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	var created model.Playlist
	if err := s.db.Preload("Songs").First(&created, playlist.ID).Error; err != nil {
		return nil, fmt.Errorf("failed to load created playlist: %w", err)
	}

	return &created, nil
}

// GetPlaylistByID 根据 ID 获取播放列表
func (s *MusicService) GetPlaylistByID(id uint) (*model.Playlist, error) {
	var playlist model.Playlist

	err := s.db.Preload("Songs").First(&playlist, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("%w: %d", ErrPlaylistNotFound, id)
		}
		return nil, fmt.Errorf("failed to get playlist: %w", err)
	}

	return &playlist, nil
}

// DeletePlaylist 删除播放列表
func (s *MusicService) DeletePlaylist(id uint) error {
	result := s.db.Delete(&model.Playlist{}, id)
	if result.Error != nil {
		return fmt.Errorf("failed to delete playlist: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("%w: %d", ErrPlaylistNotFound, id)
	}
	return nil
}

// AddSongToPlaylist 添加歌曲到播放列表
func (s *MusicService) AddSongToPlaylist(playlistID, songID uint) error {
	playlist := model.Playlist{ID: playlistID}
	song := model.Song{ID: songID}

	if err := s.db.Model(&playlist).Association("Songs").Append(&song); err != nil {
		return fmt.Errorf("failed to add song to playlist: %w", err)
	}
	return nil
}

// RemoveSongFromPlaylist 从播放列表移除歌曲
func (s *MusicService) RemoveSongFromPlaylist(playlistID, songID uint) error {
	playlist := model.Playlist{ID: playlistID}
	song := model.Song{ID: songID}

	if err := s.db.Model(&playlist).Association("Songs").Delete(&song); err != nil {
		return fmt.Errorf("failed to remove song from playlist: %w", err)
	}
	return nil
}

// GetHotSearches 获取热门搜索（从缓存）
func (s *MusicService) GetHotSearches() []string {
	s.cacheMu.RLock()
	defer s.cacheMu.RUnlock()

	hotSearches := make([]string, 0)
	for key := range s.cache {
		if len(key) > 7 && key[:7] == "search:" {
			hotSearches = append(hotSearches, key[7:])
		}
	}
	return hotSearches
}

// ClearCache 清除缓存
func (s *MusicService) ClearCache() {
	s.cacheMu.Lock()
	defer s.cacheMu.Unlock()
	s.cache = make(map[string]*CacheItem)
}
