import {createRouter, createWebHistory, RouteRecordRaw} from 'vue-router'

import About from '../views/About.vue';
import Connections from '../views/Connections.vue';
import Crawl from '../views/Crawl.vue';
import General from '../views/General.vue';
import Log from '../views/Log.vue';
import Profile from '../views/Profile.vue';
import Proxy from '../views/Proxy.vue';
import Rule from '../views/Rule.vue';

//路由记录集合
let routes: RouteRecordRaw[] = [
    {
        path: "/about",
        component: About
    },
    {
        path: "/connection",
        component: Connections
    },
    {
        path: "/crawl",
        component: Crawl
    },
    {
        path: "/general",
        component: General
    },
    {
        path: "/",
        component: General
    },
    {
        path: "/log",
        component: Log
    },
    {
        path: "/profile",
        component: Profile
    },
    {
        path: "/proxy",
        component: Proxy
    },
    {
        path: "/rule",
        component: Rule
    },

];

//创建路由器
let router = createRouter({
    history: createWebHistory(),
    routes
});

//导出
export default router;