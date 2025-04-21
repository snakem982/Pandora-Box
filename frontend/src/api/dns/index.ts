// 获取dns
const getDNS = (proxy: any) => async function () {
    return await proxy.$http.get('/pDns');
}

// 更新dns
const updateDNS = (proxy: any) => async function (configs: any) {
    return await proxy.$http.put('/pDns', configs);
}

// 切换dns
const switchDNS = (proxy: any) => async function (configs: any) {
    return await proxy.$http.post('/pDns/switch', configs);
}

export default function createDnsApi(proxy: any) {
    return {
        getDNS: getDNS(proxy),
        updateDNS: updateDNS(proxy),
        switchDNS: switchDNS(proxy),
    }
}