import axiosInstance from './axios-instance';
import {GetFreePort} from "../../wailsjs/go/main/App";

export async function get<T>(url: string, params?: any): Promise<T> {
    const baseUrl = await GetFreePort()
    const fullUrl = 'http://'+baseUrl + url
    const response = await axiosInstance.get<T>(fullUrl, params);
    return response.data;
}

export async function post<T>(url: string, data?: any): Promise<T> {
    const baseUrl = await GetFreePort()
    const fullUrl = 'http://'+baseUrl + url
    const response = await axiosInstance.post<T>(fullUrl, data);
    return response.data;
}

export async function put<T>(url: string, data?: any): Promise<T> {
    const baseUrl = await GetFreePort()
    const fullUrl = 'http://'+baseUrl + url
    const response = await axiosInstance.put<T>(fullUrl, data);
    return response.data;
}

export async function patch<T>(url: string, data?: any): Promise<T> {
    const baseUrl = await GetFreePort()
    const fullUrl = 'http://'+baseUrl + url
    const response = await axiosInstance.patch<T>(fullUrl, data);
    return response.data;
}

export async function del<T>(url: string, params?: any): Promise<T> {
    const baseUrl = await GetFreePort()
    const fullUrl = 'http://'+baseUrl + url
    const response = await axiosInstance.delete<T>(fullUrl, {params});
    return response.data;
}
