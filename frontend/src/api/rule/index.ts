// 获取规则列表
const getRules = (proxy: any) => async function () {
    const data = await proxy.$http.get('/rules');

    return data['rules'];
}


export default function createRuleApi(proxy: any) {
    return {
        getRules: getRules(proxy),
    }
}