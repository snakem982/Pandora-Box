import {app, BrowserWindow, globalShortcut, Menu, nativeImage, Tray} from 'electron';
import path from "node:path";
import {storeSet} from "./store";
import {disableAutoLaunch, enableAutoLaunch} from "./launch";
import {doChange} from "./change";
import {AppName, IsDev, onMsg, sendMsg} from "./common";

// --- 状态管理 ---
let tray: Tray = null;
let mainWindow: BrowserWindow = null;
let isQuiting = false;

// 存储菜单所需的动态数据
const state = {
    mode: 'rule',
    proxy: false,
    tun: false,
    profiles: [] as any[],
    labels: {
        show: '显示窗口',
        rule: '规则',
        global: '全局',
        direct: '直连',
        profiles: '订阅',
        proxy: '系统代理',
        tun: 'Tun模式',
        quit: '退出'
    }
};

// --- 核心逻辑 ---

// 切换规则
function switchMode(menuItem: any, mode: string) {
    if (!menuItem.checked) {
        menuItem.checked = true
        return
    }
    sendMsg(mainWindow, "switchMode", mode)
}


// 切换配置
function switchProfiles(menuItem: any, profile: any) {
    if (!menuItem.checked) {
        menuItem.checked = true
        return
    }
    sendMsg(mainWindow, "switchProfiles", profile)
}

/**
 * 根据当前 state 重新渲染托盘菜单
 */
function updateTrayMenu() {
    if (!tray) return;

    const template: any[] = [
        {label: state.labels.show, click: () => showWindow()},
        {type: 'separator'},
        {
            label: state.labels.rule,
            type: 'checkbox',
            checked: state.mode === 'rule',
            click: (menuItem: any) => switchMode(menuItem, 'rule')
        },
        {
            label: state.labels.global,
            type: 'checkbox',
            checked: state.mode === 'global',
            click: (menuItem: any) => switchMode(menuItem, 'global')
        },
        {
            label: state.labels.direct,
            type: 'checkbox',
            checked: state.mode === 'direct',
            click: (menuItem: any) => switchMode(menuItem, 'direct')
        },
        {type: 'separator'},
        {
            label: state.labels.profiles,
            submenu: state.profiles.map(p => ({
                label: p.title,
                type: 'checkbox',
                checked: !!p.selected,
                click: (menuItem: any) => switchProfiles(menuItem, p)
            }))
        },
        {type: 'separator'},
        {
            label: state.labels.proxy,
            type: 'checkbox',
            checked: state.proxy,
            click: () => sendMsg(mainWindow, "switchProxy")
        },
        {label: state.labels.tun, type: 'checkbox', checked: state.tun, click: () => sendMsg(mainWindow, "switchTun")},
        {type: 'separator'},
        {label: state.labels.quit, click: () => sendMsg(mainWindow, "readyToQuit")},
    ];

    tray.setContextMenu(Menu.buildFromTemplate(template));
}

/**
 * 初始化托盘
 */
export function initTray(win: BrowserWindow) {
    mainWindow = win;

    // 窗口关闭逻辑
    mainWindow.on('close', (event) => {
        if (!isQuiting) {
            event.preventDefault();
            process.platform === 'darwin' ? mainWindow?.hide() : mainWindow?.minimize();
        }
    });

    // 图标处理
    const trayPath = IsDev
        ? path.join(__dirname, '../../public', 'tray.png')
        : path.join(process.resourcesPath, 'tray.png');
    const size = process.platform === 'darwin' ? 16 : 32;
    const trayImage = nativeImage.createFromPath(trayPath).resize({width: size, height: size});

    tray = new Tray(trayImage);
    tray.setToolTip(AppName);

    // 初次渲染菜单
    updateTrayMenu();

    // 点击事件
    tray.on('click', () => tray.popUpContextMenu());
    tray.on('right-click', () => tray.popUpContextMenu());
}

// --- 工具函数 ---

export function showWindow() {
    mainWindow?.show();
    app.dock?.show();
    mainWindow?.focus();
}


export const doQuit = () => {
    if (mainWindow) storeSet('windowBounds', mainWindow.getBounds());
    if (app.isReady()) {
        try {
            globalShortcut.unregisterAll();
        } catch (error) {
            console.error("Failed to unregisterAll shortcut:", error);
        }
    }
    isQuiting = true;
    app.quit();
}

onMsg("doQuit", doQuit);

onMsg("mode", (val) => {
    state.mode = val;
    updateTrayMenu();
});
onMsg("proxy", (val) => {
    state.proxy = val;
    updateTrayMenu();
});
onMsg("tun", (val) => {
    state.tun = val;
    updateTrayMenu();
});
onMsg("profiles", (val) => {
    state.profiles = val;
    updateTrayMenu();
});

onMsg("translate", (labels) => {
    // 批量更新 label 名称
    Object.keys(labels).forEach(key => {
        const pureKey = key.replace('tray.', '');
        // @ts-ignore
        if (state.labels[pureKey]) state.labels[pureKey] = labels[key];
    });
    updateTrayMenu();
});

onMsg("hide", () => {
    mainWindow?.hide();
    app.dock?.hide();
});
onMsg("close", () => app.quit());
onMsg("max", () => mainWindow.isMaximized() ? mainWindow.unmaximize() : mainWindow.maximize());
onMsg("min", () => mainWindow.minimize());
onMsg("boot", (val) => val ? enableAutoLaunch() : disableAutoLaunch());
onMsg("tunAuthTip", (val) => {
        if (val) storeSet("tunAuthTip", val)
    }
)
onMsg("doChangeConfigDir", (val) => doChange(mainWindow!, val));