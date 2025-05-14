// @ts-nocheck

import {clipboard, contextBridge, ipcRenderer, shell} from 'electron';
import os from 'os';

// tray相关
contextBridge.exposeInMainWorld('pxTray', {
    on: (name, callback) => ipcRenderer.on('px_' + name, (_event, ...value: any[]) => callback(...value)),
    emit: (name: string, ...value: any[]) => ipcRenderer.send('px_' + name, ...value)
})

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

// 打开配置目录
contextBridge.exposeInMainWorld('pxConfigDir', (url: string) => shell.openPath(url));

// 获取剪贴板内容
contextBridge.exposeInMainWorld('pxClipboard', () => clipboard.readText());

// 打开外部URL地址
contextBridge.exposeInMainWorld('pxOpen', (url: string) => shell.openExternal(url));
