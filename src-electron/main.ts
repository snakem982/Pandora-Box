import {app, BrowserWindow, ipcMain, Menu, nativeImage, Tray} from 'electron';
import path from 'node:path';
import {spawn} from 'child_process';
import {startServer, storeInfo} from "./server";
import Store from 'electron-store';

// 是否在开发模式
const isDev = !app.isPackaged;

// 托盘
let tray: Electron.CrossProcessExports.Tray;

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


// 处理菜单
const createMenu = (menuTemplate: any) => {
    if (process.platform === 'darwin') {
        if (isDev) {
            menuTemplate.push(
                {
                    label: 'View',
                    submenu: [
                        {
                            label: 'Open Developer Tools',
                            accelerator: 'CmdOrCtrl+Shift+I',
                            click: () => {
                                // 获取当前聚焦的窗口
                                const win = BrowserWindow.getFocusedWindow();
                                if (win) win.webContents.openDevTools();
                            }
                        }
                    ]
                }
            )
        }
        const menu = Menu.buildFromTemplate(menuTemplate);
        Menu.setApplicationMenu(menu);
    }
};

ipcMain.on('update-menu', (event, menuTemplate) => {
    createMenu(menuTemplate);
});

createMenu([
    {
        label: 'Pandora-Box', submenu: [
            {
                label: 'Quit', accelerator: 'Cmd+Q', click: () => {
                    isQuiting = true;
                    app.quit();
                }
            }
        ]
    },
    {
        label: 'Edit',
        submenu: [
            {label: 'Undo', role: 'undo'},
            {label: 'Redo', role: 'redo'},
            {type: 'separator'},
            {label: 'Cut', role: 'cut'},
            {label: 'Copy', role: 'copy'},
            {label: 'Paste', role: 'paste'},
            {label: 'Delete', role: 'delete'},
            {type: 'separator'},
            {label: 'Select All', role: 'selectAll'}
        ]
    }
]);

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
        // 设置托盘
        const trayImage = nativeImage.createFromPath(path.join(__dirname, 'tray.png')).resize({width: 16, height: 16});
        tray = new Tray(trayImage);
        tray.setToolTip('Pandora-Box');

        // 等待后端启动
        await waitForReady;
        initStore(storeInfo.home())

        // 启动UI
        console.log('准备就绪，启动窗口，port=', storeInfo.port(), ' secret=', storeInfo.secret());
        createWindow();
    });
}

let currentMenu: any
ipcMain.on('update-tray', (event, newMenu) => {
    if (tray) {
        // 释放旧菜单对象
        if (currentMenu) {
            currentMenu = null;
        }
        // 更新菜单
        currentMenu = Menu.buildFromTemplate(newMenu);
        tray.setContextMenu(currentMenu);
    }
});