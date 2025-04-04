import {defineStore} from 'pinia';

export const useMenuStore = defineStore('menu', {
    state: () => ({
        menu: '',
        path: '',
        rule: 'rule',
        proxy: false,
        tun: false,
    }),
    actions: {
        setMenu(menu: string) {
            this.menu = menu;
        },
        setPath(path: string) {
            this.path = path;
        },
        setRule(rule: string) {
            this.rule = rule;
        },
        setProxy(proxy: string) {
            this.proxy = proxy;
        },
        setTun(tun: string) {
            this.tun = tun;
        }
    },
    persist: true,
});
