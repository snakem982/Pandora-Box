// 获取Mihomo
const getMihomo = (proxy: any) => async function () {
    return await proxy.$http.get('/mihomo');
}

// 更新Mihomo
const updateMihomo = (proxy: any) => async function (configs: any) {
    return await proxy.$http.put('/mihomo', configs);
}


export default function createMihomoApi(proxy: any) {
    return {
        getMihomo: getMihomo(proxy),
        updateMihomo: updateMihomo(proxy),
    }
}