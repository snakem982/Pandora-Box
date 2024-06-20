class WS {
    ip: string;
    port: string;
    url: string;
    ws: WebSocket;
    closure: Function;
    send: Function;

    constructor(
        ip: string,
        port: string,
        url: string,
        onopen: ((ws: WS, ev: Event) => any) | null = null,
        onmessage: ((ws: WS, ev: MessageEvent) => any) | null = null,
        onerror: ((ws: WS, ev: Event) => any) | null = null,
        onclose: ((ws: WS, ev: CloseEvent) => any) | null = null
    ) {
        this.ip = ip;
        this.port = port;
        this.url = url;
        this.ws = new WebSocket("ws://" + ip + ":" + port + url);
        this.closure = (): void => {
            this.ws.close();
        };
        this.send = (msg: any): void => {
            this.ws.send(msg);
        };

        const onopen_ = (ev: Event): any => {
            if (onopen !== null) {
                onopen(this, ev);
            }
            console.log(`websocket ${this.ip}:${this.port}${this.url} 连接开启！`);
        };

        const onmessage_ = (ev: MessageEvent): any => {
            if (onmessage !== null) {
                onmessage(this, ev);
            } else {
                console.log(`websocket ${this.ip}:${this.port}${this.url} 收到信息：${ev}！`);
            }
        };

        const onerror_ = (ev: Event): any => {
            if (onerror !== null) {
                onerror(this, ev);
            }
            console.log(`websocket ${this.ip}:${this.port}${this.url} 连接error：${ev}！`);
        };

        const onclose_ = (ev: CloseEvent): any => {
            if (onclose !== null) {
                onclose(this, ev);
            }
            console.log(`websocket ${this.ip}:${this.port}${this.url} 连接关闭！`);
        };

        this.ws.onopen = onopen_;
        this.ws.onmessage = onmessage_;
        this.ws.onerror = onerror_;
        this.ws.onclose = onclose_;
    }
}

export {WS};