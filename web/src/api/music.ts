/**
 * 音乐 API 服务
 * 封装所有与音乐相关的后端 API 调用
 */

import axios, { type AxiosInstance, type AxiosResponse } from 'axios'
import type { Song, SearchResponse, LyricResponse } from '../types/music'
import { getApiBaseUrl, isNative } from '../utils/platform'

/**
 * 获取 API 基础 URL
 * 优先级：localStorage 自定义 > 环境变量 > 默认值
 */
function resolveApiBaseUrl(): string {
  const customUrl = localStorage.getItem('api_base_url')
  if (customUrl) return customUrl
  return getApiBaseUrl() || import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080/api'
}

/**
 * 创建 axios 实例
 */
let apiClient: AxiosInstance = createApiClient()

function createApiClient(): AxiosInstance {
  const client = axios.create({
    baseURL: resolveApiBaseUrl(),
    timeout: 15000,
    headers: {
      'Content-Type': 'application/json',
    },
  })

  // 请求拦截器
  client.interceptors.request.use(
    (config) => {
      // 每次请求时重新检查 API 地址（支持运行时切换）
      const currentUrl = resolveApiBaseUrl()
      if (currentUrl && config.baseURL !== currentUrl) {
        config.baseURL = currentUrl
      }
      return config
    },
    (error) => Promise.reject(error)
  )

  // 响应拦截器
  client.interceptors.response.use(
    (response: AxiosResponse) => response.data,
    (error) => {
      console.error('API Error:', error)
      if (error.code === 'ERR_NETWORK' && isNative()) {
        console.error(
          '网络错误：请在设置中配置正确的 API 服务器地址。' +
          '当前地址：' + resolveApiBaseUrl()
        )
      }
      return Promise.reject(error)
    }
  )

  return client
}

/**
 * 刷新 API 客户端（切换 API 地址后调用）
 */
export function refreshApiClient() {
  apiClient = createApiClient()
}

/**
 * 音乐搜索 API
 */
export async function searchMusic(
  query: string,
  page: number = 1,
  limit: number = 20
): Promise<SearchResponse> {
  const response = await apiClient.get('/search', {
    params: { q: query, page, limit },
  })
  // 后端返回 { data: { songs, total, ... } }
  const data = response?.data || response
  return {
    songs: (data?.songs || []).map(mapSong),
    total: data?.total || 0,
    page,
    limit,
  }
}

/**
 * 映射后端 Song 到前端 Song 类型
 */
function mapSong(raw: any): Song {
  return {
    id: String(raw.ID || raw.id || ''),
    title: raw.Title || raw.title || '未知歌曲',
    artist: raw.Artist || raw.artist || '未知歌手',
    album: raw.Album || raw.album || '',
    cover: raw.CoverURL || raw.cover_url || raw.cover || '',
    url: raw.URL || raw.url || '',
    duration: raw.Duration || raw.duration || 0,
  }
}

/**
 * 获取歌曲详情
 */
export async function getSongDetail(songId: string): Promise<Song> {
  const response = await apiClient.get(`/song/${songId}`)
  const data = response?.data || response
  return mapSong(data)
}

/**
 * 获取歌曲播放 URL
 */
export async function getSongUrl(songId: string): Promise<{ url: string }> {
  const response = await apiClient.get(`/song/${songId}/url`)
  const data = response?.data || response
  return { url: data?.url || '' }
}

/**
 * 获取歌词
 */
export async function getLyrics(songId: string): Promise<LyricResponse> {
  const response = await apiClient.get(`/song/${songId}/lyric`)
  const data = response?.data || response
  return {
    songId,
    lrc: data?.lyric || data?.lrc || '',
    translation: data?.tlyric || data?.translation || '',
  }
}

/**
 * 获取热门搜索
 */
export async function getHotSearches(): Promise<string[]> {
  try {
    const response = await apiClient.get('/hot-searches')
    const data = response?.data || response
    return data?.hot_searches || []
  } catch {
    return []
  }
}

/**
 * 获取推荐歌单
 */
export async function getRecommendPlaylists(limit: number = 10): Promise<any> {
  const response = await apiClient.get('/playlist', {
    params: { limit },
  })
  return response?.data || response
}

/**
 * 获取热门歌曲
 */
export async function getHotSongs(limit: number = 20): Promise<Song[]> {
  try {
    const response = await apiClient.get('/music/hot', {
      params: { limit },
    })
    const data = response?.data || response
    return (Array.isArray(data) ? data : []).map(mapSong)
  } catch {
    return []
  }
}

/**
 * 解析 LRC 歌词文本
 */
export function parseLyrics(lrcText: string): Array<{ time: number; text: string; translation?: string }> {
  if (!lrcText) return []

  const lines = lrcText.split('\n')
  const lyrics: Array<{ time: number; text: string; translation?: string }> = []
  const timeRegex = /\[(\d{2}):(\d{2})\.(\d{2,3})\]/
  const lyricsMap = new Map<number, { time: number; text: string; translation?: string }>()

  lines.forEach((line) => {
    const match = timeRegex.exec(line)
    if (match) {
      const minutes = parseInt(match[1])
      const seconds = parseInt(match[2])
      const milliseconds = parseInt(match[3].padEnd(3, '0'))
      const time = minutes * 60 + seconds + milliseconds / 1000
      const text = line.replace(timeRegex, '').trim()

      if (text) {
        const existingLyric = lyricsMap.get(time)
        if (existingLyric) {
          existingLyric.translation = text
        } else {
          const lyricLine = { time, text }
          lyricsMap.set(time, lyricLine)
          lyrics.push(lyricLine)
        }
      }
    }
  })

  return lyrics.sort((a, b) => a.time - b.time)
}

export default {
  searchMusic,
  getSongDetail,
  getSongUrl,
  getLyrics,
  getHotSearches,
  getRecommendPlaylists,
  getHotSongs,
  parseLyrics,
  refreshApiClient,
}
