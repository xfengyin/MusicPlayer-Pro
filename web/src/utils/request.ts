/**
 * Axios 实例配置
 * 提供统一的 HTTP 请求配置和拦截器
 */

import axios, { type AxiosInstance, type AxiosRequestConfig, type AxiosResponse } from 'axios'

// API 基础配置
const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080/api'
const API_TIMEOUT = parseInt(import.meta.env.VITE_API_TIMEOUT || '10000', 10)

/**
 * 创建 axios 实例
 */
export const apiClient: AxiosInstance = axios.create({
  baseURL: API_BASE_URL,
  timeout: API_TIMEOUT,
  headers: {
    'Content-Type': 'application/json',
  },
  // 跨域请求时携带凭证
  withCredentials: false,
  // 最大重定向次数
  maxRedirects: 5,
})

/**
 * 请求拦截器
 */
apiClient.interceptors.request.use(
  (config) => {
    // 添加认证 token（如果存在）
    const token = localStorage.getItem('auth_token')
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }

    // 添加请求时间戳（防止缓存）
    if (config.method === 'get') {
      config.params = {
        ...config.params,
        _t: Date.now(),
      }
    }

    return config
  },
  (error) => {
    console.error('Request error:', error)
    return Promise.reject(error)
  }
)

/**
 * 响应拦截器
 */
apiClient.interceptors.response.use(
  (response: AxiosResponse) => {
    // 直接返回 data
    return response.data
  },
  (error) => {
    // 统一错误处理
    if (error.response) {
      const { status, data } = error.response

      switch (status) {
        case 400:
          console.error('请求参数错误')
          break
        case 401:
          console.error('未授权，请重新登录')
          // 可以在这里触发登录流程
          // localStorage.removeItem('auth_token')
          // window.location.href = '/login'
          break
        case 403:
          console.error('禁止访问')
          break
        case 404:
          console.error('请求的资源不存在')
          break
        case 500:
          console.error('服务器内部错误')
          break
        case 502:
          console.error('网关错误')
          break
        case 503:
          console.error('服务不可用')
          break
        default:
          console.error(`HTTP 错误：${status}`)
      }
    } else if (error.request) {
      // 请求已发送但没有收到响应
      console.error('网络错误，请检查网络连接')
    } else {
      // 请求配置出错
      console.error('请求配置错误:', error.message)
    }

    return Promise.reject(error)
  }
)

/**
 * 创建自定义配置的 axios 实例
 */
export function createApiClient(config?: AxiosRequestConfig): AxiosInstance {
  const instance = axios.create({
    baseURL: API_BASE_URL,
    timeout: API_TIMEOUT,
    ...config,
  })

  // 应用相同的拦截器
  instance.interceptors.request.use(apiClient.interceptors.request.handlers[0].fulfilled)
  instance.interceptors.response.use(
    apiClient.interceptors.response.handlers[0].fulfilled,
    apiClient.interceptors.response.handlers[0].rejected
  )

  return instance
}

export default apiClient
