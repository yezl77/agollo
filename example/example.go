package main

import (
  "fmt"
  "apollo_go/configcenter"
  "time"
)

func main() {
  cfgCenter := new(configcenter.ConfigCenter)
  appConfigPath := "app.yaml"
  err := cfgCenter.Init(appConfigPath)
  if err != nil {
    fmt.Println("cfgCenter init error:", err)
    return
  }

  cfgCenter.RegisterKeyWatchFunc("testyaml.yaml", "name",
    func(oldValue, newValue interface{}, changeType string) error {
      fmt.Printf("oldValue:%s,\n newValue:%s, \n changeType:%s\n",
        oldValue, newValue, changeType)
      return nil
    })
  cfgCenter.RegisterKeyWatchFuncDefault("apollo",
    // 监听回调
    func(oldValue, newValue interface{}, changeType string) error {
      fmt.Printf("oldValue:%s,\n newValue:%s, \n changeType:%s\n",
        oldValue, newValue, changeType)
      return nil
    })

  cfgCenter.RegisterKeyWatchFunc("testjson.json", "path",
    func(oldValue, newValue interface{}, changeType string) error {
      fmt.Printf("oldValue:%s,\n newValue:%s, \n changeType:%s\n",
        oldValue, newValue, changeType)
      return nil
    })

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

  cfgCenter.RegisterKeyWatchFunc("testyaml.yaml", "spouse",
    func(oldValue, newValue interface{}, changeType string) error {
      fmt.Printf("oldValue:%s,\n newValue:%s, \n changeType:%s\n",
        oldValue, newValue, changeType)
      return nil
    })

  value1, _, _ := cfgCenter.GetConfigValueWithNameSpace("testyaml.yaml",
    "children", "yemp3")
  fmt.Println("children:", value1)
  if val, ok := value1.(map[interface{}]interface{}); ok {
    fmt.Println(val["info"])
  }

  time.Sleep(60 * time.Second)

}
