import { createRouter, createWebHistory } from 'vue-router';

import Home from '@/views/Home.vue';
import Setting from '@/views/Setting.vue';
import Proxies from '@/views/Proxies.vue';
import Profiles from '@/views/Profiles.vue';
import Crawl from "@/views/Crawl.vue";

const routes = [
    {
        path: '/',
        name: 'Start',
        component: Crawl,
    },
    {
        path: '/Home',
        name: 'Home',
        component: Home,
    },
    {
        path: '/Setting',
        name: 'Setting',
        component: Setting,
    },
    {
        path: '/Proxies',
        name: 'Proxies',
        component: Proxies,
    },
    {
        path: '/Profiles',
        name: 'Profiles',
        component: Profiles,
    },
];

const router = createRouter({
    history: createWebHistory(),
    routes,
});

export default router;
