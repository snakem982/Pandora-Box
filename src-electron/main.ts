import {app, BrowserWindow, BrowserWindowConstructorOptions} from 'electron';
import path from 'node:path';
import {startServer, storeInfo} from "./server";
import {doQuit, initTray, showWindow} from "./tray";
import {startBackend} from "./admin";
import log from './log';
import {initStore, storeGet} from "./store";

// 是否在开发模式
const isDev = !app.isPackaged;


// 主窗口
let mainWindow: BrowserWindow;
// 屏蔽安全警告
process.env["ELECTRON_DISABLE_SECURITY_WARNINGS"] = "true";
const createWindow = () => {

    // 窗口配置
    let windowOptions: BrowserWindowConstructorOptions = {
        minWidth: 960,
        minHeight: 660,
        width: 1100,
        height: 760,
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
    }

    // 从 store 获取保存的窗口尺寸与位置
    const savedBounds: any = storeGet('windowBounds');
    if (savedBounds && savedBounds.x !== undefined && savedBounds.y !== undefined) {
        windowOptions = {
            ...windowOptions,
            x: savedBounds.x,
            y: savedBounds.y,
            width: savedBounds.width,
            height: savedBounds.height
        }
    }

    // 创建窗口
    mainWindow = new BrowserWindow(windowOptions);

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
    app.on('second-instance', showWindow);

    // 监听应用被激活
    app.on('activate', showWindow);

    app.whenReady().then(async () => {
        // 初始化前端数据库
        initStore(log.getHomeDir())

        // 启动前端静态服务
        startServer(resolveReady, startBackend)

        // 等待后端启动
        await waitForReady;

        // 启动UI
        log.info('准备就绪，启动窗口，port=', storeInfo.port(), ' secret=', storeInfo.secret());
        createWindow();
    });
}