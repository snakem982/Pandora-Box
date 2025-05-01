export const memoryCache: Record<string, string> = {};

export const customStorage = {
    getItem: (key: string): string | null => {
        return memoryCache[key] ?? null;
    },
    setItem: (key: string, value: string): void => {
        memoryCache[key] = value;
        console.log('setItem', key, value);
        window["pxSetItem"](key, value);
    }
};

export const defaultPersist = {
    enabled: true,
    storage: customStorage
};
