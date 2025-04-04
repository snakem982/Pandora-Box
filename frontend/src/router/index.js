import { createRouter, createWebHistory } from 'vue-router';

import Home from '@/views/Home.vue';
import Setting from '@/views/Setting.vue';
import Proxies from '@/views/Proxies.vue';
import Profiles from '@/views/Profiles.vue';
import Rule from "@/views/Rule.vue";

const routes = [
    {
        path: '/',
        name: 'Start',
        component: Rule,
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
    {
        path: '/Rule',
        name: 'Rule',
        component: Rule,
    },
];

const router = createRouter({
    history: createWebHistory(),
    routes,
});

export default router;
