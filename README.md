# 阿波罗配置中心golang客户端

## 功能

* 多 namespace 支持
* 容错，本地缓存
* 零依赖
* 适配多种配置格式（.yaml  .json）

## 依赖

**go 1.9** 或更新

## 安装

```sh
go get ksogit.kingsoft.net/gz_svr_dev/golib_configcenter.git
```

## 使用

### 使用 app.yaml 配置文件启动

```golang
    cfgCenter := new(configcenter.ConfigCenter)
    appConfigPath := "src/app.yaml"           //配置文件路径
    err := cfgCenter.Init(appConfigPath)
    if err != nil {
    fmt.Println("cfgCenter init error:", err)
    return
  }
```

### 使用自定义配置启动

```golang
   cfgCenter := new(configcenter.ConfigCenter)
   conf:=&agollo.Conf{
    AppID: "app-apollo-demo",
    Cluster: "default",
    NameSpaceNames:[]string{"application","testyaml.yaml","testjson.json"},
    IP: "120.131.9.219:8080",
    EnvLocal: true,
    EnvLocalPath: "catchfile",
  }
   err :=cfgCenter.InitWithConf(conf)
   if err != nil {
    fmt.Println("cfgCenter init error:", err)
    return
  }
```

### 监听配置更新

```golang
  // 默认监听 namespace=application , key=apollo  只提供对第一级的key的监听
  cfgCenter.RegisterKeyWatchFuncDefault("apollo",
    // 监听回调
    func(oldValue, newValue interface{}, changeType string) error {
      fmt.Printf("oldValue:%s,\n newValue:%s, \n changeType:%s\n",
        oldValue, newValue, changeType.String())
      return nil
    })
  
  // 监听 namespace=testyaml.yaml , key=name   只提供对第一级的key的监听
  cfgCenter.RegisterKeyWatchFunc("testyaml.yaml", "name",
    func(oldValue, newValue interface{}, changeType string) error {
      fmt.Printf("oldValue:%s,\n newValue:%s, \n changeType:%s\n",
        oldValue, newValue, changeType.String())
      return nil
    })
  
  // 监听 namespace=testjson.json , key=path   只提供对第一级的key的监听
  cfgCenter.RegisterKeyWatchFunc("testjson.json", "path",
    func(oldValue, newValue interface{}, changeType string) error {
      fmt.Printf("oldValue:%s,\n newValue:%s, \n changeType:%s\n",
        oldValue, newValue, changeType.String())
      return nil
    })
  
```

### 获取配置

```golang
  // 获取默认 namespace=application , key=apollo  默认返回值=yemp
  value, sourceType, err := cfgCenter.GetConfigValue("apollo", "yemp")
  if err != nil {
    fmt.Println(err)
  }
  fmt.Println(sourceType.String())
  if sourceType == configcenter.LOCAL{
    fmt.Println("本地")
  }else if sourceType == configcenter.REMOTE{
    fmt.Println("远程")
  }else if sourceType == configcenter.DEFAULT{
    fmt.Println("默认")
  }
  fmt.Println("value:", value)
  // 获取 namespace=testyaml.yaml , key=children  默认返回值=yemp
  value1, err := cfgCenter.GetConfigValueWithNameSpace("testyaml.yaml",
    "children", "yemp")
  if err != nil{
    fmt.Println(err)
  }
  fmt.Println("children:", value1)
  if val, ok := value1.(map[interface{}]interface{}); ok {
    fmt.Println(val["info"])
  }
```

注：新建项目的默认application如果没有第一次发布，那么就会阻塞客户端对其他namespace的配置的更新监听和查询