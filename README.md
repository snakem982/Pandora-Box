<div align="center">

  <img src="build/appicon.png" width="160px" alt="Pandora-Box"/>

  <h1>Pandora-Box</h1>

  <p>🌈 A simple desktop client for <strong>Mihomo</strong></p>
  <p>✨ 一个简易的 <strong>Mihomo</strong> 桌面客户端</p>
  <p>✨ Простой настольный клиент для <strong>Mihomo</strong></p>

  <p>
    🇨🇳 <a href="doc/README.zh-CN.md">简体中文</a> | 🇺🇸 <a href="doc/README.en.md">English</a> | 🇷🇺 <a href="doc/README.ru.md">Русский</a>
  </p>

</div>

---

## 📦 Project Overview | 项目简介 | Обзор проекта

**Pandora-Box** is a lightweight and user-friendly cross-platform client for [Mihomo](https://github.com/MetaCubeX/mihomo), supporting multiple proxy protocols, automatic rule grouping, and TUN mode.  
It is designed for both casual and advanced users to easily manage and convert proxy subscriptions.

**Pandora-Box** 是一个跨平台的轻量桌面客户端，适配 [Mihomo](https://github.com/MetaCubeX/mihomo) 内核，支持多种代理协议、规则自动分组与 TUN 模式。界面简洁，功能强大，适合轻量与进阶用户使用。

**Pandora-Box** — это легкий и удобный кроссплатформенный клиент для [Mihomo](https://github.com/MetaCubeX/mihomo), поддерживающий различные прокси-протоколы, автоматическую группировку правил и режим TUN.  
Он разработан как для обычных пользователей, так и для продвинутых, чтобы облегчить управление и конвертацию подписок прокси.

---

## 🌐 Language ｜ 语言选择 ｜ Выбор языка

- 🇨🇳 [查看中文文档](doc/README.zh-CN.md)
- 🇺🇸 [View English Documentation](doc/README.en.md)
- 🇷🇺 [Просмотр русской документации](doc/README.ru.md)

---

## 📥 Get Started ｜ 快速开始 ｜ Начало работы

👉 [Download the Latest Release / 下载最新版本 / Скачать последнюю версию](https://github.com/snakem982/Pandora-Box/releases)

---

## 🛠 Development ｜ 开发 ｜ Разработка

If you want to contribute or build Pandora-Box locally, refer to the resources below:  
如果你想参与开发或构建 Pandora-Box，可以参考以下资源：  
Если вы хотите принять участие в разработке или собрать Pandora-Box локально, воспользуйтесь следующими ресурсами:

### 🔧 Prerequisites | 前置依赖 | Предварительные требования

- [Node.js](https://nodejs.org/) ≥ 18 (for building UI components or tooling)
- [Go](https://go.dev/) ≥ 1.24 (for integration with Mihomo or backend modules)

### 🧪 Build Instructions | 构建指南 | Инструкции по сборке

```bash
# Install dependencies
npm install
cd src-go
go mod tidy

# Build px backend
CGO_ENABLED=0 go build -tags=with_gvisor -trimpath -ldflags "-X github.com/snakem982/pandora-box/api.Version=v-test" -o px(.exe)
cd ..

# Build desktop app
npm run package

# Run in dev mode
npm run start
```

---

## 🧭 More Information ｜ 更多信息 ｜ Дополнительная информация

- ❤️ *Special thanks to [ZMTO](https://console.zmto.com/?affid=1579&oid=8) for providing the VPS hosting support for this project.*
- ✅ [Project Issues](https://github.com/snakem982/Pandora-Box/issues)
- 📄 [License (GPL-3.0)](./LICENSE)

---

📝 This README was generated with the assistance of AI and reviewed by the developer.  
📝 本文档内容由 AI 辅助生成，并由开发者校对。  
📝 Этот README создан при поддержке ИИ и проверен разработчиком.
