import type {ForgeConfig} from '@electron-forge/shared-types';
import {MakerWix} from '@electron-forge/maker-wix';
import {MakerDMG} from '@electron-forge/maker-dmg';
import {MakerDeb} from '@electron-forge/maker-deb';
import {MakerRpm} from '@electron-forge/maker-rpm';
import {VitePlugin} from '@electron-forge/plugin-vite';
import {FusesPlugin} from '@electron-forge/plugin-fuses';
import {FuseV1Options, FuseVersion} from '@electron/fuses';

const isWindows = process.platform === 'win32';
const extraResource = isWindows ? ['src-go/px.exe'] : ['src-go/px'];
const arch = process.env.ARCH || process.arch;

const config: ForgeConfig = {
    packagerConfig: {
        asar: true,
        icon: 'build/appicon',
        extraResource,
        extendInfo: {
            LSMinimumSystemVersion: "10.13.0"
        }
    },
    rebuildConfig: {},
    makers: [
        new MakerWix({
            manufacturer: 'snakem982',
            description: 'A Simple Mihomo GUI',
            icon: 'build/appicon.ico',
            ui: {
                chooseDirectory: true,
                localizations: ["zh-CN", "en-US"],
            }
        }),
        new MakerDMG({
            icon: 'build/appicon.icns',
            title: `Pandora-Box-${arch}`,  // dmg 挂载卷名称
        }),
        new MakerRpm({
            options: {
                icon: 'build/appicon.png',
                homepage: 'https://github.com/snakem982/Pandora-Box',
            }
        }),
        new MakerDeb({
            options: {
                icon: 'build/appicon.png',
                maintainer: 'snakem982',
                homepage: 'https://github.com/snakem982/Pandora-Box',
            }
        })
    ],
    plugins: [
        new VitePlugin({
            // `build` can specify multiple entry builds, which can be Main process, Preload scripts, Worker process, etc.
            // If you are familiar with Vite configuration, it will look really familiar.
            build: [
                {
                    // `entry` is just an alias for `build.lib.entry` in the corresponding file of `config`.
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
        // Fuses are used to enable/disable various Electron functionality
        // at package time, before code signing the application
        new FusesPlugin({
            version: FuseVersion.V1,
            [FuseV1Options.RunAsNode]: false,
            [FuseV1Options.EnableCookieEncryption]: true,
            [FuseV1Options.EnableNodeOptionsEnvironmentVariable]: false,
            [FuseV1Options.EnableNodeCliInspectArguments]: false,
            [FuseV1Options.EnableEmbeddedAsarIntegrityValidation]: true,
            [FuseV1Options.OnlyLoadAppFromAsar]: true,
        }),
    ],
};

export default config;
