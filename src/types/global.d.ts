import {AxiosRequest} from "@/util/axiosRequest";


// 为 '@/api' 模块提供类型声明
declare module '@/api' {
    export interface Api {
        // 使用 unknown 替代 any，符合 ESLint 规范，同时支持传入泛型指定具体返回类型
        proxies: <T = unknown>() => Promise<T>;
    }
}

// 为 Vue 的全局属性添加类型声明
declare module '@vue/runtime-core' {
    export interface ComponentCustomProperties {
        $http: AxiosRequest; // 声明全局 $http 的类型
        $t: (key: string) => string; // i18n
    }
}

// 绑定函数
declare global {
    interface Window {
        pxOs?: () => string;
        pxClipboard?: () => Promise<string>;
        pxDeepLink?: {
            onImportProfile: (callback: (data: {
                rawUrl?: string;
                url?: string;
                name?: string
            } | string) => void) => void;
            notifyReady?: () => void | Promise<void>;
        };
        pxStore?: {
            // 1. 使用 unknown 替代 any：安全且符合规范
            // 2. 结合泛型 <T = unknown>：调用时可以手动指定返回类型，不指定时默认为 unknown
            get: <T = unknown>(key: string) => Promise<T>;
            set: (key: string, value: unknown) => Promise<void> | void;
            [key: string]: unknown;
        };
        pxCommon?: {
            emit: (name: string, data: unknown) => void;
            on: (name: string, callback: (...args: unknown[]) => void) => void;
            [key: string]: unknown; // 允许存在其他扩展属性或方法
        };
    }
}