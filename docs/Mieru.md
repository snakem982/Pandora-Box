# yaml配置
v0.2.23 后支持<br>
Supported after v0.2.23
```yaml
proxies:
  - name: mieru
    type: mieru
    server: 1.2.3.4
    port: 2999
    # port-range: 2090-2099 #（不可同时填写 port 和 port-range）
    transport: TCP # Available: "TCP".
    username: user
    password: password
```
## server
必须 
## port
必须
## username
必须，用于Mieru的用户唯一识别码
## password
必须，用于Mieru的用户密码
## transport
必须，填写 TCP
## port-range
跳跃端口，不可同时填写 port 和 port-range
