package configcenter

import (
  "testing"
  "fmt"
  "time"
  "apollo_go/configcenter"
)

//测试前先启动模拟的配置中心HttpConfigServer

func TestConfigCenter_GetConfigValueWithNameSpace_json(t *testing.T) {

  cfgCenter := new(ConfigCenter)
  appConfigPath := "testapp.yaml"

  err := cfgCenter.Init(appConfigPath)
  if err != nil {
    fmt.Println("cfgCenter init error:", err)
    t.Error("测试开启失败")
    return
  }
  value, sourceType,  err := cfgCenter.GetConfigValueWithNameSpace("testjson.json",
    "number", "default")
  fmt.Println(sourceType.String())
  fmt.Println("number:",value)
  if value == float64(888) {
    t.Log("测试json成功")
  } else {
    t.Error("测试json失败")
  }
  value, sourceType,  err = cfgCenter.GetConfigValueWithNameSpace("testjson.json",
    "test", "default")
  fmt.Println("test:",value)
  if value == "default" {
    t.Log("测试key = test成功")
  } else {
    t.Error("测试key = test失败")
  }


}

//测试前先启动模拟的配置中心HttpConfigServer
func TestConfigCenter_GetConfigValueWithNameSpace_properties(t *testing.T) {

  cfgCenter := new(ConfigCenter)
  appConfigPath := "testapp.yaml"

  err := cfgCenter.Init(appConfigPath)
  if err != nil {
    fmt.Println("cfgCenter init error:", err)
    return
  }
  value, sourceType,  err := cfgCenter.GetConfigValue("apollo", "default")
  fmt.Println(sourceType.String())
  if value == "admin" {
    t.Log("测试properties成功")
  } else {
    fmt.Println(value)
    t.Error("测试properties失败")
  }

}

//测试前先启动模拟的配置中心HttpConfigServer
func TestConfigCenter_GetConfigValueWithNameSpace_yaml(t *testing.T) {

  cfgCenter := new(ConfigCenter)
  appConfigPath := "testapp.yaml"

  err := cfgCenter.Init(appConfigPath)
  if err != nil {
    fmt.Println("cfgCenter init error:", err)
    return
  }
  value, sourceType,  err := cfgCenter.GetConfigValueWithNameSpace("testyaml.yaml",
    "spouse", "default")
  fmt.Println(sourceType.String())
  fmt.Println(value)
  if val, ok := value.(map[interface{}]interface{}); ok {
    info := val["info"]
    if byte := info.([]interface{}); ok {
      if byte[0] == 2048 {
        t.Log("测试yaml成功")
      }
    }
  } else {
    t.Error("测试yaml失败")
  }

}

//测试前先启动模拟的配置中心HttpConfigServer
func TestConfigCenter_GetConfigValueWithNameSpace_yaml_update(t *testing.T) {
  cfgCenter := new(ConfigCenter)
  appConfigPath := "testapp.yaml"

  err := cfgCenter.Init(appConfigPath)
  if err != nil {
    fmt.Println("cfgCenter init error:", err)
    return
  }
  cfgCenter.RegisterKeyWatchFunc("testyaml.yaml", "timeout",
    func(oldValue, newValue interface{}, changeType string) error {
      fmt.Printf(" oldValue:%s,\n newValue:%s, \n changeType:%s\n",
        oldValue, newValue, changeType)
      return nil
    })

  value, sourceType,  err := cfgCenter.GetConfigValueWithNameSpace("testyaml.yaml",
    "spouse", "default")
  fmt.Println(sourceType.String())
  fmt.Println(value)
  if val, ok := value.(map[interface{}]interface{}); ok {
    info := val["info"]
    if bytes := info.([]interface{}); ok {
      if bytes[0] == 2048 {
        t.Log("测试yaml_update成功")
      }
    }
  } else {
    t.Error("测试yaml_update失败")
  }
  time.Sleep(5 * time.Second)
}

func TestConfigCenter_appyaml(t *testing.T) {

  cfgCenter := new(ConfigCenter)
  appConfigPath := "testapp2.yaml"

  err := cfgCenter.Init(appConfigPath)
  if err != nil {
    fmt.Println("cfgCenter init error:", err)
    return
  }
  value, sourceType,  err := cfgCenter.GetConfigValue("apollo", "default")
  fmt.Println(sourceType.String())
  if value == "admin" {
    t.Log("测试testapp2.yaml 配置项部分缺失成功")
  } else {
    fmt.Println(value)
    t.Error("测试testapp2.yaml 配置项部分缺失 失败")
  }
}

func TestConfigCenter_GetConfigValue_SourceType(t *testing.T) {
  cfgCenter := new(ConfigCenter)
  appConfigPath := "testapp.yaml"

  err := cfgCenter.Init(appConfigPath)
  if err != nil {
    fmt.Println("cfgCenter init error:", err)
    return
  }
  cfgCenter.RegisterKeyWatchFunc("testyaml.yaml", "timeout",
    func(oldValue, newValue interface{}, changeType string) error {
      fmt.Printf(" oldValue:%s,\n newValue:%s, \n changeType:%s\n",
        oldValue, newValue, changeType)
      return nil
    })

  value, sourceType,  err := cfgCenter.GetConfigValueWithNameSpace("testyaml.yaml",
    "root", "default")
  fmt.Println(sourceType.String())
  if sourceType == configcenter.DEFAULT{
    t.Log("测试sourceType成功")
  }else{
    t.Error("测试sourceType失败")
  }
  fmt.Println(value)


  time.Sleep(5 * time.Second)
  value, sourceType,  err = cfgCenter.GetConfigValueWithNameSpace("testyaml.yaml",
    "spouse", "default")
  fmt.Println(sourceType.String())
  if sourceType == configcenter.REMOTE{
    t.Log("测试sourceType成功")
  }else{
    t.Error("测试sourceType失败")
  }
  fmt.Println(value)

}

