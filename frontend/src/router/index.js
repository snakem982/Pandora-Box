import {createRouter, createWebHistory} from 'vue-router';

import Home from '@/views/Home.vue';
import Setting from '@/views/Setting.vue';
import Proxies from '@/views/Proxies.vue';
import Profiles from '@/views/Profiles.vue';
import Rule from "@/views/Rule.vue";
import Now from '@/views/rule/Now.vue';
import Group from '@/views/rule/Group.vue';
import Ignore from '@/views/rule/Ignore.vue';
import Connection from "@/views/Connection.vue";
import Log from "@/views/Log.vue";
import Crawl from "@/views/Crawl.vue";

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
        children: [
            {
                path: 'Now',
                name: 'Now',
                component: Now,
            },
            {
                path: 'Group',
                name: 'Group',
                component: Group,
            },
            {
                path: 'Ignore',
                name: 'Ignore',
                component: Ignore,
            }
        ],
    },
    {
        path: '/Connection',
        name: 'Connection',
        component: Connection,
    },
    {
        path: '/Log',
        name: 'Log',
        component: Log,
    },
    {
        path: '/Crawl',
        name: 'Crawl',
        component: Crawl,
    },
];

const router = createRouter({
    history: createWebHistory(),
    routes,
});

export default router;
