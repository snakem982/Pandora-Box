# 统一规则分组 <br> Unified Rule Grouping

## 实现原理 （Implementation principle）
### 使用以下代码导入订阅节点 <br> Using the following code to import the subscription node 
```yaml
# 以下代码不能缺少
# The following code cannot be missing
proxy-providers:
  pandora-box:
    type: file
    path: {{PANDORA-BOX}}
# 以上代码不能缺少
# The above code cannot be missing
```

### 使用以下代码定义分组 <br> Use the following code to define the grouping
```yaml
proxy-groups:
  - name: 🚀 节点选择
    type: select
    proxies:
      - ♻️ 自动选择
    include-all: true

  - name: ♻️ 自动选择
    type: url-test
    url: https://www.google.com/blank.html
    interval: 600
    tolerance: 30
    include-all: true

  - name: 🎯 全球直连
    type: select
    proxies:
      - DIRECT
      - 🚀 节点选择
      - ♻️ 自动选择

  - name: 🛑 全球拦截
    type: select
    proxies:
      - REJECT
      - DIRECT

  - name: 🐟 漏网之鱼
    type: select
    proxies:
      - 🚀 节点选择
      - 🎯 全球直连
      - 🛑 全球拦截
      - ♻️ 自动选择
    include-all: true
```

### 使用以下代码定义规则 <br> Using the following code to define the rule
```yaml
rules:
  - DOMAIN-SUFFIX,googlevideo.com,🚀 节点选择
  - DOMAIN-SUFFIX,youtube.com,🚀 节点选择
  - DOMAIN-SUFFIX,baidujs.cnys.com,🛑 全球拦截
  - DOMAIN-SUFFIX,aliimg.com,🎯 全球直连
  - GEOIP,CN,🎯 全球直连
  - MATCH,🐟 漏网之鱼
```

## 其他可参考 Others
https://wiki.metacubex.one/example/conf/#__tabbed_3_1
