package crawler

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/xfengyin/MusicPlayer-Pro/server/internal/model"
)

// KuGouConfig 酷狗音乐配置
type KuGouConfig struct {
	BaseURL    string
	SearchPath string
	PlayPath   string
}

// DefaultKuGouConfig 默认酷狗配置
func DefaultKuGouConfig() *KuGouConfig {
	return &KuGouConfig{
		BaseURL:    "https://mobileservice.kugou.com",
		SearchPath: "/api/v3/search/song",
		PlayPath:   "/api/v3/url/song",
	}
}

// KuGouCrawler 酷狗音乐爬虫
type KuGouCrawler struct {
	config  *KuGouConfig
	client  *http.Client
	timeout time.Duration
}

// kugouSearchResponse 酷狗搜索响应
type kugouSearchResponse struct {
	Data struct {
		Info []struct {
			SongID     int    `json:"song_id"`
			SongName   string `json:"song_name"`
			AuthorName string `json:"author_name"`
			AlbumName  string `json:"album_name"`
			Duration   int    `json:"duration"`
			Img        string `json:"img"`
			Hash       string `json:"hash"`
		} `json:"info"`
		Total int `json:"total"`
	} `json:"data"`
}

// NewKuGouCrawler 创建酷狗爬虫
func NewKuGouCrawler(config *KuGouConfig) *KuGouCrawler {
	if config == nil {
		config = DefaultKuGouConfig()
	}
	return &KuGouCrawler{
		config:  config,
		client:  &http.Client{Timeout: 10 * time.Second},
		timeout: 10 * time.Second,
	}
}

// GetName 获取音源名称
func (kc *KuGouCrawler) GetName() string {
	return "kugou"
}

// Search 搜索歌曲
func (kc *KuGouCrawler) Search(ctx context.Context, query string, limit int) (*SearchResult, error) {
	// 构建搜索 URL
	searchURL := fmt.Sprintf("%s%s", kc.config.BaseURL, kc.config.SearchPath)
	params := url.Values{}
	params.Set("keyword", query)
	params.Set("page", "1")
	params.Set("pagesize", strconv.Itoa(limit))
	params.Set("showtype", "1")

	fullURL := searchURL + "?" + params.Encode()

	// 创建请求
	req, err := http.NewRequestWithContext(ctx, "GET", fullURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// 设置请求头
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")

	// 发送请求
	resp, err := kc.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	// 解析响应
	var searchResp kugouSearchResponse
	if err := json.Unmarshal(body, &searchResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	// 转换为通用模型
	songs := make([]model.Song, 0, len(searchResp.Data.Info))
	for _, song := range searchResp.Data.Info {
		songs = append(songs, model.Song{
			Title:    song.SongName,
			Artist:   song.AuthorName,
			Album:    song.AlbumName,
			Duration: song.Duration,
			CoverURL: song.Img,
			URL:      fmt.Sprintf("https://www.kugou.com/song/%s.html", song.Hash),
		})
	}

	return &SearchResult{
		Songs: songs,
		Total: searchResp.Data.Total,
	}, nil
}

// GetSongDetail 获取歌曲详情
func (kc *KuGouCrawler) GetSongDetail(ctx context.Context, songID string) (*model.Song, error) {
	// 注意：由于版权原因，我们不能直接提供真实音乐数据
	// 这里提供一个框架，实际部署时应使用官方API或合法途径
	return nil, fmt.Errorf("酷狗音乐详情获取功能暂未启用，因版权原因")
}

// GetPlayURL 获取播放 URL
func (kc *KuGouCrawler) GetPlayURL(ctx context.Context, songID string) (string, error) {
	// 注意：由于版权原因，我们不能直接提供真实音乐播放链接
	// 这里提供一个框架，实际部署时应使用官方API或合法途径
	return "", fmt.Errorf("酷狗音乐播放链接获取功能暂未启用，因版权原因")
}
