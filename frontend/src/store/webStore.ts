import {defineStore} from 'pinia'

export const useWebStore = defineStore('useWebStore', {
    state: () => {
        return {
            port: 0,
            secret: '',
        }
    },
    getters: {
        baseUrl: (state) => "http://127.0.0.1:" + state.port,
    },
    actions: {
        getPortAndSecret() {
            this.port = 9898
            this.secret = '123456'
        },
    },
    persist: true,
})
