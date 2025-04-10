export class WS {
    url: string;
    ws: WebSocket;
    closure: Function;
    send: Function;

    constructor(
        url: string,
        onopen: ((ev: Event) => any) | null = null,
        onmessage: ((ev: MessageEvent) => any) | null = null,
        onerror: ((ev: Event) => any) | null = null,
        onclose: ((ev: CloseEvent) => any) | null = null
    ) {
        this.url = url;
        this.ws = new WebSocket(url);
        this.closure = (): void => {
            this.ws.close();
        };
        this.send = (msg: any): void => {
            this.ws.send(msg);
        };

        // 绑定事件
        this.ws.onopen = (ev: Event) => {
            onopen?.(ev);
            console.log(`websocket ${this.url} 连接开启！`);
        };

        this.ws.onmessage = (ev: MessageEvent) => {
            onmessage?.(ev);
        };

        this.ws.onerror = (ev: Event) => {
            onerror?.(ev);
            console.log(`websocket ${this.url} 连接发生错误：`, ev);
        };

        this.ws.onclose = (ev: CloseEvent) => {
            onclose?.(ev);
            console.log(`websocket ${this.url} 连接关闭！`);
            this.ws.onmessage = null; // 清除监听
        };
    }

    // 强制清理
    close() {
        if (
            this.ws.readyState === WebSocket.OPEN ||
            this.ws.readyState === WebSocket.CONNECTING
        ) {
            this.ws.close();
        }
        this.ws.onmessage = null; // 移除监听
        this.ws.onerror = null;
        this.ws.onclose = null;
        console.log("WebSocket 已完全清理");
    }
}
