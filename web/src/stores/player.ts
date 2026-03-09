import { defineStore } from 'pinia'
import { ref, computed, watch } from 'vue'
import type { Song } from '../types/music'
import { getSongUrl, getLyrics, parseLyrics } from '../api/music'

export const usePlayerStore = defineStore('player', () => {
  // State
  const currentSong = ref<Song | null>(null)
  const playlist = ref<Song[]>([])
  const currentIndex = ref(-1)
  const isPlaying = ref(false)
  const currentTime = ref(0)
  const duration = ref(0)
  const volume = ref(0.8)
  const isMuted = ref(false)
  const loopMode = ref<'single' | 'list' | 'random'>('list')
  const lyrics = ref<Array<{ time: number; text: string; translation?: string }>>([])
  const isLoading = ref(false)
  const error = ref<string | null>(null)

  // HTML5 Audio 实例
  let audio: HTMLAudioElement | null = null

  // 初始化 Audio
  function initAudio() {
    if (!audio) {
      audio = new Audio()
      audio.preload = 'metadata'

      // 监听音频事件
      audio.addEventListener('timeupdate', onTimeUpdate)
      audio.addEventListener('loadedmetadata', onLoadedMetadata)
      audio.addEventListener('ended', onEnded)
      audio.addEventListener('error', onError)
      audio.addEventListener('canplay', onCanPlay)
    }
  }

  // 事件处理函数
  function onTimeUpdate() {
    if (audio) {
      currentTime.value = audio.currentTime
    }
  }

  function onLoadedMetadata() {
    if (audio) {
      duration.value = audio.duration
    }
  }

  function onEnded() {
    // 根据播放模式处理播放结束
    if (loopMode.value === 'single') {
      // 单曲循环
      if (audio) {
        audio.currentTime = 0
        audio.play()
      }
    } else if (loopMode.value === 'random') {
      // 随机播放
      playRandomSong()
    } else {
      // 列表循环
      nextSong()
    }
  }

  function onError(e: Event) {
    console.error('Audio error:', e)
    error.value = '播放失败'
    isLoading.value = false
  }

  function onCanPlay() {
    isLoading.value = false
  }

  // Getters
  const progress = computed(() => {
    if (duration.value === 0) return 0
    return (currentTime.value / duration.value) * 100
  })

  const hasNextSong = computed(() => {
    if (loopMode.value === 'random') return playlist.value.length > 1
    return currentIndex.value < playlist.value.length - 1
  })

  const hasPrevSong = computed(() => {
    if (loopMode.value === 'random') return playlist.value.length > 1
    return currentIndex.value > 0
  })

  // 从 localStorage 加载音量和播放列表
  function loadFromStorage() {
    try {
      // 加载音量
      const savedVolume = localStorage.getItem('player_volume')
      if (savedVolume) {
        volume.value = parseFloat(savedVolume)
      }

      // 加载播放列表
      const savedPlaylist = localStorage.getItem('player_playlist')
      if (savedPlaylist) {
        playlist.value = JSON.parse(savedPlaylist)
      }

      // 加载当前播放索引
      const savedIndex = localStorage.getItem('player_current_index')
      if (savedIndex) {
        currentIndex.value = parseInt(savedIndex)
        if (currentIndex.value >= 0 && currentIndex.value < playlist.value.length) {
          currentSong.value = playlist.value[currentIndex.value]
        }
      }

      // 加载播放模式
      const savedLoopMode = localStorage.getItem('player_loop_mode')
      if (savedLoopMode && ['single', 'list', 'random'].includes(savedLoopMode)) {
        loopMode.value = savedLoopMode as 'single' | 'list' | 'random'
      }
    } catch (e) {
      console.error('Failed to load from storage:', e)
    }
  }

  // 保存到 localStorage
  function saveToStorage() {
    try {
      localStorage.setItem('player_volume', volume.value.toString())
      localStorage.setItem('player_playlist', JSON.stringify(playlist.value))
      localStorage.setItem('player_current_index', currentIndex.value.toString())
      localStorage.setItem('player_loop_mode', loopMode.value)
    } catch (e) {
      console.error('Failed to save to storage:', e)
    }
  }

  // 监听状态变化并保存
  watch([volume, playlist, currentIndex, loopMode], () => {
    saveToStorage()
  })

  // Actions
  /**
   * 播放歌曲
   */
  async function playSong(song: Song, index?: number) {
    initAudio()
    isLoading.value = true
    error.value = null

    try {
      // 更新当前歌曲
      currentSong.value = song
      if (index !== undefined) {
        currentIndex.value = index
      } else {
        currentIndex.value = playlist.value.findIndex((s) => s.id === song.id)
      }

      // 获取播放 URL
      const { url } = await getSongUrl(song.id)

      // 设置音频源并播放
      audio!.src = url
      audio!.volume = isMuted.value ? 0 : volume.value
      await audio!.play()
      isPlaying.value = true

      // 加载歌词
      await loadLyrics(song.id)
    } catch (e: any) {
      console.error('Failed to play song:', e)
      error.value = e.message || '播放失败'
      isLoading.value = false
      isPlaying.value = false
    }
  }

  /**
   * 加载歌词
   */
  async function loadLyrics(songId: string) {
    try {
      const response = await getLyrics(songId)
      if (response.lrc) {
        lyrics.value = parseLyrics(response.lrc)
      } else {
        lyrics.value = []
      }
    } catch (e) {
      console.error('Failed to load lyrics:', e)
      lyrics.value = []
    }
  }

  /**
   * 切换播放/暂停
   */
  async function togglePlay() {
    if (!audio || !currentSong.value) return

    if (isPlaying.value) {
      audio.pause()
      isPlaying.value = false
    } else {
      try {
        await audio.play()
        isPlaying.value = true
      } catch (e) {
        console.error('Failed to play:', e)
      }
    }
  }

  /**
   * 播放
   */
  async function play() {
    if (!audio || !currentSong.value) return
    try {
      await audio.play()
      isPlaying.value = true
    } catch (e) {
      console.error('Failed to play:', e)
    }
  }

  /**
   * 暂停
   */
  function pause() {
    if (!audio) return
    audio.pause()
    isPlaying.value = false
  }

  /**
   * 下一首
   */
  async function nextSong() {
    if (playlist.value.length === 0) return

    let nextIndex: number

    if (loopMode.value === 'random') {
      // 随机播放
      do {
        nextIndex = Math.floor(Math.random() * playlist.value.length)
      } while (playlist.value.length > 1 && nextIndex === currentIndex.value)
    } else {
      // 顺序播放
      nextIndex = currentIndex.value + 1
      if (nextIndex >= playlist.value.length) {
        nextIndex = 0 // 列表循环
      }
    }

    if (nextIndex >= 0 && nextIndex < playlist.value.length) {
      await playSong(playlist.value[nextIndex], nextIndex)
    }
  }

  /**
   * 上一首
   */
  async function prevSong() {
    if (playlist.value.length === 0) return

    let prevIndex: number

    if (loopMode.value === 'random') {
      // 随机播放
      do {
        prevIndex = Math.floor(Math.random() * playlist.value.length)
      } while (playlist.value.length > 1 && prevIndex === currentIndex.value)
    } else {
      // 顺序播放
      prevIndex = currentIndex.value - 1
      if (prevIndex < 0) {
        prevIndex = playlist.value.length - 1 // 列表循环
      }
    }

    if (prevIndex >= 0 && prevIndex < playlist.value.length) {
      await playSong(playlist.value[prevIndex], prevIndex)
    }
  }

  /**
   * 随机播放一首
   */
  async function playRandomSong() {
    if (playlist.value.length === 0) return

    let randomIndex: number
    do {
      randomIndex = Math.floor(Math.random() * playlist.value.length)
    } while (playlist.value.length > 1 && randomIndex === currentIndex.value)

    await playSong(playlist.value[randomIndex], randomIndex)
  }

  /**
   * 设置音量
   */
  function setVolume(value: number) {
    volume.value = Math.max(0, Math.min(1, value))
    isMuted.value = volume.value === 0

    if (audio) {
      audio.volume = isMuted.value ? 0 : volume.value
    }

    localStorage.setItem('player_volume', volume.value.toString())
  }

  /**
   * 切换静音
   */
  function toggleMute() {
    isMuted.value = !isMuted.value
    if (audio) {
      audio.volume = isMuted.value ? 0 : volume.value
    }
  }

  /**
   * 跳转到指定时间
   */
  function seekTo(time: number) {
    if (!audio) return
    audio.currentTime = Math.max(0, Math.min(duration.value, time))
    currentTime.value = audio.currentTime
  }

  /**
   * 设置播放模式
   */
  function setLoopMode(mode: 'single' | 'list' | 'random') {
    loopMode.value = mode
    localStorage.setItem('player_loop_mode', mode)
  }

  /**
   * 添加到播放列表
   */
  function addToPlaylist(song: Song) {
    const exists = playlist.value.find((s) => s.id === song.id)
    if (!exists) {
      playlist.value.push(song)
      saveToStorage()
    }
  }

  /**
   * 从播放列表移除
   */
  function removeFromPlaylist(index: number) {
    if (index < 0 || index >= playlist.value.length) return

    playlist.value.splice(index, 1)

    // 调整当前索引
    if (index < currentIndex.value) {
      currentIndex.value--
    } else if (index === currentIndex.value) {
      // 如果移除的是当前歌曲
      if (playlist.value.length === 0) {
        currentSong.value = null
        currentIndex.value = -1
        if (audio) {
          audio.pause()
          audio.src = ''
        }
        isPlaying.value = false
      } else if (currentIndex.value >= playlist.value.length) {
        currentIndex.value = playlist.value.length - 1
        currentSong.value = playlist.value[currentIndex.value]
      }
    }

    saveToStorage()
  }

  /**
   * 清空播放列表
   */
  function clearPlaylist() {
    playlist.value = []
    currentSong.value = null
    currentIndex.value = -1
    if (audio) {
      audio.pause()
      audio.src = ''
    }
    isPlaying.value = false
    currentTime.value = 0
    duration.value = 0
    lyrics.value = []
    saveToStorage()
  }

  /**
   * 重排播放列表
   */
  function reorderPlaylist(fromIndex: number, toIndex: number) {
    if (
      fromIndex < 0 ||
      fromIndex >= playlist.value.length ||
      toIndex < 0 ||
      toIndex >= playlist.value.length
    ) {
      return
    }

    const [removed] = playlist.value.splice(fromIndex, 1)
    playlist.value.splice(toIndex, 0, removed)

    // 调整当前索引
    if (fromIndex === currentIndex.value) {
      currentIndex.value = toIndex
    } else if (fromIndex < currentIndex.value && toIndex >= currentIndex.value) {
      currentIndex.value--
    } else if (fromIndex > currentIndex.value && toIndex <= currentIndex.value) {
      currentIndex.value++
    }

    saveToStorage()
  }

  /**
   * 设置播放列表
   */
  function setPlaylist(songs: Song[], startIndex?: number) {
    playlist.value = songs
    currentIndex.value = startIndex ?? 0
    if (startIndex !== undefined && songs.length > 0) {
      currentSong.value = songs[startIndex]
    }
    saveToStorage()
  }

  // 初始化时加载存储的数据
  loadFromStorage()

  return {
    // State
    currentSong,
    playlist,
    currentIndex,
    isPlaying,
    currentTime,
    duration,
    volume,
    isMuted,
    loopMode,
    lyrics,
    isLoading,
    error,
    // Getters
    progress,
    hasNextSong,
    hasPrevSong,
    // Actions
    playSong,
    togglePlay,
    play,
    pause,
    nextSong,
    prevSong,
    playRandomSong,
    setVolume,
    toggleMute,
    seekTo,
    setLoopMode,
    addToPlaylist,
    removeFromPlaylist,
    clearPlaylist,
    reorderPlaylist,
    setPlaylist,
  }
})
