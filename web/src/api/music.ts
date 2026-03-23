/**
 * 音乐 API 服务
 * 封装所有与音乐相关的后端 API 调用
 */

import axios, { type AxiosInstance, type AxiosResponse } from 'axios'
import type { Song, SearchResponse, LyricResponse } from '../types/music'

// API 基础配置
const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080/api'

/**
 * 创建 axios 实例
 */
const apiClient: AxiosInstance = axios.create({
  baseURL: API_BASE_URL,
  timeout: 10000, // 10 秒超时
  headers: {
    'Content-Type': 'application/json',
  },
})

// 请求拦截器
apiClient.interceptors.request.use(
  (config) => {
    // 可以在这里添加 token 等认证信息
    // const token = localStorage.getItem('token')
    // if (token) {
    //   config.headers.Authorization = `Bearer ${token}`
    // }
    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

// 响应拦截器
apiClient.interceptors.response.use(
  (response: AxiosResponse) => {
    return response.data
  },
  (error) => {
    console.error('API Error:', error)
    // 统一错误处理
    if (error.response) {
      switch (error.response.status) {
        case 401:
          // 未授权，跳转登录
          break
        case 403:
          // 禁止访问
          break
        case 404:
          // 资源不存在
          break
        case 500:
          // 服务器错误
          break
      }
    }
    return Promise.reject(error)
  }
)

/**
 * 音乐搜索 API
 * @param query 搜索关键词
 * @param page 页码（默认 1）
 * @param limit 每页数量（默认 20）
 */
export async function searchMusic(
  query: string,
  page: number = 1,
  limit: number = 20
): Promise<SearchResponse> {
  const response = await apiClient.get<SearchResponse>('/search', {
    params: { q: query, page, limit },
  })
  return response
}

/**
 * 获取歌曲详情
 * @param songId 歌曲 ID
 */
export async function getSongDetail(songId: string): Promise<Song> {
  const response = await apiClient.get<Song>(`/song/${songId}`)
  return response
}

/**
 * 获取歌曲播放 URL
 * @param songId 歌曲 ID
 */
export async function getSongUrl(songId: string): Promise<{ url: string }> {
  const response = await apiClient.get<{ url: string }>(`/song/${songId}/url`)
  return response
}

/**
 * 获取歌词
 * @param songId 歌曲 ID
 */
export async function getLyrics(songId: string): Promise<LyricResponse> {
  const response = await apiClient.get<LyricResponse>(`/song/${songId}/lyric`)
  return response
}

/**
 * 获取推荐歌单
 * @param limit 数量限制
 */
export async function getRecommendPlaylists(limit: number = 10): Promise<any> {
  const response = await apiClient.get('/music/playlist/recommend', {
    params: { limit },
  })
  return response
}

/**
 * 获取热门歌曲
 * @param limit 数量限制
 */
export async function getHotSongs(limit: number = 20): Promise<Song[]> {
  const response = await apiClient.get<Song[]>('/music/hot', {
    params: { limit },
  })
  return response
}

/**
 * 解析 LRC 歌词文本
 * @param lrcText LRC 格式歌词文本
 */
export function parseLyrics(lrcText: string): Array<{ time: number; text: string; translation?: string }> {
  if (!lrcText) return []

  const lines = lrcText.split('\n')
  const lyrics: Array<{ time: number; text: string; translation?: string }> = []
  const timeRegex = /\[(\d{2}):(\d{2})\.(\d{2,3})\]/

  // 用于合并翻译的临时存储
  const lyricsMap = new Map<number, { time: number; text: string; translation?: string }>()

  lines.forEach((line) => {
    const match = timeRegex.exec(line)
    if (match) {
      const minutes = parseInt(match[1])
      const seconds = parseInt(match[2])
      const milliseconds = parseInt(match[3].padEnd(3, '0'))
      const time = minutes * 60 + seconds + milliseconds / 1000

      // 移除时间标签获取歌词文本
      const text = line.replace(timeRegex, '').trim()

      if (text) {
        // 检查是否是翻译（通常翻译行没有时间标签或格式不同）
        // 这里简化处理，假设所有带时间标签的都是主歌词
        const existingLyric = lyricsMap.get(time)
        if (existingLyric) {
          // 如果已有相同时间的歌词，可能是翻译
          existingLyric.translation = text
        } else {
          const lyricLine = { time, text }
          lyricsMap.set(time, lyricLine)
          lyrics.push(lyricLine)
        }
      }
    }
  })

  // 按时间排序
  return lyrics.sort((a, b) => a.time - b.time)
}

export default {
  searchMusic,
  getSongDetail,
  getSongUrl,
  getLyrics,
  getRecommendPlaylists,
  getHotSongs,
  parseLyrics,
}
