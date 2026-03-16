const { contextBridge, ipcRenderer } = require('electron');

contextBridge.exposeInMainWorld('electronAPI', {
  getAppVersion: () => ipcRenderer.invoke('get-app-version'),
  getPlatform: () => ipcRenderer.invoke('get-platform'),
  showOpenDialog: (options) => ipcRenderer.invoke('show-open-dialog', options),
  
  // 播放器控制
  onPlay: (callback) => ipcRenderer.on('player-play', callback),
  onPause: (callback) => ipcRenderer.on('player-pause', callback),
  onNext: (callback) => ipcRenderer.on('player-next', callback),
  onPrevious: (callback) => ipcRenderer.on('player-previous', callback),
  
  // 系统
  minimize: () => ipcRenderer.send('window-minimize'),
  maximize: () => ipcRenderer.send('window-maximize'),
  close: () => ipcRenderer.send('window-close')
});
