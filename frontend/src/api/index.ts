import createProxiesApi from './proxies';

export default function createApi(proxy: any) {
    return {
        getDelay: createProxiesApi(proxy).getDelay,
        getGroups: createProxiesApi(proxy).getGroups,
        getProxies: createProxiesApi(proxy).getProxies,
        setProxy: createProxiesApi(proxy).setProxy,
    };
}
