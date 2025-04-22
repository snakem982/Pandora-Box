import axios, {AxiosInstance, AxiosResponse} from 'axios';

// axios的封装
export class AxiosRequest {

    private instance: AxiosInstance;

    // 默认配置baseURL等
    constructor(baseURL: string, secret: string = "", timeout = 30000) {
        this.instance = axios.create({
            baseURL,
            timeout,
            headers: {
                'Content-Type': 'application/json',
                'Authorization': 'Bearer ' + secret
            }
        })

        // 添加请求拦截器
        this.instance.interceptors.request.use(
            (config: any) => {
                // 在发送请求之前做些什么
                return config;
            },
            (error: any) => {
                // 处理请求错误
                // @ts-ignore
                return Promise.reject(error);
            },
        );

        // 添加响应拦截器
        this.instance.interceptors.response.use(
            (response: AxiosResponse) => {
                // 对响应数据做处理，只有状态码为 200 或 204 时返回数据
                if (response.status === 200) {
                    return response.data;
                }
            },
            (e: any) => {
                if (e['response'] && e['response']['data']) {
                    return Promise.reject(e['response']['data']);
                }
                // 处理响应错误
                // @ts-ignore
                return Promise.reject(e);
            },
        );
    }

    get<T>(url: string, params?: any): Promise<AxiosResponse<T, any>> {
        return this.instance.get<T>(url, params);
    }

    post<T>(url: string, params?: any): Promise<AxiosResponse<T, any>> {
        return this.instance.post<T>(url, params);
    }

    put<T>(url: string, params?: any): Promise<AxiosResponse<T, any>> {
        return this.instance.put<T>(url, params);
    }

    patch<T>(url: string, params?: any): Promise<AxiosResponse<T, any>> {
        return this.instance.patch<T>(url, params);
    }

    delete<T>(url: string, params?: any): Promise<AxiosResponse<T, any>> {
        return this.instance.delete<T>(url, params);
    }
}