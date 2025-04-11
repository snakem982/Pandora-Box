import createProxiesApi from './proxies';
import createHomeApi from "@/api/home";
import createConnApi from "@/api/connections";

export default function createApi(proxy: any) {
    return {
        getDelay: createProxiesApi(proxy).getDelay,
        getGroups: createProxiesApi(proxy).getGroups,
        getProxies: createProxiesApi(proxy).getProxies,
        setProxy: createProxiesApi(proxy).setProxy,
        getVersion: createHomeApi(proxy).getVersion,
        getConfigs: createHomeApi(proxy).getConfigs,
        getGroupMd5: createHomeApi(proxy).getGroupMd5,
        closeConnection: createConnApi(proxy).closeConnection,
        closeAllConnection: createConnApi(proxy).closeAllConnection,
    };
}
