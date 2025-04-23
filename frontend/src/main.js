import {createApp} from 'vue'
import App from './App.vue'
import router from './router';
import {createPinia} from 'pinia'
import piniaPluginPersistedstate from 'pinia-plugin-persistedstate'
import {createI18n} from 'vue-i18n'
import messages from '@intlify/unplugin-vue-i18n/messages'
import ElementPlus from 'element-plus'
import VueApexCharts from "vue3-apexcharts";
import 'element-plus/dist/index.css'
import './styles/global.css'
import './styles/basic.css'
import {useMenuStore} from "@/store/menuStore";
import {useWebStore} from "@/store/webStore";
import {AxiosRequest} from "./util/axiosRequest";
import {useHomeStore} from "@/store/homeStore.js";


// 国际化设置
const i18n = createI18n({
    locale: 'zh',
    messages,
    legacy: false,
    globalInjection: true,
})

// 全局状态管理
const pinia = createPinia()
pinia.use(piniaPluginPersistedstate)

// 加载所需组件
const app = createApp(App)
app.use(pinia)
app.use(ElementPlus)
app.use(VueApexCharts);
app.use(i18n)
app.use(router);

// 获取api地址、端口、密钥
const url = window.location.search;
const params = new URLSearchParams(url);
const webStore = useWebStore();
const host = params.get('host');
const port = params.get('port');
const secret = params.get('secret');
if (host) {
    webStore.setHost(host);
}
if (port) {
    webStore.setPort(port);
}
if (secret) {
    webStore.setSecret(secret);
}

// 注册 Axios 实例到全局
app.config.globalProperties.$http = new AxiosRequest(webStore.baseUrl, webStore.secret);

// 激活menu
const menuStore = useMenuStore()
router.afterEach((to) => {
    const split = to.path.split("/");
    menuStore.setMenu(split[1]);
    if (split.length > 2 && split[1] === "Rule") {
        menuStore.setRuleMenu(split[2]);
    }
});

// 设置起始时间 和 操作系统类型
const homeStore = useHomeStore();

function getOSFromUserAgent() {
    const userAgent = navigator.userAgent.toLowerCase();

    if (userAgent.includes("win")) {
        return "Windows";
    } else if (userAgent.includes("mac")) {
        return "MacOS";
    } else if (userAgent.includes("linux")) {
        return "Linux";
    } else {
        return "Unknown OS";
    }
}

homeStore.setOS(getOSFromUserAgent())
homeStore.setStartTime(Date.now())

// 将vue挂载到页面app元素
app.mount('#app')

