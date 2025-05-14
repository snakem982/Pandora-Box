import type {ForgeConfig} from '@electron-forge/shared-types';
import {MakerSquirrel} from '@electron-forge/maker-squirrel';
import {MakerDMG} from '@electron-forge/maker-dmg';
import {MakerDeb} from '@electron-forge/maker-deb';
import {MakerRpm} from '@electron-forge/maker-rpm';
import {VitePlugin} from '@electron-forge/plugin-vite';
import {FusesPlugin} from '@electron-forge/plugin-fuses';
import {FuseV1Options, FuseVersion} from '@electron/fuses';

const isWindows = process.platform === 'win32';
const extraResource = isWindows ? ['src-go/px.exe'] : ['src-go/px'];

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
        new MakerSquirrel({
            authors: 'snakem982',
            description: 'A Simple Mihomo GUI'
        }),
        new MakerDMG({
            icon: 'build/appicon.icns',
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
