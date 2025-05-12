import {app, BrowserWindow, ipcMain} from 'electron';
import path from 'node:path';
import {spawn} from 'child_process';
import {startServer, storeInfo} from "./server";
import Store from 'electron-store';

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


function getBackendPath() {
    const execName = 'px';
    const isDev = !app.isPackaged;

    return isDev
        ? path.join(__dirname, '../../src-go', execName)
        : path.join(process.resourcesPath, execName);
}

function startBackend(addr: string) {
    const backendPath = getBackendPath();

    const backend = spawn(backendPath, ['-addr=' + addr], {
        stdio: ['ignore', 'pipe', 'pipe']
    });

    backend.stdout.on('data', () => {
        // console.log(`[backend stdout]: ${data}`);
    });

    backend.stderr.on('data', (data) => {
        console.error(`[backend stderr]: ${data}`);
    });

    backend.on('error', (err) => console.error('Backend error:', err));
    backend.on('exit', (code) => console.log('Backend exited with code:', code));
}

let isQuiting = false;
ipcMain.on('quit-app', () => {
    isQuiting = true;
    app.quit();
});


let mainWindow: BrowserWindow;
const createWindow = () => {
    mainWindow = new BrowserWindow({
        width: 1100,
        height: 760,
        minWidth: 960,
        minHeight: 660,
        webPreferences: {
            sandbox: false,
            disableBlinkFeatures: "Autofill",
            preload: path.join(__dirname, 'preload.js'),
            contextIsolation: true,
        },
        titleBarStyle: 'hiddenInset',
    });

    // 隐藏菜单栏
    mainWindow.setMenu(null);

    const isDev = !app.isPackaged;  // 判断是否在开发模式

    if (isDev) {
        const filePath = `http://localhost:5173?port=${storeInfo.port()}&secret=${storeInfo.secret()}`;
        mainWindow.loadURL(filePath);
    } else {
        const filePath = `http://${storeInfo.listenAddr()}/index.html?port=${storeInfo.port()}&secret=${storeInfo.secret()}`;
        console.log('准备就绪，加载窗口，url=', filePath);
        mainWindow.loadURL(filePath);
    }

    mainWindow.on('close', (event) => {
        if (!isQuiting) {
            event.preventDefault();
            mainWindow.hide();
        }
    });
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
    isQuiting = true;
    app.quit();
} else {
    app.on('second-instance', () => {
        if (mainWindow) {
            if (mainWindow.isMinimized()) mainWindow.restore();
            mainWindow.focus();
        }
    });

    startServer(resolveReady, startBackend)

    app.whenReady().then(async () => {
        await waitForReady;
        initStore(storeInfo.home())
        console.log('准备就绪，启动窗口，port=', storeInfo.port(), ' secret=', storeInfo.secret());
        createWindow();
    });
}
