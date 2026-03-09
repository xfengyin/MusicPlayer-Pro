package model

import (
	"time"

	"gorm.io/gorm"
)

// Song 歌曲数据模型
type Song struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Title     string         `gorm:"size:255;not null" json:"title"`      // 歌曲标题
	Artist    string         `gorm:"size:255;not null" json:"artist"`     // 艺术家
	Album     string         `gorm:"size:255" json:"album"`               // 专辑名称
	Duration  int            `gorm:"not null" json:"duration"`            // 时长 (秒)
	URL       string         `gorm:"size:512;not null" json:"url"`        // 播放链接
	CoverURL  string         `gorm:"size:512" json:"cover_url"`           // 封面图片链接
	LYRIC     string         `gorm:"size:1024" json:"lyric"`              // 歌词
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (Song) TableName() string {
	return "songs"
}

// Playlist 播放列表数据模型
type Playlist struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"size:255;not null" json:"name"`        // 播放列表名称
	Description string         `gorm:"size:512" json:"description"`          // 描述
	CoverURL    string         `gorm:"size:512" json:"cover_url"`            // 封面图片链接
	UserID      uint           `gorm:"not null" json:"user_id"`              // 用户 ID
	Songs       []Song         `gorm:"many2many:playlist_songs;" json:"songs"` // 歌曲列表
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (Playlist) TableName() string {
	return "playlists"
}

// PlaylistSong 播放列表与歌曲的关联表
type PlaylistSong struct {
	PlaylistID uint `gorm:"primaryKey"`
	SongID     uint `gorm:"primaryKey"`
}

// TableName 指定关联表表名
func (PlaylistSong) TableName() string {
	return "playlist_songs"
}
