# palworld-okg-playground
一起去palworld捕获一直帕鲁吧。

## 编译打包

```bash
make build-binary
```
编译后的二进制文件 `pal` 在 ./bin 目录下

## 使用说明

### 当前支持的Cloud Provider

AlibabaCloud
- 【PRE】环境变量配置
  ```bash
  export ACCESS_KEY_ID={填写阿里云账号ACCESS KEY ID}
  export ACCESS_KEY_SECRET={填写阿里云账号ACCESS KEY SECRET}
  export REGION_ID={填写集群Region ID, 如cn-hangzhou}
  ```

### 集群管理

```bash
# 列出所有帕鲁游戏服集群信息
pal cluster list

# 创建一个新的集群
pal cluster create --cluster_name [即将创建的集群名称]

# 输入集群ID，切换默认集群
pal cluster choose [集群id]
```

### 游戏服管理

```bash
# 列出默认集群下的所有游戏服
pal server list

# 在默认集群下创建一个新的游戏服
pal server create

# 删除默认集群下指定的游戏服
pal server delete -n [游戏服名称]

# 更新默认集群下指定的游戏服
pal server upgrade -n [游戏服名称]  -r [资源规格类型, 可选择: small / medium / large] -p [规定最大的玩家数量]
```