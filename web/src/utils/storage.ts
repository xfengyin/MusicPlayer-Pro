/**
 * 本地存储工具函数
 * 提供统一的 localStorage 操作接口
 */

/**
 * 从 localStorage 获取数据并解析 JSON
 */
export function getFromStorage<T>(key: string, defaultValue: T): T {
  try {
    const item = localStorage.getItem(key)
    if (item === null) {
      return defaultValue
    }
    return JSON.parse(item) as T
  } catch (error) {
    console.error(`Failed to parse localStorage item "${key}":`, error)
    return defaultValue
  }
}

/**
 * 将数据序列化并存储到 localStorage
 */
export function saveToStorage<T>(key: string, value: T): void {
  try {
    localStorage.setItem(key, JSON.stringify(value))
  } catch (error) {
    console.error(`Failed to save to localStorage "${key}":`, error)
    // 处理存储空间不足的情况
    if (error instanceof DOMException && error.name === 'QuotaExceededError') {
      console.warn('LocalStorage quota exceeded, clearing old data...')
      // 可以选择清理一些旧数据
      clearStorage()
    }
  }
}

/**
 * 从 localStorage 移除数据
 */
export function removeFromStorage(key: string): void {
  try {
    localStorage.removeItem(key)
  } catch (error) {
    console.error(`Failed to remove from localStorage "${key}":`, error)
  }
}

/**
 * 清空所有 localStorage 数据
 */
export function clearStorage(): void {
  try {
    localStorage.clear()
  } catch (error) {
    console.error('Failed to clear localStorage:', error)
  }
}

/**
 * 检查 localStorage 是否可用
 */
export function isStorageAvailable(): boolean {
  try {
    const test = '__storage_test__'
    localStorage.setItem(test, test)
    localStorage.removeItem(test)
    return true
  } catch (error) {
    return false
  }
}

/**
 * 获取 localStorage 使用量（字节）
 */
export function getStorageUsage(): number {
  let total = 0
  try {
    for (const key in localStorage) {
      if (localStorage.hasOwnProperty(key)) {
        total += key.length + localStorage[key].length
      }
    }
    return total * 2 // UTF-16 每个字符 2 字节
  } catch (error) {
    console.error('Failed to calculate storage usage:', error)
    return 0
  }
}

/**
 * 获取 localStorage 剩余空间（估算）
 * 大多数浏览器限制为 5-10MB
 */
export function getStorageRemaining(limitMB: number = 5): number {
  const used = getStorageUsage()
  const limit = limitMB * 1024 * 1024
  return Math.max(0, limit - used)
}

/**
 * 存储键名常量
 */
export const StorageKeys = {
  // 播放器相关
  PLAYER_VOLUME: 'player_volume',
  PLAYER_PLAYLIST: 'player_playlist',
  PLAYER_CURRENT_INDEX: 'player_current_index',
  PLAYER_LOOP_MODE: 'player_loop_mode',
  PLAYER_PLAYBACK_RATE: 'player_playback_rate',

  // 搜索相关
  SEARCH_HISTORY: 'search_history',

  // 用户设置
  USER_SETTINGS: 'user_settings',
  THEME: 'theme',
  LANGUAGE: 'language',

  // 缓存
  LYRIC_CACHE: 'lyric_cache',
  SONG_CACHE: 'song_cache',
} as const
