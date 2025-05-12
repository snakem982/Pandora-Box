// @ts-nocheck

// 内存缓存
export const memoryCache: Record<string, string> = {};

// 自定义存储（模拟同步）
export const customStorage = {
    getItem: (key: string): string | null => {
        return memoryCache[key] ?? null;
    },

    setItem: (key: string, value: string): void => {
        memoryCache[key] = value; // 先存入缓存
        if (window["pxStore"]) {
            window["pxStore"].set(key, value)
        }
    }
};

// 持久化配置
export const defaultPersist = {
    enabled: true,
    storage: customStorage
};
