package crawler

import (
	"context"
	"errors"
	"time"

	"github.com/xfengyin/MusicPlayer-Pro/server/internal/model"
)

// 定义错误
var (
	ErrSourceNotAvailable = errors.New("music source not available")
	ErrSearchFailed       = errors.New("search failed")
	ErrParseFailed        = errors.New("failed to parse response")
	ErrRequestTimeout     = errors.New("request timeout")
)

// SearchResult 搜索结果
type SearchResult struct {
	Songs []model.Song `json:"songs"`
	Total int          `json:"total"`
}

// MusicSource 音乐源接口
type MusicSource interface {
	// Search 搜索歌曲
	Search(ctx context.Context, query string, limit int) (*SearchResult, error)
	
	// GetSongDetail 获取歌曲详情
	GetSongDetail(ctx context.Context, songID string) (*model.Song, error)
	
	// GetPlayURL 获取播放 URL
	GetPlayURL(ctx context.Context, songID string) (string, error)
	
	// GetName 获取音源名称
	GetName() string
}

// CrawlerConfig 爬虫配置
type CrawlerConfig struct {
	Timeout time.Duration // 请求超时时间
	Retry   int           // 重试次数
}

// DefaultConfig 默认配置
func DefaultConfig() *CrawlerConfig {
	return &CrawlerConfig{
		Timeout: 10 * time.Second,
		Retry:   2,
	}
}

// MusicCrawler 音乐爬虫管理器
type MusicCrawler struct {
	sources []MusicSource
	config  *CrawlerConfig
}

// NewMusicCrawler 创建新的爬虫管理器
func NewMusicCrawler(config *CrawlerConfig) *MusicCrawler {
	if config == nil {
		config = DefaultConfig()
	}
	return &MusicCrawler{
		sources: make([]MusicSource, 0),
		config:  config,
	}
}

// RegisterSource 注册音乐源
func (mc *MusicCrawler) RegisterSource(source MusicSource) {
	mc.sources = append(mc.sources, source)
}

// Search 搜索歌曲（从所有音源）
func (mc *MusicCrawler) Search(ctx context.Context, query string, limit int) (*SearchResult, error) {
	if len(mc.sources) == 0 {
		return nil, ErrSourceNotAvailable
	}

	// 从第一个可用音源搜索
	for _, source := range mc.sources {
		result, err := source.Search(ctx, query, limit)
		if err == nil && result != nil && len(result.Songs) > 0 {
			return result, nil
		}
	}

	return nil, ErrSearchFailed
}

// GetSongDetail 获取歌曲详情
func (mc *MusicCrawler) GetSongDetail(ctx context.Context, sourceName, songID string) (*model.Song, error) {
	for _, source := range mc.sources {
		if source.GetName() == sourceName {
			return source.GetSongDetail(ctx, songID)
		}
	}
	return nil, ErrSourceNotAvailable
}

// GetPlayURL 获取播放 URL
func (mc *MusicCrawler) GetPlayURL(ctx context.Context, sourceName, songID string) (string, error) {
	for _, source := range mc.sources {
		if source.GetName() == sourceName {
			return source.GetPlayURL(ctx, songID)
		}
	}
	return "", ErrSourceNotAvailable
}
