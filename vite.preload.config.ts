// @ts-ignore
import {defineConfig, type Plugin} from 'vite';

const preloadBundleConfig: Plugin = {
    name: 'preload-bundle-config',
    config(config) {
        const output = config.build?.rollupOptions?.output;
        if (!output) {
            throw new Error('Expected the Electron Forge preload output configuration');
        }

        for (const options of Array.isArray(output) ? output : [output]) {
            delete options.inlineDynamicImports;
            options.codeSplitting = false;
        }
    },
};

// https://vitejs.dev/config
export default defineConfig({
    plugins: [preloadBundleConfig],
});
