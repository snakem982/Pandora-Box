import type { ForgeConfig } from '@electron-forge/shared-types';
import { MakerWix } from '@electron-forge/maker-wix';
import { MakerZIP } from '@electron-forge/maker-zip';
import { MakerDMG } from '@electron-forge/maker-dmg';
import { MakerDeb } from '@electron-forge/maker-deb';
import { MakerRpm } from '@electron-forge/maker-rpm';
import { VitePlugin } from '@electron-forge/plugin-vite';
import { FusesPlugin } from '@electron-forge/plugin-fuses';
import { FuseV1Options, FuseVersion } from '@electron/fuses';

const isWindows = process.platform === 'win32';
const extraResource = isWindows
    ? ['src-go/px.exe', 'public/tray.png']
    : ['src-go/px', 'public/tray.png'];

const arch = process.env.ARCH || process.arch;
const isDev = process.env.NODE_ENV === 'development';

const config: ForgeConfig = {
    packagerConfig: {
        asar: true,
        icon: 'build/appicon',
        extraResource,
        extendInfo: {
            LSMinimumSystemVersion: '11.0',
        },
        appBundleId: 'com.snakem982.pandora-box',
        protocols: [
            {
                name: 'Pandora-Box Protocol',
                schemes: ['pandora-box'],
            },
        ],
    },
    rebuildConfig: {},
    makers: [
        ...(isWindows && arch !== 'arm64'
            ? [
                new MakerWix({
                    manufacturer: 'snakem982',
                    description: 'A Simple Mihomo GUI',
                    icon: 'build/appicon.ico',
                    ui: { chooseDirectory: true },
                }),
            ]
            : []),
        ...(isWindows && arch === 'arm64' ? [new MakerZIP()] : []),
        new MakerDMG({
            icon: 'build/appicon.icns',
            title: `Pandora-Box-${arch}`,
        }),
        new MakerRpm({
            options: {
                icon: 'build/appicon.png',
                homepage: 'https://github.com/snakem982/Pandora-Box',
            },
        }),
        new MakerDeb({
            options: {
                icon: 'build/appicon.png',
                maintainer: 'snakem982',
                homepage: 'https://github.com/snakem982/Pandora-Box',
                mimeType: ['x-scheme-handler/pandora-box'],
            },
        }),
    ],
    plugins: [
        new VitePlugin({
            build: [
                {
                    entry: 'src-electron/main.ts',
                    config: 'vite.main.config.ts',
                    target: 'main',
                },
                {
                    entry: 'src-electron/preload.ts',
                    config: 'vite.preload.config.ts',
                    target: 'preload',
                },
            ],
            renderer: [
                {
                    name: 'px_window',
                    config: 'vite.config.ts',
                },
            ],
        }),
        // ⚡️ 只在生产环境启用 FusesPlugin，避免和 Vite 冲突
        ...(isDev
            ? []
            : [
                new FusesPlugin({
                    version: FuseVersion.V1,
                    [FuseV1Options.RunAsNode]: false,
                    [FuseV1Options.EnableCookieEncryption]: true,
                    [FuseV1Options.EnableNodeOptionsEnvironmentVariable]: false,
                    [FuseV1Options.EnableNodeCliInspectArguments]: false,
                    [FuseV1Options.EnableEmbeddedAsarIntegrityValidation]: true,
                    [FuseV1Options.OnlyLoadAppFromAsar]: true,
                }),
            ]),
    ],
};

export default config;
