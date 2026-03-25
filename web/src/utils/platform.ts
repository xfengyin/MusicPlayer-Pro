/**
 * 平台检测工具
 */

import { Capacitor } from '@capacitor/core'

/**
 * 是否运行在原生容器中（Capacitor WebView）
 */
export function isNative(): boolean {
  return Capacitor.isNativePlatform()
}

/**
 * 获取当前平台
 */
export function getPlatform(): string {
  return Capacitor.getPlatform() // 'web' | 'android' | 'ios'
}

/**
 * 是否是移动端（包括移动浏览器）
 */
export function isMobile(): boolean {
  return isNative() || /Android|iPhone|iPad|iPod/i.test(navigator.userAgent)
}

/**
 * 获取 API 基础地址
 * - 原生环境：使用远程 API
 * - 开发环境：使用本地代理
 */
export function getApiBaseUrl(): string {
  // 优先使用环境变量
  const envUrl = import.meta.env.VITE_API_BASE_URL
  
  if (isNative()) {
    // 原生 App 中，不能用 localhost，必须用远程地址
    // 如果环境变量是 localhost，则使用备用地址
    if (!envUrl || envUrl.includes('localhost') || envUrl.includes('127.0.0.1')) {
      // 从 localStorage 读取用户配置的 API 地址
      const customUrl = localStorage.getItem('api_base_url')
      if (customUrl) return customUrl
      // 默认提示需要配置
      return ''
    }
    return envUrl
  }
  
  // Web 开发环境
  return envUrl || '/api'
}
