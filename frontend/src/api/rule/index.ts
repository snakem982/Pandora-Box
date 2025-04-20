// 获取规则列表

const getRules = (proxy: any) => async function () {
    const data = await proxy.$http.get('/rules');
    return data['rules'];
}

// 获取 ignore
const getIgnore = (proxy: any) => async function (): Promise<string[]> {
    return await proxy.$http.get('/rule/ignore');
}

// 保存 ignore
const updateIgnore = (proxy: any) => async function (ignores: any) {
    return await proxy.$http.put('/rule/ignore', ignores);
}

export default function createRuleApi(proxy: any) {
    return {
        getRules: getRules(proxy),
        getIgnore: getIgnore(proxy),
        updateIgnore: updateIgnore(proxy),
    }
}