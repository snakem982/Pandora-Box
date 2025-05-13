// 只监听 electron 的消息
export const Events = {
    // 只向electron发消息
    Emit: ({name, data}: { name: string; data: any }) => {
        // @ts-ignore
        window.pxTray.emit(name, data)
        console.log("emit========", name)
    },
    // 只收electron的消息
    On: (name: string, callback: (...args: any[]) => void) => {
        // @ts-ignore
        window.pxTray.on(name, callback);
        console.log("on========", name)
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