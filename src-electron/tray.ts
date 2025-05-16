// @ts-nocheck

import {app, BrowserWindow, ipcMain, Menu, nativeImage, Tray} from 'electron';
import path from "node:path";

// 是否在开发模式
const isDev = !app.isPackaged;

// 退出app
let isQuiting = false;
export const doQuit = () => {
    isQuiting = true;
    app.quit();
}
const readyToQuit = () => emitWindow("readyToQuit");
onWindow("doQuit", doQuit)

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
                            click: () => {
                                // 获取当前聚焦的窗口
                                const win = BrowserWindow.getFocusedWindow();
                                if (win) win.webContents.openDevTools();
                            }
                        },
                        {
                            label: 'Reload',
                            click: () => {
                                const win = BrowserWindow.getFocusedWindow();
                                if (win) win.webContents.reload();
                            }
                        },
                        {
                            label: 'Force Reload',
                            click: () => {
                                const win = BrowserWindow.getFocusedWindow();
                                if (win) win.webContents.reloadIgnoringCache();
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

const initMenu = () => createMenu([
    {
        label: 'Pandora-Box', submenu: [
            {
                label: 'Quit', accelerator: 'Cmd+Q', click: readyToQuit
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


// 托盘
let tray: Tray;
// 托盘菜单
let currentMenu: any
// 当前窗口
let mainWindow: BrowserWindow

// 显示窗口
function showWindow() {
    mainWindow.show();
    app.dock?.show();
}

// 切换规则
function switchMode(menuItem, mode) {
    if (!menuItem.checked) {
        menuItem.checked = true
        return
    }
    emitWindow("switchMode", mode);
}

// 切换配置
function switchProfiles(menuItem, profile) {
    if (!menuItem.checked) {
        menuItem.checked = true
        return
    }
    emitWindow("switchProfiles", profile);
}

const trayMap: Map<any, any> = new Map();
trayMap.set('tray.show', {
    id: 'tray.show',
    label: '显示窗口',
    type: 'normal',
    click: showWindow
});
trayMap.set('tray.rule', {
    id: 'tray.rule',
    label: '规则',
    type: 'checkbox',
    checked: false,
    click: (menuItem) => switchMode(menuItem, 'rule')
});
trayMap.set('tray.global', {
    id: 'tray.global',
    label: '全局',
    type: 'checkbox',
    checked: false,
    click: (menuItem) => switchMode(menuItem, 'global')
});
trayMap.set('tray.direct', {
    id: 'tray.direct',
    label: '直连',
    type: 'checkbox',
    checked: false,
    click: (menuItem) => switchMode(menuItem, 'direct')
});
trayMap.set('tray.profiles', {id: 'tray.profiles', label: '订阅', submenu: []});
trayMap.set('tray.proxy', {
    id: 'tray.proxy',
    label: '系统代理',
    type: 'checkbox',
    checked: false,
    click: () => emitWindow("switchProxy")
});
trayMap.set('tray.tun', {
    id: 'tray.tun',
    label: 'Tun模式',
    type: 'checkbox',
    checked: false,
    click: () => emitWindow("switchTun")
});
trayMap.set('tray.quit', {id: 'tray.quit', label: '退出', type: 'normal', click: readyToQuit});

const createTrayMenu = () => [
    trayMap.get('tray.show'),
    {type: 'separator'},
    trayMap.get('tray.rule'),
    trayMap.get('tray.global'),
    trayMap.get('tray.direct'),
    {type: 'separator'},
    trayMap.get('tray.profiles'),
    {type: 'separator'},
    trayMap.get('tray.proxy'),
    trayMap.get('tray.tun'),
    {type: 'separator'},
    trayMap.get('tray.quit'),
]

// 初始化托盘菜单
currentMenu = Menu.buildFromTemplate(createTrayMenu());

// 初始化托盘
export function initTray(browserWindow: BrowserWindow): void {
    // 初始化左上角菜单
    initMenu()

    // 初始化窗口事件
    mainWindow = browserWindow
    mainWindow.on('close', (event) => {
        if (!isQuiting) {
            event.preventDefault();
            if (process.platform !== 'darwin') {
                mainWindow.minimize()
            } else {
                mainWindow.hide();
            }
        }
    });

    // 初始化tray
    let trayImage: any;
    if (process.platform === 'darwin') {
        trayImage = nativeImage.createFromPath(path.join(__dirname, 'tray.png')).resize({width: 16, height: 16});
    } else {
        trayImage = nativeImage.createFromPath(path.join(__dirname, 'tray.png')).resize({width: 32, height: 32});
    }
    tray = new Tray(trayImage);
    tray.setToolTip('Pandora-Box');
    tray.setContextMenu(Menu.buildFromTemplate(createTrayMenu()))
}


// 接收浏览器消息
function onWindow(name, cb) {
    ipcMain.on('px_' + name, (_event, value) => {
        if (cb) {
            cb(value)
        }
    })
}

// 发送消息到浏览器
function emitWindow(name: string, ...value: any[]) {
    if (mainWindow) {
        mainWindow.webContents.send('px_' + name, ...value);
    }
}


// 监听消息
onWindow("translate", function (trayOptions) {
    for (const [key, value] of Object.entries(trayOptions)) {
        trayMap.get(key).label = value
    }
    currentMenu = Menu.buildFromTemplate(createTrayMenu());
    tray.setContextMenu(currentMenu);
})
onWindow("mode", function (value) {
    currentMenu.getMenuItemById('tray.rule').checked = false;
    currentMenu.getMenuItemById('tray.global').checked = false;
    currentMenu.getMenuItemById('tray.direct').checked = false;
    trayMap.get('tray.rule').checked = false;
    trayMap.get('tray.global').checked = false;
    trayMap.get('tray.direct').checked = false;
    const key = 'tray.' + value
    currentMenu.getMenuItemById(key).checked = true;
    trayMap.get(key).checked = true
})
onWindow("proxy", function (value) {
    const key = 'tray.proxy'
    currentMenu.getMenuItemById(key).checked = value;
    trayMap.get(key).checked = value
})
onWindow("tun", function (value) {
    const key = 'tray.tun'
    currentMenu.getMenuItemById(key).checked = value;
    trayMap.get(key).checked = value
})
onWindow("profiles", function (profiles) {
    const key = 'tray.profiles'
    const pList = []
    for (let profile of profiles) {
        pList.push({
            label: profile.title,
            type: 'checkbox',
            checked: !!profile.selected,
            click: (menuItem) => switchProfiles(menuItem, profile)
        })
    }
    trayMap.get(key).submenu = pList
    currentMenu = Menu.buildFromTemplate(createTrayMenu());
    tray.setContextMenu(currentMenu);
})

// 窗口控制
onWindow("hide", function () {
    mainWindow.hide();
    app.dock?.hide()
})
onWindow("close", function () {
    app.quit()
})
onWindow("max", function () {
    mainWindow.isMaximized() ? mainWindow.unmaximize() : mainWindow.maximize()
})
onWindow("min", function () {
    mainWindow.minimize()
})





