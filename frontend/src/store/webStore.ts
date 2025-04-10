import { defineStore } from 'pinia';

export const useWebStore = defineStore('useWebStore', {
    state: () => ({
        host: '127.0.0.1', // 默认值
        port: '9966',       // 默认端口
        secret: '',         // 默认密钥
    }),
    getters: {
        // 确保使用 state 参数引用正确
        baseUrl: (state) => `http://${state.host}:${state.port}`,
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
    },
    persist: {
        enabled: true,
        storage: sessionStorage,
    },
});
