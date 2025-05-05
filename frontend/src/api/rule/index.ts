// 获取规则列表
const getRules = (proxy: any) =>
    async function () {
        const data = await proxy.$http.get("/rules");
        return data["rules"];
    };

// 获取规则数
const getRuleNum = (proxy: any) =>
    async function () {
        const data = await proxy.$http.get("/rule/num");
        return data["data"];
    };

// 获取 ignore
const getIgnore = (proxy: any) =>
    async function (): Promise<string[]> {
        return await proxy.$http.get("/rule/ignore");
    };

// 保存 ignore
const updateIgnore = (proxy: any) =>
    async function (ignores: any) {
        return await proxy.$http.put("/rule/ignore", ignores);
    };

// 获取 template list
const getTemplateList = (proxy: any) =>
    async function (): Promise<any[]> {
        return await proxy.$http.get("/rule/list");
    };

// 获取 template by id
const getTemplateById = (proxy: any) =>
    async function (id: string): Promise<any> {
        return await proxy.$http.get(`/rule/template/${id}`);
    };

// 删除 template by id
const deleteTemplateById = (proxy: any) =>
    async function (id: string) {
        return await proxy.$http.delete(`/rule/template/${id}`);
    };

// 更新 template
const updateTemplate = (proxy: any) =>
    async function (data: any) {
        return await proxy.$http.put(`/rule/template`, data);
    };

// 创建 template
const createTemplate = (proxy: any) =>
    async function (data: any) {
        return await proxy.$http.post(`/rule/template`, data);
    };

// 测试 template
const testTemplate = (proxy: any) =>
    async function (data: any) {
        return await proxy.$http.post(`/rule/test`, data);
    };

// 开启 template
const switchTemplate = (proxy: any) =>
    async function (data: any) {
        return await proxy.$http.post(`/rule/switch`, data);
    };

export default function createRuleApi(proxy: any) {
    return {
        getRules: getRules(proxy),
        getRuleNum: getRuleNum(proxy),
        getIgnore: getIgnore(proxy),
        updateIgnore: updateIgnore(proxy),
        getTemplateList: getTemplateList(proxy),
        getTemplateById: getTemplateById(proxy),
        deleteTemplateById: deleteTemplateById(proxy),
        updateTemplate: updateTemplate(proxy),
        createTemplate: createTemplate(proxy),
        testTemplate: testTemplate(proxy),
        switchTemplate: switchTemplate(proxy),
    };
}
