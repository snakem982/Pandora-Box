const isWindows = process.platform === 'win32';

module.exports = {
    appId: 'com.snakem982.pandora-box',
    productName: 'Pandora-Box',
    asar: true,
    electronLanguages: ['en-US', 'zh-CN', 'ru'],
    directories: {
        output: 'out',
    },
    files: [
        'dist/**/*',
        'package.json',
    ],
    extraResources: [
        {
            from: isWindows ? 'src-go/px.exe' : 'src-go/px',
            to: isWindows ? 'px.exe' : 'px',
        },
        {
            from: 'public/tray.png',
            to: 'tray.png',
        },
    ],
    protocols: [
        {
            name: 'Pandora-Box Protocol',
            schemes: ['pandora-box'],
        },
    ],
    electronFuses: {
        resetAdHocDarwinSignature: process.platform === 'darwin'
            && (process.env.ARCH ?? process.arch) === 'arm64',
        runAsNode: false,
        enableCookieEncryption: true,
        enableNodeOptionsEnvironmentVariable: false,
        enableNodeCliInspectArguments: false,
        enableEmbeddedAsarIntegrityValidation: true,
        onlyLoadAppFromAsar: true,
    },
    win: {
        icon: 'build/appicon.ico',
        target: ['nsis', 'zip'],
        artifactName: 'Pandora-Box-v${version}-windows-${arch}.${ext}',
    },
    nsis: {
        artifactName: 'Pandora-Box-v${version}-windows-${arch}-setup.${ext}',
        oneClick: false,
        allowToChangeInstallationDirectory: true,
    },
    mac: {
        icon: 'build/appicon.icns',
        target: ['dmg'],
        artifactName: 'Pandora-Box-v${version}-macos-${arch}.${ext}',
        extendInfo: {
            LSMinimumSystemVersion: '11.0',
        },
    },
    dmg: {
        title: 'Pandora-Box-${arch}',
    },
    linux: {
        icon: 'build/appicon.png',
        category: 'Network',
        maintainer: 'snakem982 <127644738+snakem982@users.noreply.github.com>',
        target: ['deb', 'rpm', 'pacman'],
        artifactName: 'Pandora-Box-v${version}-linux-${arch}.${ext}',
    },
    deb: {
        mimeTypes: ['x-scheme-handler/pandora-box'],
    },
    rpm: {
        vendor: 'snakem982',
        compression: 'gzip',
    },
    pacman: {
        mimeTypes: ['x-scheme-handler/pandora-box'],
    },
};