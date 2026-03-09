/**
 * 音乐相关类型定义
 */

/**
 * 歌曲信息
 */
export interface Song {
  id: string
  title: string
  artist: string
  album: string
  cover: string
  url: string
  duration: number
  [key: string]: any // 允许扩展字段
}

/**
 * 搜索结果
 */
export interface SearchResponse {
  songs: Song[]
  total: number
  page: number
  limit: number
}

/**
 * 歌词响应
 */
export interface LyricResponse {
  songId: string
  lrc?: string // LRC 格式歌词
  translation?: string // 翻译歌词（可选）
}

/**
 * 歌词行（解析后）
 */
export interface LyricLine {
  time: number // 时间（秒）
  text: string // 歌词文本
  translation?: string // 翻译（可选）
}

/**
 * 播放列表
 */
export interface Playlist {
  id: string
  name: string
  cover: string
  description: string
  songs: Song[]
  creator: {
    id: string
    name: string
    avatar: string
  }
  playCount: number
  createTime: number
  updateTime: number
}

/**
 * 播放器状态
 */
export interface PlayerState {
  currentSong: Song | null
  playlist: Song[]
  currentIndex: number
  isPlaying: boolean
  currentTime: number
  duration: number
  volume: number
  isMuted: boolean
  loopMode: 'single' | 'list' | 'random'
  playbackRate: number
}
