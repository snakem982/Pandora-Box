<div align="center">
  <img src="../build/appicon.png" style="width:160px" alt="Pandora-Box"/>
  <h1>Pandora-Box</h1>
  <p>A simple desktop client for Mihomo</p>
</div>

## Download

[Download the App](https://github.com/snakem982/Pandora-Box/releases)

## Features

- Supports local HTTP/HTTPS/SOCKS proxies
- Supports Vmess, Vless, Shadowsocks, Trojan, Tuic, Hysteria, Hysteria2, Wireguard, and Mieru protocols
- Supports parsing of share links, subscription links, Base64 format, and YAML format
- Built-in subscription converter to convert various subscription types into Mihomo-compatible configurations
- Automatically adds minimal rule groups to unruly subscriptions
- DNS overwrite option to prevent DNS leaks
- Unified rules and group settings for all subscriptions
- Supports TUN mode

## Supported Platforms

- Windows 10/11 (AMD64 / ARM64)
- macOS 11.0+ (AMD64 / ARM64)
- Linux (AMD64 / ARM64)

## How to Enable TUN Mode

- Go to `Settings` → `Enable Authorization` → Restart the app → When the authorization prompt appears, grant
  permission → TUN mode can then be enabled in the app

## How to Enable LAN Access

Simply change the Listen Address in Settings to 0.0.0.0

## Deeplink Profile Import

Pandora-Box supports importing profiles via deeplink URLs, allowing users to easily add subscriptions from external sources.

### URL Scheme

The deeplink uses the custom protocol `pandora-box://` with the following format:

```
pandora-box://install-config?url=SUBSCRIPTION_URL
```

## Note: Px Requires Network Access

- When prompted, click "Allow" to grant network access

## Common macOS Issues

- See [mac.md](mac/mac.md)


## Preview

| Tab      | New Interface with Different Themes |
|----------|-------------------------------------|
| Home     | ![General](img/home.png)            |
| Settings | ![Setting](img/setting.png)         |
| Proxies  | ![Proxies](img/proxies.png)         |
| Profiles | ![Profiles](img/profiles.png)       |
