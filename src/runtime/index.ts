// 只监听 electron 的消息
export const Events = {
    // 只向electron发消息
    Emit: ({name, data}: { name: string; data: any }) => {
        window.pxCommon?.emit(name, data)
        console.log("emit========", name)
    },
    // 只收electron的消息
    On: (name: string, callback: (...args: any[]) => void) => {
        window.pxCommon?.on(name, callback);
        console.log("on========", name)
    },
};

// 获取剪贴板内容
export const Clipboard = {
    Text: async (): Promise<string> => {
        if (typeof window.pxClipboard === 'function') {
            return window.pxClipboard();
        }
        if (navigator.clipboard?.readText) {
            return navigator.clipboard.readText();
        }
        throw new Error('Clipboard access is unavailable');
    }
}

// 打开地址
export const Browser = {
    // @ts-ignore
    OpenURL: (url: string) => window["pxOpen"](url)
}