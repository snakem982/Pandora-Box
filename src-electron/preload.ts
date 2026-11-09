// @ts-nocheck

import {contextBridge, ipcRenderer, shell} from 'electron';
import os from 'os';

// tray相关
contextBridge.exposeInMainWorld('pxCommon', {
    on: (name, callback) => {
        const eventName = 'px_' + name;
        // 移除旧监听器，确保只注册一次
        ipcRenderer.removeAllListeners(eventName);
        ipcRenderer.on(eventName, (_event, ...value) => callback(...value));
    },
    emit: (name, ...value) => ipcRenderer.send('px_' + name, ...value)
});

// 深度链接相关
contextBridge.exposeInMainWorld('pxDeepLink', {
    onImportProfile: (callback) => {
        // 移除旧监听器，确保只注册一次
        ipcRenderer.removeAllListeners('import-profile-from-deeplink');
        ipcRenderer.on('import-profile-from-deeplink', (_event, data) => callback(data));
    },
    notifyReady: () => {
        ipcRenderer.send('deeplink-handler-ready');
    }
});


// 缓存接口
contextBridge.exposeInMainWorld('pxStore', {
    get: (key) => ipcRenderer.invoke('store:get', key),
    set: (key, value) => ipcRenderer.invoke('store:set', key, value)
});

// 获取系统信息
contextBridge.exposeInMainWorld('pxOs', () => {
    switch (os.type()) {
        case 'Darwin':
            return "MacOS " + os.arch()
        case 'Linux':
            return "Linux " + os.arch()
        case 'Windows_NT':
            return "Windows " + os.arch()
        default:
            return "Unknown";
    }
});

// 打开配置目录 changePxConfigDir
contextBridge.exposeInMainWorld('openPxConfigDir', (url: string) => shell.openPath(url));
contextBridge.exposeInMainWorld('changePxConfigDir', (options = {}) => ipcRenderer.invoke('select-directory', options));

// 获取剪贴板内容
contextBridge.exposeInMainWorld('pxClipboard', () => ipcRenderer.invoke('clipboard:readText'));

// 打开外部URL地址
contextBridge.exposeInMainWorld('pxOpen', (url: string) => shell.openExternal(url));

// 控制标题栏
if (process.platform !== 'darwin') {
    contextBridge.exposeInMainWorld('pxShowBar', () => console.log('pxShowBar'));
}