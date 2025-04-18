import {Profile} from "@/types/profile";

// 添加配置从input
const addProfileFromInput = (proxy: any) => async function (profile: Profile): Promise<Profile[]> {
    return await proxy.$http.post('/profile', profile);
}

// 添加配置从文件
const addProfileFromFile = (proxy: any) => async function (profile: Profile) {
    return await proxy.$http.post('/profile/file', profile);
}

// 删除配置
const deleteProfile = (proxy: any) => async function (profile: Profile) {
    return await proxy.$http.post('/profile/delete', profile);
}

// 修改配置
const updateProfile = (proxy: any) => async function (profile: Profile) {
    return await proxy.$http.put('/profile', profile);
}

// 获取配置列表
const getProfileList = (proxy: any) => async function (): Promise<Profile[]> {
    return await proxy.$http.get('/profile');
}

// 刷新配置
const refreshProfile = (proxy: any) => async function (profile: Profile) {
    return await proxy.$http.put('/profile/refresh', profile);
}

// 切换配置
const switchProfile = (proxy: any) => async function (profile: Profile) {
    return await proxy.$http.patch('/profile', profile);
}


export default function createProfilesApi(proxy: any) {
    return {
        addProfileFromInput: addProfileFromInput(proxy),
        addProfileFromFile: addProfileFromFile(proxy),
        deleteProfile: deleteProfile(proxy),
        updateProfile: updateProfile(proxy),
        getProfileList: getProfileList(proxy),
        refreshProfile: refreshProfile(proxy),
        switchProfile: switchProfile(proxy),
    }
}