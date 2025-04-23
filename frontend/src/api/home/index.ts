import SparkMD5 from "spark-md5";
import { WebTest } from "@/types/webtest";

// 获取版本
const getVersion = (proxy: any) =>
  async function () {
    const data = await proxy.$http.get("/version");
    return data["version"];
  };

// 获取 Mihomo 基本配置
const getConfigs = (proxy: any) =>
  async function () {
    return await proxy.$http.get("/configs");
  };

// 更新 Mihomo 基本配置
const updateConfigs = (proxy: any) =>
  async function (configs: any) {
    return await proxy.$http.patch("/configs", configs);
  };

// 获取 分组 md5
const getGroupMd5 = (proxy: any) =>
  async function () {
    const data = await proxy.$http.get("/group");
    const proxies = data["proxies"];
    const end = [];
    for (const proxy of proxies) {
      end.push({
        name: proxy["name"],
        now: proxy["now"],
      });
    }
    // 排序
    end.sort((obj1, obj2) => obj1.name.localeCompare(obj2.name));
    // 服务器md5
    const jsonString = JSON.stringify(end);
    return SparkMD5.hash(jsonString);
  };

// 获取 WebTest 列表
const getWebTest = (proxy: any) =>
  async function () {
    return await proxy.$http.get("/webtest");
  };

// 删除 WebTest
const deleteWebTest = (proxy: any) =>
  async function (webTest: WebTest) {
    return await proxy.$http.post("/webtest/delete", webTest);
  };

// 修改 WebTest
const updateWebTest = (proxy: any) =>
  async function (webTest: WebTest) {
    return await proxy.$http.put("/webtest", webTest);
  };

// 获取 delay
const getWebTestDelay = (proxy: any) =>
  async function (list: WebTest[]): Promise<WebTest[]> {
    return await proxy.$http.post("/webtest/delay", list);
  };

//  获取落地 ip
const getWebTestIp = (proxy: any) =>
  async function (data: any) {
    return await proxy.$http.post("/webtest/ip", data);
  };

export default function createHomeApi(proxy: any) {
  return {
    getVersion: getVersion(proxy),
    getConfigs: getConfigs(proxy),
    getGroupMd5: getGroupMd5(proxy),
    updateConfigs: updateConfigs(proxy),
    getWebTest: getWebTest(proxy),
    deleteWebTest: deleteWebTest(proxy),
    updateWebTest: updateWebTest(proxy),
    getWebTestDelay: getWebTestDelay(proxy),
    getWebTestIp: getWebTestIp(proxy),
  };
}
