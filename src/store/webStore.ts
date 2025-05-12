import {defineStore} from 'pinia';

export const useWebStore = defineStore('web', {
    state: () => ({
        host: '127.0.0.1', // 默认值
        port: '9686',       // 默认端口
        secret: 'Y8IUaPeFLTRvsrdf2mUJkLMBuphVZRE5',         // 默认密钥
        logs: [],         // 日志
        dnd: false,         // 拖拽显示
        dProfile: [],         // 传输文件 拖拽添加文件用
        fProfile: {}, // 更新profile 配置切换用
    }),
    getters: {
        // 确保使用 state 参数引用正确
        baseUrl: (state) => `http://${state.host}:${state.port}`,
        wsUrl: (state) => `ws://${state.host}:${state.port}`,
    },
    actions: {
        setHost(host: string) {
            if (host) this.host = host;
        },
        setPort(port: string) {
            if (port) this.port = port;
        },
        setSecret(secret: string) {
            if (secret) this.secret = secret;
        },
        addLog(log: any) {
            // 只保留最近的100条日志
            if (this.logs.length >= 100) {
                this.logs.pop();
            }
            // 在头部添加新日志
            this.logs.unshift(log);
        },
    },
    persist: true,
});
