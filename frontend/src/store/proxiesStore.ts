import {defineStore} from 'pinia';
import {defaultPersist} from "@/types/persist";

export const useProxiesStore = defineStore('proxies', {
    state: () => ({
        isHide: false,
        isSort: false,
        isVertical: false,
        active: '',
    }),
    actions: {
        setHide(isHide: boolean) {
            this.isHide = isHide;
        },
        setSort(isSort: boolean) {
            this.isSort = isSort;
        },
        setVertical(isVertical: boolean) {
            this.isVertical = isVertical;
        },
        setActive(active: string) {
            this.active = active;
        },
    },
    persist: defaultPersist,
});
