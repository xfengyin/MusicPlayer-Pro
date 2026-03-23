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

// NetEaseConfig 网易云音乐配置
type NetEaseConfig struct {
	BaseURL    string
	APIPath    string
	SearchPath string
}

// DefaultNetEaseConfig 默认网易云配置
func DefaultNetEaseConfig() *NetEaseConfig {
	return &NetEaseConfig{
		BaseURL:    "https://music.163.com",
		APIPath:    "/api",
		SearchPath: "/search/get",
	}
}

// NetEaseCrawler 网易云音乐爬虫
type NetEaseCrawler struct {
	config  *NetEaseConfig
	client  *http.Client
	timeout time.Duration
}

// neteaseSearchResponse 网易云搜索响应
type neteaseSearchResponse struct {
	Result struct {
		Songs []struct {
			ID        int    `json:"id"`
			Name      string `json:"name"`
			Artists   []struct {
				Name string `json:"name"`
			} `json:"artists"`
			Album struct {
				Name string `json:"name"`
			} `json:"album"`
			Duration int    `json:"duration"`
			PicURL   string `json:"picUrl"`
		} `json:"songs"`
		SongCount int `json:"songCount"`
	} `json:"result"`
}

// NewNetEaseCrawler 创建网易云爬虫
func NewNetEaseCrawler(config *NetEaseConfig) *NetEaseCrawler {
	if config == nil {
		config = DefaultNetEaseConfig()
	}
	return &NetEaseCrawler{
		config:  config,
		client:  &http.Client{Timeout: 10 * time.Second},
		timeout: 10 * time.Second,
	}
}

// GetName 获取音源名称
func (nc *NetEaseCrawler) GetName() string {
	return "netease"
}

// Search 搜索歌曲
func (nc *NetEaseCrawler) Search(ctx context.Context, query string, limit int) (*SearchResult, error) {
	// 构建搜索 URL
	searchURL := fmt.Sprintf("%s%s%s", nc.config.BaseURL, nc.config.APIPath, nc.config.SearchPath)
	params := url.Values{}
	params.Set("s", query)
	params.Set("type", "1") // 1=单曲
	params.Set("limit", strconv.Itoa(limit))
	params.Set("offset", "0")

	fullURL := searchURL + "?" + params.Encode()

	// 创建请求
	req, err := http.NewRequestWithContext(ctx, "GET", fullURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// 设置请求头（模拟浏览器）
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	req.Header.Set("Referer", "https://music.163.com/")

	// 发送请求
	resp, err := nc.client.Do(req)
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
	var searchResp neteaseSearchResponse
	if err := json.Unmarshal(body, &searchResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	// 转换为通用模型
	songs := make([]model.Song, 0, len(searchResp.Result.Songs))
	for _, song := range searchResp.Result.Songs {
		artist := ""
		if len(song.Artists) > 0 {
			artist = song.Artists[0].Name
		}

		songs = append(songs, model.Song{
			Title:    song.Name,
			Artist:   artist,
			Album:    song.Album.Name,
			Duration: song.Duration / 1000, // 毫秒转秒
			CoverURL: song.PicURL,
			URL:      fmt.Sprintf("https://music.163.com/song/media/outer/url?id=%d.mp3", song.ID),
		})
	}

	return &SearchResult{
		Songs: songs,
		Total: searchResp.Result.SongCount,
	}, nil
}

// GetSongDetail 获取歌曲详情
func (nc *NetEaseCrawler) GetSongDetail(ctx context.Context, songID string) (*model.Song, error) {
	// 注意：由于版权原因，我们不能直接提供真实音乐数据
	// 这里提供一个框架，实际部署时应使用官方API或合法途径
	return nil, fmt.Errorf("网易云音乐详情获取功能暂未启用，因版权原因")
}

// GetPlayURL 获取播放 URL
func (nc *NetEaseCrawler) GetPlayURL(ctx context.Context, songID string) (string, error) {
	// 注意：由于版权原因，我们不能直接提供真实音乐播放链接
	// 这里提供一个框架，实际部署时应使用官方API或合法途径
	return "", fmt.Errorf("网易云音乐播放链接获取功能暂未启用，因版权原因")
}
