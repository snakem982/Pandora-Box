import {defineStore} from 'pinia';

export const useMenuStore = defineStore('menu', {
    state: () => ({
        menu: 'Home',
        path: '/Home',
        rule: '规则',
        proxy: false,
        tun: false,
        language: 'zh',
        ruleMenu: 'Now',
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
        },
        setLanguage(language: string) {
            this.language = language;
        },
        setRuleMenu(ruleMenu: string) {
            this.ruleMenu = ruleMenu;
        }
    },
    persist: true,
});
