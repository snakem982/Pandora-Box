import {defineStore} from 'pinia';

export const useSettingStore = defineStore('setting', {
    state: () => ({
        testUrl: 'https://www.google.com/blank.html',
    }),
    actions: {
        setTestUrl(testUrl: any) {
            this.testUrl = testUrl;
        },
    },
    persist: true,
});
