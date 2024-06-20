<h1>Pandora-Box</h1>
<p>A Simple Mihomo/Clash.Meta/Clash GUI.</p>

## Build

1、Build Environment

- Node.js [link](https://nodejs.org/en)
- Go [link](https://go.dev/)
- Wails [link](https://wails.io/) ：`go install github.com/wailsapp/wails/v2/cmd/wails@latest`

2、Pull and Build

```bash
git clone -b v2 --single-branch https://github.com/snakem982/Pandora-Box.git

cd Pandora-Box/frontend

npm install

npm build

cd ..

wails build -tags with_gvisor
```