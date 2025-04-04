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


const i18n = createI18n({
    locale: 'zh',
    messages,
    fallbackWarn: false,
    missingWarn: false,
})

const pinia = createPinia()
pinia.use(piniaPluginPersistedstate)

const app = createApp(App)
app.use(pinia)
app.use(ElementPlus)
app.use(VueApexCharts);
app.use(i18n)
app.use(router);
app.mount('#app')
