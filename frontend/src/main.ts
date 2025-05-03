import {createApp} from "vue";
import App from "./App.vue";
import router from "@/router";
import {createPinia} from "pinia";
import piniaPluginPersistence from "pinia-plugin-persistedstate";
import {createI18n} from "vue-i18n";
import messages from "@intlify/unplugin-vue-i18n/messages";
import ElementPlus from "element-plus";
import VueApexCharts from "vue3-apexcharts";
import "element-plus/dist/index.css";
import "./styles/global.css";
import "./styles/basic.css";
import {useMenuStore} from "@/store/menuStore";
import {useWebStore} from "@/store/webStore";
import {AxiosRequest} from "@/util/axiosRequest";
import {useHomeStore} from "@/store/homeStore";
import {memoryCache} from "@/types/persist";

const app = createApp(App);

async function bootstrap() {
    // 加载缓存数据
    const keys = ['menu', 'home', 'proxies', 'setting', 'web'];
    for (const key of keys) {
        const val = await window["pxGetItem"](key);
        if (val) {
            memoryCache[key] = val;
        }
    }

    // 国际化设置
    const i18n = createI18n({
        locale: "zh",
        messages,
        globalInjection: true,
    });

    // 全局状态管理
    const pinia = createPinia();
    pinia.use(piniaPluginPersistence);


    // 加载所需组件
    app.use(pinia);
    app.use(ElementPlus);
    app.use(VueApexCharts);
    app.use(i18n);
    app.use(router);

    // 获取api地址、端口、密钥
    const url = window.location.search;
    const params = new URLSearchParams(url);
    const webStore = useWebStore();
    const host = params.get("host");
    const port = params.get("port");
    const secret = params.get("secret");
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
    app.config.globalProperties.$http = new AxiosRequest(
        webStore.baseUrl,
        webStore.secret
    );

    // 激活menu
    const menuStore = useMenuStore();
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

    homeStore.setOS(getOSFromUserAgent());
    homeStore.setStartTime(Date.now());

    // 绑定拖拽指令
    app.directive('pxDrag', {
        mounted(el) {
            let isDragging = false;

            const startDrag = (event: any) => {
                // 如果拖拽元素内部按下且鼠标左键按下
                if (event.target === el && event.buttons === 1) {
                    isDragging = true;
                    document.body.style.cursor = 'move';
                    if (window['pxDrag']) {
                        window['pxDrag']();
                    }
                }
            };

            const endDrag = () => {
                if (isDragging) {
                    isDragging = false;
                    document.body.style.cursor = '';
                }
            };

            // 绑定 mousedown 在元素上
            el.addEventListener('mousedown', startDrag);
            // 绑定 mouseup 在全局，这样能确保捕捉到鼠标在任意位置释放的事件
            document.addEventListener('mouseup', endDrag);
            // 如果需要也可以在 document 上绑定 mousemove 来处理拖拽逻辑

            // 存储 cleanup 函数，方便在组件卸载时移除绑定监听器
            el._cleanup = () => {
                el.removeEventListener('mousedown', startDrag);
                document.removeEventListener('mouseup', endDrag);
            };
        },
        unmounted(el) {
            if (el._cleanup) {
                el._cleanup();
            }
        }
    });


}

// 🚀 启动应用
bootstrap().then(() => app.mount("#app"));



