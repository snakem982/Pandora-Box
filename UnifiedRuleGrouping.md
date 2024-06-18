# 统一规则分组 <br> Unified Rule Grouping 
需要版本 v0.2.10+ <br>
Since version v0.2.10+

## 实现原理 （Implementation principle）
### 使用以下代码导入订阅节点 <br> Using the following code to import the subscription node 
```yaml
######### 锚点 start #######
# proxy 相关
pr: &pr { type: select, proxies: [ ♻️ 自动选择 ] }
pr2: &pr2 { type: select, proxies: [ 🚀 节点选择,🎯 全球直连,♻️ 自动选择 ] }

#这里是订阅更新和延迟测试相关的
p: &p { type: file, health-check: { enable: true, url: https://www.gstatic.com/generate_204, interval: 900 } }

use: &use
  type: select
  use:
    - crawling

######### 锚点 end #######

proxy-providers:
  crawling:
    <<: *p
    path: {{PANDORA-BOX}}
```

### 使用以下代码定义分组 <br> Use the following code to define the grouping
```yaml
proxy-groups:
  - { name: 🚀 节点选择, <<: [ *pr,*use ], }
  - { name: ♻️ 自动选择, <<: *use, url: https://www.gstatic.com/generate_204,lazy: true,interval: 3600,tolerance: 50, type: url-test }
  - { name: 🎯 全球直连, proxies: [ DIRECT,🚀 节点选择,♻️ 自动选择 ], type: select }
  - { name: 🛑 全球拦截, proxies: [ REJECT,DIRECT ], type: select }
  - { name: 🐟 漏网之鱼, <<: [ *pr2,*use ], type: select }
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

### 以下代码不能缺少 <br> The following code cannot be missing
```yaml
#这里是订阅更新和延迟测试相关的
p: &p { type: file, health-check: { enable: true, url: https://www.gstatic.com/generate_204, interval: 900 } }

proxy-providers:
  crawling:
    <<: *p
    path: {{PANDORA-BOX}}
```

## 其他可参考 Others
https://wiki.metacubex.one/example/conf/#__tabbed_3_1
