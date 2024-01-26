# palworld-okg-playground
一起去palworld捕获一直帕鲁吧。

## 命令介绍
```shell 

pal init [环境名称]
pal use  [环境名称]

# gameserver 生命周期管理
pal create server 
pal list server 
pal server describe -n [游戏服名称]
pal server delete -n [游戏服名称]
pal server upgrade -n [游戏服名称]  -r [资源规格类型, 可选择: small / medium / large]

# 用户管理 
pal server [游戏服名称] kick user [用户id]
```