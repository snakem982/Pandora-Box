import {defineConfig} from 'vite'
import vue from '@vitejs/plugin-vue'
import Icons from "unplugin-icons/vite";
import IconsResolver from "unplugin-icons/resolver";
import AutoImport from 'unplugin-auto-import/vite'
import Components from 'unplugin-vue-components/vite'
import {ElementPlusResolver} from 'unplugin-vue-components/resolvers'
import VueI18nPlugin from '@intlify/unplugin-vue-i18n/vite'
import path from 'path'

const pathSrc = path.resolve(__dirname, 'src')

const host = process.env.TAURI_DEV_HOST;

// https://vitejs.dev/config/
export default defineConfig({
    build: {
        emptyOutDir: true,
        outDir: '../backend/static/dist',
    },
    resolve: {
        alias: {
            '@': pathSrc,
        },
    },
    plugins: [vue(),
        AutoImport({
            imports: ["vue"],
            resolvers: [
                ElementPlusResolver(),
                IconsResolver({
                    prefix: "Icon",
                }),
            ],
            dts: path.resolve(pathSrc, "auto-imports.d.ts"),
        }),
        Components({
            resolvers: [
                IconsResolver({
                    prefix: 'icon',
                    enabledCollections: ["ep", "mdi"],
                }),
                ElementPlusResolver()
            ],
            dts: path.resolve(pathSrc, 'components.d.ts'),
        }),
        Icons({
            autoInstall: true,
            compiler: "vue3",
        }),
        VueI18nPlugin({
            include: [path.resolve(pathSrc, './locales/**')],
        }),
    ],
    // Vite options tailored for Tauri development and only applied in `tauri dev` or `tauri build`
    //
    // 1. prevent vite from obscuring rust errors
    clearScreen: false,
    // 2. tauri expects a fixed port, fail if that port is not available
    server: {
        open: false,
        port: 1420,
        strictPort: true,
        host: host || false
    },
})
