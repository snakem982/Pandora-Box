import {app, BrowserWindow, BrowserWindowConstructorOptions, ipcMain, session} from 'electron';
import path from 'node:path';
import {startServer, storeInfo} from "./server";
import {doQuit, initTray, showWindow} from "./tray";
import {initMenu} from "./menu"
import {startBackend} from "./admin";
import log from './log';
import {initStore, storeGet} from "./store";
import {isBootAutoLaunch, updateAutoLaunchRegistration, waitForNetworkReady} from "./launch";
import {deeplink} from "./deeplink";
import {setRandomUA} from "./ua";
import {selectDirectory} from "./selector";
import {IsDev} from "./common";
import {initShortcut} from "./shortcut";

// 初始化前端数据库
initStore()

// 初始化日志
log.initLog()

// 主窗口
let mainWindow: BrowserWindow;

// 屏蔽安全警告
process.env["ELECTRON_DISABLE_SECURITY_WARNINGS"] = "true";
const createWindow = (isBoot: boolean) => {
    let windowOptions: BrowserWindowConstructorOptions = {
        minWidth: 960,
        minHeight: 660,
        width: 1100,
        height: 760,
        show: false, // 先不显示窗口
        center: true,
        webPreferences: {
            preload: path.join(__dirname, 'preload.js'),
            contextIsolation: true,
            webSecurity: false,
            nodeIntegrationInWorker: true
        },
        ...(process.platform !== 'darwin' ? {
            titleBarStyle: 'hidden'
        } : {
            titleBarStyle: 'hiddenInset'
        })
    };

    // 恢复上次窗口位置
    const savedBounds: any = storeGet('windowBounds');
    if (savedBounds && savedBounds.x !== undefined && savedBounds.y !== undefined) {
        windowOptions = {
            ...windowOptions,
            ...savedBounds
        };
    }

    mainWindow = new BrowserWindow(windowOptions);

    // 隐藏窗口菜单栏
    mainWindow.setMenu(null);

    // 顶部菜单
    initMenu(mainWindow);
    // 托盘
    initTray(mainWindow);

    // 页面加载
    const filePath = IsDev
        ? `http://localhost:5173?port=${storeInfo.port()}&secret=${storeInfo.secret()}`
        : `http://${storeInfo.listenAddr()}/index.html?port=${storeInfo.port()}&secret=${storeInfo.secret()}`;

    log.info('准备加载页面');
    mainWindow
        .loadURL(filePath)
        .then(() => {
            initShortcut(mainWindow);
        })
        .catch((err) => {
            log.error('加载页面失败:', err);
        });

    // 页面加载完成再显示，避免白屏
    mainWindow.webContents.once('did-finish-load', () => {
        if (isBoot) {
            log.info('静默启动完成');
        } else {
            mainWindow.show();
            mainWindow.focus();
            log.info('页面加载成功');
        }
    });

    mainWindow.on('closed', () => {
        deeplink.DeepLinkHandlerReady(false);
        mainWindow = null;
    });
};


// 监听深度链接事件
ipcMain.on(deeplink.DEEP_LINK_READY_EVENT, (event) => {
    if (!mainWindow || event.sender !== mainWindow.webContents) {
        return;
    }
    deeplink.DeepLinkHandlerReady(true);
    deeplink.processPendingDeepLinks(mainWindow);
});
for (const arg of process.argv) {
    deeplink.pushDeeplink(arg)
}

// 修改配置目录
ipcMain.handle('select-directory', async (event, options = {}) => {
    const win = BrowserWindow.fromWebContents(event.sender);
    const paths = await selectDirectory(win, options);
    return paths.length > 0 ? paths : null; // 返回 null 表示取消
});

// ⚠️ 必须放在 app.whenReady() 之前！否则开关不会生效 单位为字节 (200MB)
app.commandLine.appendSwitch('disk-cache-size', '209715200');

// 单例模式
const gotTheLock = app.requestSingleInstanceLock();
if (!gotTheLock) {
    doQuit()
} else {
    // 试图启动第二个应用实例
    app.on('second-instance', (_event, commandLine) => {
        showWindow();
        const urls = commandLine.filter(deeplink.isDeepLinkUrl);
        if (urls.length > 0) {
            urls.forEach((url) => deeplink.handleDeepLink(url, mainWindow));
        }
    });

    // 监听应用被激活
    app.on('activate', showWindow);

    if (process.platform === 'darwin') {
        app.on('open-url', (event, url) => {
            event.preventDefault();
            showWindow();
            deeplink.handleDeepLink(url, mainWindow);
        });
    }

    app.whenReady().then(async () => {
        // 判断是否开机启动
        const isBoot = await isBootAutoLaunch();
        log.info('是否开机启动:', isBoot);

        // 如果是开机启动，则等待网络就绪（最多30秒）
        if (isBoot) {
            // 先隐藏dock
            app.dock?.hide()

            log.info('开机启动，等待网络准备...');
            const networkReady = await waitForNetworkReady(30000, 'bing.com');
            if (!networkReady) {
                log.warn('网络检测超时，继续启动但可能无网络');
            } else {
                log.info('网络已准备好');
            }
        }

        // 等待 backend 传来的 port 和 secret
        let resolveReady: () => void;
        const waitForReady = new Promise<void>((resolve) => {
            resolveReady = resolve;
        });

        // 启动前端静态服务
        startServer(resolveReady, startBackend)

        // 等待后端启动
        await waitForReady;

        // 注册深度链接
        deeplink.registerDeepLink();

        // 设置请求头 Referer
        setRandomUA(session)

        // 启动UI
        log.info('准备就绪，启动窗口，port=', storeInfo.port(), ' secret=', storeInfo.secret());
        createWindow(isBoot);

        // 更新开机自启路径
        await updateAutoLaunchRegistration()
    });
}