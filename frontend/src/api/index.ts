import createProxiesApi from './proxies';

export default function createApi(proxy: any) {
    return {
        getGroups: createProxiesApi(proxy).getGroups,
        getProxies: createProxiesApi(proxy).getProxies,
        setProxy: createProxiesApi(proxy).setProxy,
    };
}
