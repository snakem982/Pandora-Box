import {defineStore} from 'pinia';

export const useSettingStore = defineStore('setting', {
    state: () => ({
        testUrl: 'https://www.google.com/blank.html',
        port: 12345,
        stack: 'Mixed',
        lan: false,
        ipv6: false,
        dns: false,
        startup: false,
    }),
    actions: {
        setTestUrl(testUrl: any) {
            this.testUrl = testUrl;
        },
        setPort(port: any) {
            this.port = port;
        },
        setStack(stack: any) {
            this.stack = stack;
        },
        setLan(lan: any) {
            this.lan = lan; 
        },
        setIpv6(ipv6: any) {
            this.ipv6 = ipv6;
        },
        setDns(dns: any) {
            this.dns = dns;
        },
        setStartup(startup: any) {
            this.startup = startup;
        },
    },
    persist: true,
});
