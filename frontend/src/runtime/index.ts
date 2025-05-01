// 只监听 webview 的消息
export const Events = {
    // 只向webview发消息
    Emit: ({name, data}: { name: string; data: any }) => {
        if (window["px_" + name]) {
            (window as any)["px_" + name](data);
        }
    },
    // 只收webview的消息
    On: (name: string, callback: (...args: any[]) => void) => {
        console.log('on');
        window["px_" + name] = callback;
    },
};

// 获取剪贴板内容
export const Clipboard = {
    Text: () => {
        if (window["pxClipboard"]) {
            return (window as any)["pxClipboard"]();
        }

        return ""
    }
}

// 打开地址
export const Browser = {
    OpenURL: (url: string) => {
        console.log('browser');
        if (window["pxOpen"]) {
            return (window as any)["pxOpen"](url);
        }
    }
}