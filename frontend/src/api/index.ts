import createProxiesApi from './proxies';
import createHomeApi from "@/api/home";
import createConnApi from "@/api/connections";
import createRuleApi from './rule';
import createProfilesApi from "@/api/profiles";

export default function createApi(proxy: any) {
    return {
        getDelay: createProxiesApi(proxy).getDelay,
        getGroups: createProxiesApi(proxy).getGroups,
        getProxies: createProxiesApi(proxy).getProxies,
        setProxy: createProxiesApi(proxy).setProxy,
        getVersion: createHomeApi(proxy).getVersion,
        getConfigs: createHomeApi(proxy).getConfigs,
        getGroupMd5: createHomeApi(proxy).getGroupMd5,
        updateConfigs: createHomeApi(proxy).updateConfigs,
        closeConnection: createConnApi(proxy).closeConnection,
        closeAllConnection: createConnApi(proxy).closeAllConnection,
        getRules: createRuleApi(proxy).getRules,
        addProfileFromInput: createProfilesApi(proxy).addProfileFromInput,
        addProfileFromFile: createProfilesApi(proxy).addProfileFromFile,
        deleteProfile: createProfilesApi(proxy).deleteProfile,
        updateProfile: createProfilesApi(proxy).updateProfile,
        getProfileList: createProfilesApi(proxy).getProfileList,
        refreshProfile: createProfilesApi(proxy).refreshProfile,
        switchProfile: createProfilesApi(proxy).switchProfile,
    };
}
