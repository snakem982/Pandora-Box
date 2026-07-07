import {app, BrowserWindow, globalShortcut} from 'electron';
import {onMsg, sendMsg} from "./common";

const scs = {
    "showOrHide": showOrHide
};

function replaceCtrlCmd(key: string) {
    return key.replaceAll("Ctrl", "CommandOrControl").replaceAll("Cmd", "CommandOrControl");
}

// 修改为异步函数，处理潜在的竞态条件
async function registerOne(acc: any, mainWindow: BrowserWindow | null): Promise<void> {
    // 即使主进程逻辑里写了 whenReady，但在 IPC 回调里再次确认可以防止“抢跑”
    if (!app.isReady()) {
        await app.whenReady();
    }

    if (!scs[acc.name as keyof typeof scs]) {
        return;
    }

    const key = replaceCtrlCmd(acc.key);

    try {
        // 先注销旧的，防止重复注册导致的冲突[cite: 1]
        if (acc.old) {
            globalShortcut.unregister(replaceCtrlCmd(acc.old));
        }
        globalShortcut.unregister(key);

        // 执行注册
        const success = globalShortcut.register(key, () => {
            // @ts-ignore
            scs[acc.name](mainWindow);
        });

        sendMsg(mainWindow, "shortcut:result", success);
    } catch (error) {
        console.error("Failed to register shortcut:", error);
        sendMsg(mainWindow, "shortcut:result", false);
    }
}

export function initShortcut(mainWindow: BrowserWindow) {
    // 监听改为异步处理[cite: 1]
    onMsg('shortcut:register', async (acc) => {
        await registerOne(acc, mainWindow);
    });

    onMsg('shortcut:unregister-all', async () => {
        if (app.isReady()) {
            try {
                globalShortcut.unregisterAll();
            } catch (error) {
                console.error("Failed to unregisterAll shortcut:", error);
            }
        }
        return true;
    });
}

function showOrHide(mainWindow: BrowserWindow) {
    if (!mainWindow || mainWindow.isDestroyed()) return;

    if (!mainWindow.isVisible() || !mainWindow.isFocused()) {
        mainWindow?.show();
        mainWindow?.focus();
        app.dock?.show();
    }
}