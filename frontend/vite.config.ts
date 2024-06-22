import {defineConfig} from 'vite'
import vue from '@vitejs/plugin-vue'
import viteImagemin from 'vite-plugin-imagemin';

// https://vitejs.dev/config/
export default defineConfig({
    plugins: [vue(),
        viteImagemin({
            gifsicle: {
                optimizationLevel: 7, // 设置GIF图片的优化等级为7
                interlaced: false // 不启用交错扫描
            },
            optipng: {
                optimizationLevel: 7 // 设置PNG图片的优化等级为7
            },
            mozjpeg: {
                quality: 20 // 设置JPEG图片的质量为20
            },
            pngquant: {
                quality: [0.8, 0.9], // 设置PNG图片的质量范围为0.8到0.9之间
                speed: 4 // 设置PNG图片的优化速度为4
            },
            svgo: {
                plugins: [
                    {
                        name: 'removeViewBox' // 启用移除SVG视图框的插件
                    },
                    {
                        name: 'removeEmptyAttrs',
                        active: false // 不启用移除空属性的插件
                    }
                ]
            }
        })
    ],
    build: {
        chunkSizeWarningLimit: 1000, // 单个模块文件大小限制1000KB
        terserOptions: {
            // 清除代码中console和debugger
            compress: {
                drop_console: true,
                drop_debugger: true,
            },
        },
        rollupOptions: {
            output: {
                manualChunks(id) {
                    // 静态资源拆分
                    // @ts-ignore
                    if (id.includes("node_modules")) {
                        return id
                            .toString()
                            .split("node_modules/")[1]
                            .split("/")[0]
                            .toString();
                    }
                },
                // 设置chunk的文件名格式
                chunkFileNames: (chunkInfo) => {
                    const facadeModuleId = chunkInfo.facadeModuleId
                        ? chunkInfo.facadeModuleId.split("/")
                        : [];
                    const fileName1 =
                        facadeModuleId[facadeModuleId.length - 2] || "[name]";
                    // 根据chunk的facadeModuleId（入口模块的相对路径）生成chunk的文件名
                    return `js/${fileName1}/[name].[hash].js`;
                },
                // 设置入口文件的文件名格式
                entryFileNames: "js/[name].[hash].js",
                // 设置静态资源文件的文件名格式
                assetFileNames: "[ext]/[name].[hash].[ext]",
            },
        },
    },
});
