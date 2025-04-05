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
import './style.css'
import './assets/theme/basic.css'
import {useMenuStore} from "@/store/menuStore.js";

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

// 将vue挂载到页面app元素
const app = createApp(App)
app.use(pinia)
app.use(ElementPlus)
app.use(VueApexCharts);
app.use(i18n)
app.use(router);
app.mount('#app')

// 激活menu
const menuStore = useMenuStore()
router.afterEach((to) => {
    const split = to.path.split("/");
    menuStore.setMenu(split[1]);
    if (split.length > 2) {
        menuStore.setRuleMenu(split[2]);
    }
});

