package configcenter

import (
  "sync"
  "github.com/huchangwei/agollo"
)

type CallBackFunc func(oldValue, newValue interface{},
  changeType string) error

type configInstance map[string]CallBackFunc

var (
  // the value is from REMOTE
  REMOTE = agollo.REMOTE
  // LOCAL a old value
  LOCAL = agollo.LOCAL
  // DEFAULT ...
  DEFAULT = agollo.DEFAULT
)

type ConfigCenter struct {
  sync.RWMutex
  watches   map[string]configInstance
  watchChan <-chan *agollo.ChangeEvent
  stopChan  chan struct{}
}

func (c *ConfigCenter) Init(appConfigPath string) error {

  err := agollo.StartWithConfFile(appConfigPath)
  if err != nil {
    return err
  }

  c.stopChan = make(chan struct{})
  c.watches = make(map[string]configInstance)
  c.watchChan = agollo.WatchUpdate()

  go c.watchConfigUpdatesProc()

  return err
}

func (c *ConfigCenter) InitWithConf(conf *agollo.Conf) error {
  err := agollo.StartWithConf(conf)
  if err != nil {
    return err
  }

  c.stopChan = make(chan struct{})
  c.watches = make(map[string]configInstance)
  c.watchChan = agollo.WatchUpdate()

  go c.watchConfigUpdatesProc()

  return err
}

func (c *ConfigCenter) UnInit() {
  close(c.stopChan)
  agollo.Stop()
}

// key 为 指定的监控 key
func (c *ConfigCenter) RegisterKeyWatchFuncDefault(key string,
  callback CallBackFunc) {
  defaultNamespace := "application"
  c.RegisterKeyWatchFunc(defaultNamespace, key, callback)
}

func (c *ConfigCenter) RegisterKeyWatchFunc(namespace, key string,
  callback CallBackFunc) {
  c.Lock()
  defer c.Unlock()

  instance := make(map[string]CallBackFunc)
  instance[key] = callback

  c.watches[namespace] = instance
}

func (c *ConfigCenter) GetConfigValue(key, defaultValue string) (interface{},agollo.SourceType, error) {
  return agollo.GetStringValue(key, defaultValue)
}

func (c *ConfigCenter) GetConfigValueWithNameSpace(namespace string, key,
defaultValue interface{}) (interface{}, agollo.SourceType, error)  {
  return agollo.GetStringValueWithNameSpace(namespace, key, defaultValue)
}

func (c *ConfigCenter) watchConfigUpdatesProc() {
  for {
    select {
    case <-c.stopChan:
      break
    case updates := <-c.watchChan:
      c.triggerConfigInstanceCallBack(updates)
    }
  }
}

func (c *ConfigCenter) triggerConfigInstanceCallBack(
  updates *agollo.ChangeEvent) {
  upNameSpace := updates.Namespace
  c.RLock()
  cfgInstances, ok := c.watches[upNameSpace]
  c.RUnlock()

  if !ok {
    return
  }
  changes := updates.Changes
  for uKey, upItem := range changes {
    for key, callback := range cfgInstances {
      if uKey == key {
        go callback(upItem.OldValue, upItem.NewValue, upItem.ChangeType.String())
      }
    }
  }

}
