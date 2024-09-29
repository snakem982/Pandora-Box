import axiosInstance from './axios-instance';
import {GetFreePort, GetSecret} from "../../wailsjs/go/main/App";

async function getFullUrl(url: string) {
    const baseUrl = await GetFreePort()
    return 'http://' + baseUrl + url;
}

async function getConfig(params?: any) {
    const secret = await GetSecret()
    const config = {
        headers: {
            'Authorization': 'Bearer ' + secret,
        }
    };
    Object.assign(config, params);

    return config;
}

export async function get<T>(url: string, params?: any): Promise<T> {
    const config = await getConfig(params);
    const fullUrl = await getFullUrl(url)
    const response = await axiosInstance.get<T>(fullUrl, config);
    return response.data;
}

export async function post<T>(url: string, data?: any): Promise<T> {
    const config = await getConfig();
    const fullUrl = await getFullUrl(url)
    const response = await axiosInstance.post<T>(fullUrl, data, config);
    return response.data;
}

export async function put<T>(url: string, data?: any): Promise<T> {
    const config = await getConfig();
    const fullUrl = await getFullUrl(url)
    const response = await axiosInstance.put<T>(fullUrl, data, config);
    return response.data;
}

export async function patch<T>(url: string, data?: any): Promise<T> {
    const config = await getConfig();
    const fullUrl = await getFullUrl(url)
    const response = await axiosInstance.patch<T>(fullUrl, data, config);
    return response.data;
}

export async function del<T>(url: string, params?: any): Promise<T> {
    const config = await getConfig(params);
    const fullUrl = await getFullUrl(url)
    const response = await axiosInstance.delete<T>(fullUrl, config);
    return response.data;
}
