// 获取Mihomo
const getMihomo = (proxy: any) => async function () {
    return await proxy.$http.get('/mihomo');
}

// 更新Mihomo
const updateMihomo = (proxy: any) => async function (configs: any) {
    return await proxy.$http.put('/mihomo', configs);
}

// 等待 Mihomo 切换完成
const waitRunning = (proxy: any) => async function () {
    return await proxy.$http.get('/wait');
}

// 获取Mihomo
const getAdmin = (proxy: any) => async function () {
    return await proxy.$http.get('/mihomo/admin');
}


export default function createMihomoApi(proxy: any) {
    return {
        getMihomo: getMihomo(proxy),
        updateMihomo: updateMihomo(proxy),
        waitRunning: waitRunning(proxy),
        getAdmin: getAdmin(proxy),
    }
}