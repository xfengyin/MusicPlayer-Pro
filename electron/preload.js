/**
 * Preload Script
 * 安全地暴露 API 给渲染进程
 */

const { contextBridge, ipcRenderer } = require('electron');

// 定义暴露的 API
const api = {
  // 应用信息
  app: {
    getVersion: () => ipcRenderer.invoke('app:getVersion'),
    getPlatform: () => ipcRenderer.invoke('app:getPlatform'),
  },
  
  // 文件对话框
  dialog: {
    openFile: () => ipcRenderer.invoke('dialog:openFile'),
  },
  
  // 事件监听
  onAppError: (callback) => {
    ipcRenderer.on('app-error', (event, error) => callback(error));
  },
  
  // 移除事件监听
  removeAllListeners: (channel) => {
    ipcRenderer.removeAllListeners(channel);
  },
};

// 安全地暴露 API
contextBridge.exposeInMainWorld('electronAPI', api);

// 类型声明（供 TypeScript 使用）
if (process.env.NODE_ENV === 'development') {
  console.log('[Preload] API exposed securely');
}
