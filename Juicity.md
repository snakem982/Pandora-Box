# yaml配置
```yaml
proxies:
    - name: "juicity-client"
      server: 192.168.1.1
      port: 12345
      type: juicity
      uuid: 6005528d-0cb4-406f-affa-5c4d5c69b29e
      password: gogogo
      sni: gogogo.com
      skip-cert-verify: false
      congestion-controller: bbr
      pinned-certchain-sha256: ""
      udp: true
```
## server
必须 
## port
必须
## uuid
必须，用于Juicity的用户唯一识别码
## password
必须，用于Juicity的用户密码，
## sni
服务器名称指示，如果为空，则为server中的地址
## skip-cert-verify
跳过证书验证
## congestion-controller
设置拥塞控制算法，可选项为 bbr
##  pinned-certchain-sha256
自签证书hash

# 分享链接格式 <br> Link Format
```html
juicity://uuid:password@122.12.31.66:port?congestion_control=bbr&sni=www.example.com&allow_insecure=0&pinned_certchain_sha256=CERT_HASH
```