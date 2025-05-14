import {app, BrowserWindow, ipcMain} from 'electron';
import path from 'node:path';
import {startServer, storeInfo} from "./server";
import Store from 'electron-store';
import {initTray, quitApp} from "./tray";
import {startBackend} from "./admin";

// 是否在开发模式
const isDev = !app.isPackaged;


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

    console.log("数据库初始化完成")
}



// 窗口
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
        titleBarStyle: 'hiddenInset',
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

app.on('activate', () => {
    if (mainWindow && !mainWindow.isVisible()) {
        mainWindow.show();
    }
});

app.on('window-all-closed', () => {
    if (process.platform !== 'darwin') {
        app.quit();
    }
});


// 等待 backend 传来的 port 和 secret
let resolveReady: () => void;
const waitForReady = new Promise<void>((resolve) => {
    resolveReady = resolve;
});

// 单例模式
const gotTheLock = app.requestSingleInstanceLock();
if (!gotTheLock) {
    quitApp()
} else {
    app.on('second-instance', () => {
        if (mainWindow) {
            if (mainWindow.isMinimized()) mainWindow.restore();
            mainWindow.focus();
        }
    });

    startServer(resolveReady, startBackend)

    app.whenReady().then(async () => {
        // 等待后端启动
        await waitForReady;
        initStore(storeInfo.home())

        // 启动UI
        console.log('准备就绪，启动窗口，port=', storeInfo.port(), ' secret=', storeInfo.secret());
        createWindow();
    });
}