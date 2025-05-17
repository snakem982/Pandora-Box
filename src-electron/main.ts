import {app, BrowserWindow, ipcMain} from 'electron';
import path from 'node:path';
import {startServer, storeInfo} from "./server";
import Store from 'electron-store';
import {doQuit, initTray} from "./tray";
import {startBackend} from "./admin";
import log from './log';

// 是否在开发模式
const isDev = !app.isPackaged;

// 初始化数据库
function initStore(home: string) {
    const store = new Store({
        cwd: path.join(home, 'px-electron.db')
    });

    ipcMain.handle('store:get', (event, key) => {
        return store.get(key);
    });

    ipcMain.handle('store:set', (event, key, value) => {
        store.set(key, value);
    });

    log.info("数据库初始化完成")
}

// 主窗口
let mainWindow: BrowserWindow;
// 屏蔽安全警告
process.env["ELECTRON_DISABLE_SECURITY_WARNINGS"] = "true";
const createWindow = () => {
    mainWindow = new BrowserWindow({
        width: 1100,
        height: 760,
        minWidth: 960,
        minHeight: 660,
        center: true,
        webPreferences: {
            preload: path.join(__dirname, 'preload.js'),
            contextIsolation: true,
            webSecurity: false,
            nodeIntegrationInWorker: true
        },
        // expose window controls in Windows/Linux
        ...(process.platform !== 'darwin' ? {
            titleBarStyle: 'hidden'
        } : {titleBarStyle: 'hiddenInset'})
    });

    // 加载托盘
    initTray(mainWindow)

    // 隐藏菜单栏
    mainWindow.setMenu(null);

    if (isDev) {
        const filePath = `http://localhost:5173?port=${storeInfo.port()}&secret=${storeInfo.secret()}`;
        mainWindow.loadURL(filePath);
    } else {
        const filePath = `http://${storeInfo.listenAddr()}/index.html?port=${storeInfo.port()}&secret=${storeInfo.secret()}`;
        console.log('准备就绪，加载窗口，url=', filePath);
        mainWindow.loadURL(filePath);
    }
};

// 等待 backend 传来的 port 和 secret
let resolveReady: () => void;
const waitForReady = new Promise<void>((resolve) => {
    resolveReady = resolve;
});

// 单例模式
const gotTheLock = app.requestSingleInstanceLock();
if (!gotTheLock) {
    doQuit()
} else {
    // 试图启动第二个应用实例
    app.on('second-instance', () => {
        if (mainWindow) {
            mainWindow.show();
            app.dock?.show();
            mainWindow.focus();
        }
    });

    // 监听应用被激活
    app.on('activate', () => {
        if (mainWindow && !mainWindow.isVisible()) {
            mainWindow.show();
            app.dock?.show();
            mainWindow.focus();
        }
    });

    app.whenReady().then(async () => {
        // 启动前端静态服务
        startServer(resolveReady, startBackend)

        // 等待后端启动后初始化前端数据库
        await waitForReady;
        initStore(log.getHomeDir())

        // 启动UI
        log.info('准备就绪，启动窗口，port=', storeInfo.port(), ' secret=', storeInfo.secret());
        createWindow();
    });
}