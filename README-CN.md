# Pandora-Box
一个简易的 Mihomo/Clash.Meta/Clash 桌面客户端

[下载 APP](https://github.com/snakem982/Pandora-Box/releases)

## 功能特点
- 支持 本地 HTTP/HTTPS/SOCKS 代理
- 支持 Vmess, Vless, Shadowsocks, Trojan, Tuic, Hysteria, Hysteria2, Wireguard, Juicity 协议
- 支持 Mihomo 配置文件
- 支持 分享链接, 订阅链接, Base64格式
- 内置将节点和订阅转换为 clash(meta) 配置
- 支持 节点爬取
- 自动添加规则分组
- 【实验阶段】支持统一所有订阅的规则和分组 需要版本v0.2.10+

##  支持的系统平台
- Windows 10/11 AMD64/ARM64
- MacOS 10.13+ AMD64
- MacOS 11.0+ ARM64

##  使用手册
- [基本使用](Manual-CN.md)
- [Juicity 配置详解](Juicity.md)
- [统一所有订阅的规则和分组](UnifiedRuleGrouping.md)
- 自定义配置同 [Mihomo](https://wiki.metacubex.one/config/)

## 友情提示
- 因为作者没有Windows电脑，<br>所以软件在Windows上测试不充分，<br>如果运行不了请更换其他软件。
- 首次启动时若页面空白，请以管理员身份运行。
- 提示需要网络连接，请点击允许。
- 如有疑问请留言。

## 为什么不公开源代码？
- 因为代码写得太烂了

## 关于tun模式
### 怎样开启tun？
- Win 右键以管理员身份运行
- Mac 终端输入命令 sudo /Applications/pandora-box-amd64.app/Contents/MacOS/Pandora-Box 
### 为什么不建议开启？
- 支持tunnel模式，软件需要高级权限
- tunnel模式有时会导致 CPU 使用率爆表
- 关闭tunnel，软件内存使用比较稳定

## 界面预览
### 白色主题-通用
![general.png](img%2F1.png)
### 白色主题-节点
![proxies.png](img%2F2.png)
### 白色主题-配置
![proxies.png](img%2F3.png)
### 白色主题-连接
![proxies.png](img%2F4.png)
### 黑色主题-通用
![general.png](img%2Fdark1.png)
### 黑色主题-节点
![general.png](img%2Fdark2.png)
### 黑色主题-配置
![general.png](img%2Fdark3.png)
### 黑色主题-连接
![general.png](img%2Fdark4.png)
