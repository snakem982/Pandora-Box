// scripts/sync-version.js
const fs = require('fs');
const path = require('path');

const version = process.env.VERSION || '';
if (!version.startsWith('v')) {
    console.error('VERSION 环境变量必须是类似 v1.0.0 的格式');
    process.exit(1);
}

const cleanVersion = version.replace(/^v/, '');
const pkgPath = path.resolve(__dirname, '../package.json');

const pkg = JSON.parse(fs.readFileSync(pkgPath, 'utf-8'));
pkg.version = cleanVersion;
fs.writeFileSync(pkgPath, JSON.stringify(pkg, null, 2) + '\n');

console.log(`✔ 已将 package.json 的 version 设置为 ${cleanVersion}`);
