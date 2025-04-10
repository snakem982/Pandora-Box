import {defineStore} from 'pinia';

export const useProxiesStore = defineStore('proxies', {
    state: () => ({
        isHide: false,
        isSort: false,
        active: '',
    }),
    actions: {
        setHide(isHide: boolean) {
            this.isHide = isHide;
        },
        setSort(isSort: boolean) {
            this.isSort = isSort;
        },
        setActive(active: string) {
            this.active = active;
        },
    },
    persist: true,
});
