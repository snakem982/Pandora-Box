<div align="center">
<img src="build/540x540.png"  style="width:160px" />
<h1>Pandora-Box</h1>
<p>A Simple Mihomo GUI.</p>
</div>

## Introduction

- [English](docs/README.md)
- [中文](docs/README-CN.md)


## Build

1、Build Environment

- Node.js [link](https://nodejs.org/en)
- Go [link](https://go.dev/)
- Wails [link](https://wails.io/) ：`go install github.com/wailsapp/wails/v2/cmd/wails@latest`

2、Pull and Build

```bash
git clone https://github.com/snakem982/Pandora-Box.git

cd Pandora-Box/frontend

npm install

npm run build

cd ..

wails build -tags with_gvisor
```

## Commit Submission Specification
```yaml
feat: 新功能（feature）
fix: 修补bug
docs: 文档（documentation）
style: 格式（不影响代码运行的变动）
refactor: 重构（即不是新增功能，也不是修改bug的代码变动）
chore: 构建过程或辅助工具的变动
revert: 撤销，版本回退
perf: 性能优化
test: 测试
improvement: 改进
build: 打包
ci: 持续集成
```