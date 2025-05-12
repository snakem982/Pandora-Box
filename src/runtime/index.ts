// 只监听 webview 的消息
export const Events = {
    // 只向webview发消息
    Emit: ({name, data}: { name: string; data: any }) => {
        // @ts-ignore
        if (window["px_" + name]) {
            (window as any)["px_" + name](data);
        }
    },
    // 只收webview的消息
    On: (name: string, callback: (...args: any[]) => void) => {
        // @ts-ignore
        window["px_" + name] = callback;
    },
};

// 获取剪贴板内容
export const Clipboard = {
    // @ts-ignore
    Text: window["pxClipboard"]
}

// 打开地址
export const Browser = {
    // @ts-ignore
    OpenURL: (url: string) => window["pxOpen"](url)
}