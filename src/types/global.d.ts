import {AxiosRequest} from "@/util/axiosRequest";

// 为 '@/api' 模块提供类型声明
declare module '@/api' {
    // 定义一个类型接口，替代 `any`，根据项目实际的 API 结构定义
    export interface Api {
        proxies: () => Promise<any>;
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
        pxOs: () => string;
    }
}

