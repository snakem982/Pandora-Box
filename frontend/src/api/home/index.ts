import SparkMD5 from "spark-md5";

// 获取版本
const getVersion = (proxy: any) => async function () {
    const data = await proxy.$http.get('/version');
    return data['version'];
}

// 获取 Mihomo 基本配置
const getConfigs = (proxy: any) => async function () {
    return await proxy.$http.get('/configs');
}

// 更新 Mihomo 基本配置
const updateConfigs = (proxy: any) => async function (configs: any) {
    return await proxy.$http.patch('/configs', configs);
}

// 获取 分组 md5
const getGroupMd5 = (proxy: any) => async function () {
    const data = await proxy.$http.get('/group');
    const proxies = data['proxies']
    const end = []
    for (const proxy of proxies) {
        end.push({
            name: proxy['name'],
            now: proxy['now'],
        });
    }
    // 排序
    end.sort((obj1, obj2) => obj1.name.localeCompare(obj2.name));
    // 服务器md5
    const jsonString = JSON.stringify(end);
    return SparkMD5.hash(jsonString);
}


export default function createHomeApi(proxy: any) {
    return {
        getVersion: getVersion(proxy),
        getConfigs: getConfigs(proxy),
        getGroupMd5: getGroupMd5(proxy),
        updateConfigs: updateConfigs(proxy),
    }
}