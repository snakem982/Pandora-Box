import {defineConfig} from 'vite';

export default defineConfig({
    build: {
        ssr: 'src-electron/preload.ts',
        outDir: 'dist/electron',
        emptyOutDir: false,
        rollupOptions: {
            external: (id) => !id.startsWith('.') && !id.startsWith('/') && !/^[A-Za-z]:[\\/]/.test(id),
            output: {
                format: 'cjs',
                entryFileNames: 'preload.cjs',
                codeSplitting: false,
            },
        },
    },
});
