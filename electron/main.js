const { app, BrowserWindow, ipcMain, Tray, Menu, nativeImage, dialog, shell } = require('electron');
const path = require('path');
const log = require('electron-log');

let mainWindow = null;
let tray = null;

// 启用 Electron 安全警告
process.env.ELECTRON_DISABLE_SECURITY_WARNINGS = 'false';

// 捕获未处理的异常 - 增强版
process.on('uncaughtException', (error) => {
  log.error('Uncaught Exception:', error);
  if (mainWindow) {
    mainWindow.webContents.send('app-error', {
      type: 'uncaughtException',
      message: error.message,
      stack: error.stack
    });
  }
  // 生产环境显示友好错误，开发环境显示详细错误
  if (process.env.NODE_ENV === 'production') {
    dialog.showErrorBox('应用程序错误', '发生了一个意外错误，请重启应用程序。');
  } else {
    dialog.showErrorBox('Error', error.message);
  }
});

process.on('unhandledRejection', (reason, promise) => {
  log.error('Unhandled Rejection at:', promise, 'reason:', reason);
  // 不中断应用，只记录日志
});

function createWindow() {
  mainWindow = new BrowserWindow({
    width: 1200,
    height: 800,
    minWidth: 800,
    minHeight: 600,
    show: false,
    // 安全配置 - 增强版
    webPreferences: {
      nodeIntegration: false,
      contextIsolation: true,
      preload: path.join(__dirname, 'preload.js'),
      webSecurity: true,
      sandbox: true,  // 启用沙箱模式
      allowRunningInsecureContent: false,  // 禁止混合内容
      experimentalFeatures: false,  // 禁用实验性功能
      enableRemoteModule: false,  // 禁用 remote 模块
      safeDialogs: true,  // 启用安全对话框
    },
    icon: path.join(__dirname, 'build', 'icon.png'),
    title: 'MusicPlayerPro',
    // 性能优化
    backgroundColor: '#1a1a2e',
    paintWhenInitiallyHidden: true,
  });

  // 加载前端构建产物
  const isDev = process.env.NODE_ENV === 'development' || process.env.ELECTRON_ENV === 'development';
  
  if (isDev) {
    mainWindow.loadURL('http://localhost:5173');
    mainWindow.webContents.openDevTools();
    // 开发环境禁用 CSP 以便调试
  } else {
    const indexPath = path.join(__dirname, 'renderer', 'index.html');
    mainWindow.loadFile(indexPath);
    // 生产环境设置 Content Security Policy
    session.defaultSession.webRequest.onHeadersReceived((details, callback) => {
      callback({
        responseHeaders: {
          ...details.responseHeaders,
          'Content-Security-Policy': ["default-src 'self'; script-src 'self' 'unsafe-inline'; style-src 'self' 'unsafe-inline'"]
        }
      });
    });
  }

  // 窗口准备就绪后显示
  mainWindow.once('ready-to-show', () => {
    mainWindow.show();
    log.info('Main window shown');
  });

  // 窗口关闭事件
  mainWindow.on('closed', () => {
    mainWindow = null;
  });

  // 阻止新窗口创建，使用默认浏览器打开
  mainWindow.webContents.setWindowOpenHandler(({ url }) => {
    // 只允许外部链接，阻止所有新窗口
    if (url.startsWith('http://') || url.startsWith('https://')) {
      shell.openExternal(url);
      return { action: 'deny' };
    }
    return { action: 'deny' };
  });

  // 导航拦截 - 防止加载外部内容
  mainWindow.webContents.on('will-navigate', (event, url) => {
    const parsedUrl = new URL(url);
    if (parsedUrl.origin !== 'http://localhost:5173' && !url.startsWith('file://')) {
      event.preventDefault();
      shell.openExternal(url);
    }
  });
}

function createTray() {
  const iconPath = path.join(__dirname, 'build', 'icon.png');
  try {
    tray = new Tray(nativeImage.createFromPath(iconPath).resize({ width: 16, height: 16 }));
  } catch (e) {
    log.warn('Tray icon not found, using empty icon');
    tray = new Tray(nativeImage.createEmpty());
  }

  const contextMenu = Menu.buildFromTemplate([
    { label: '显示', click: () => mainWindow?.show() },
    { label: '隐藏', click: () => mainWindow?.hide() },
    { type: 'separator' },
    { label: '退出', click: () => app.quit() }
  ]);

  tray.setToolTip('MusicPlayerPro');
  tray.setContextMenu(contextMenu);
  tray.on('click', () => mainWindow?.show());
}

// IPC 处理器 - 安全版本
ipcMain.handle('app:getVersion', () => {
  return app.getVersion();
});

ipcMain.handle('app:getPlatform', () => {
  return process.platform;
});

ipcMain.handle('dialog:openFile', async () => {
  const result = await dialog.showOpenDialog(mainWindow, {
    properties: ['openFile'],
    filters: [
      { name: 'Audio Files', extensions: ['mp3', 'flac', 'wav', 'ogg', 'm4a'] },
      { name: 'All Files', extensions: ['*'] }
    ]
  });
  
  if (!result.canceled && result.filePaths.length > 0) {
    return result.filePaths[0];
  }
  return null;
});

// 应用生命周期
app.whenReady().then(() => {
  createWindow();
  createTray();
  log.info('Application ready');
});

app.on('window-all-closed', () => {
  if (process.platform !== 'darwin') {
    app.quit();
  }
});

app.on('activate', () => {
  if (BrowserWindow.getAllWindows().length === 0) {
    createWindow();
  }
});

// 安全退出
app.on('before-quit', () => {
  log.info('Application quitting');
  if (mainWindow) {
    mainWindow.removeAllListeners('close');
    mainWindow.close();
  }
});
