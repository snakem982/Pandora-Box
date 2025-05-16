// 开启代理
const enableProxy = (proxy: any) => async function (configs: any) {
    return await proxy.$http.put('/pandora/enableProxy', configs);
}

// 关闭代理
const disableProxy = (proxy: any) => async function () {
    return await proxy.$http.get('/pandora/disableProxy');
}

// 检测地址端口是否可用
const checkAddressPort = (proxy: any) => async function (configs: any) {
    return await proxy.$http.put('/pandora/checkAddressPort', configs);
}

// 获取配置文件目录
const configDir = (proxy: any) => async function () {
    return await proxy.$http.get('/pandora/configDir');
}

// 退出Px
const exit = (proxy: any) => async function () {
    return await proxy.$http.get('/pandora/exit');
}

export default function createPandoraApi(proxy: any) {
    return {
        enableProxy: enableProxy(proxy),
        disableProxy: disableProxy(proxy),
        checkAddressPort: checkAddressPort(proxy),
        configDir: configDir(proxy),
        exit: exit(proxy),
    }
}