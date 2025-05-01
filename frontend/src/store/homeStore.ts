import {defineStore} from 'pinia';
import {defaultPersist} from "@/types/persist";

export const useHomeStore = defineStore('home', {
    state: () => ({
        startTime: 0,
        os: '',
        md5: '',
        md6: '',
        ip: {
            query: '',
            regionName: '',
            country: '',
            city: '',
            isp: '',
            timezone: '',
            as: '',
        },
    }),
    actions: {
        setStartTime(startTime: number) {
            this.startTime = startTime;
        },
        setOS(os: string) {
            this.os = os;
        },
        setMd5(md5: string) {
            this.md5 = md5;
        },
        setMd6(md6: string) {
            this.md6 = md6;
        },
        setIp(ip: any) {
            this.ip = ip;
        },
    },
    persist: defaultPersist
});
