import createProxiesApi from "./proxies";
import createHomeApi from "@/api/home";
import createConnApi from "@/api/connections";
import createRuleApi from "./rule";
import createProfilesApi from "@/api/profiles";
import createDnsApi from "@/api/dns";
import createMihomoApi from "@/api/mihomo";
import createPandoraApi from "@/api/pandora";

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
        getWebTest: createHomeApi(proxy).getWebTest,
        deleteWebTest: createHomeApi(proxy).deleteWebTest,
        updateWebTest: createHomeApi(proxy).updateWebTest,
        getWebTestDelay: createHomeApi(proxy).getWebTestDelay,
        getWebTestIp: createHomeApi(proxy).getWebTestIp,
        closeConnection: createConnApi(proxy).closeConnection,
        closeAllConnection: createConnApi(proxy).closeAllConnection,
        getRules: createRuleApi(proxy).getRules,
        getRuleNum: createRuleApi(proxy).getRuleNum,
        getIgnore: createRuleApi(proxy).getIgnore,
        updateIgnore: createRuleApi(proxy).updateIgnore,
        getTemplateList: createRuleApi(proxy).getTemplateList,
        getTemplateById: createRuleApi(proxy).getTemplateById,
        deleteTemplateById: createRuleApi(proxy).deleteTemplateById,
        updateTemplate: createRuleApi(proxy).updateTemplate,
        createTemplate: createRuleApi(proxy).createTemplate,
        testTemplate: createRuleApi(proxy).testTemplate,
        switchTemplate: createRuleApi(proxy).switchTemplate,
        addProfileFromInput: createProfilesApi(proxy).addProfileFromInput,
        addProfileFromFile: createProfilesApi(proxy).addProfileFromFile,
        deleteProfile: createProfilesApi(proxy).deleteProfile,
        updateProfile: createProfilesApi(proxy).updateProfile,
        getProfileList: createProfilesApi(proxy).getProfileList,
        refreshProfile: createProfilesApi(proxy).refreshProfile,
        switchProfile: createProfilesApi(proxy).switchProfile,
        getDNS: createDnsApi(proxy).getDNS,
        updateDNS: createDnsApi(proxy).updateDNS,
        switchDNS: createDnsApi(proxy).switchDNS,
        getMihomo: createMihomoApi(proxy).getMihomo,
        updateMihomo: createMihomoApi(proxy).updateMihomo,
        waitRunning: createMihomoApi(proxy).waitRunning,
        getAdmin: createMihomoApi(proxy).getAdmin,
        enableProxy: createPandoraApi(proxy).enableProxy,
        disableProxy: createPandoraApi(proxy).disableProxy,
        checkAddressPort: createPandoraApi(proxy).checkAddressPort,
        configDir: createPandoraApi(proxy).configDir,
        exit: createPandoraApi(proxy).exit,
    };
}
