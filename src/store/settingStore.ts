import {defineStore} from 'pinia';
import {defaultPersist} from "@/types/persist";

export const useSettingStore = defineStore('setting', {
    state: () => ({
        testUrl: 'https://www.google.com/blank.html',
        port: 12345,
        bindAddress: "127.0.0.1",
        stack: 'Mixed',
        ipv6: false,
        dns: false,
        startup: false,
        auth: true
    }),
    actions: {
        setTestUrl(testUrl: any) {
            this.testUrl = testUrl;
        },
        setPort(port: any) {
            this.port = Number(port);
        },
        setStack(stack: any) {
            this.stack = stack;
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
        setBindAddress(bindAddress: any) {
            this.bindAddress = bindAddress;
        },
        setAuth(auth: any) {
            this.auth = auth;
        },
    },
    persist: defaultPersist,
});
