package agollo

import (
  "encoding/gob"
  "os"
  "sync"
)

type namespaceCache struct {
  lock   sync.RWMutex
  caches map[string]*cache
}

func newNamespaceCache() *namespaceCache {
  return &namespaceCache{
    caches: map[string]*cache{},
  }
}



func (n *namespaceCache) mustGetCache(namespace string) *cache {
  n.lock.RLock()
  if ret, ok := n.caches[namespace]; ok {
    n.lock.RUnlock()
    return ret
  }
  n.lock.RUnlock()

  n.lock.Lock()
  defer n.lock.Unlock()

  cache := newCache()
  n.caches[namespace] = cache
  return cache
}

func (n *namespaceCache) drain() {
  for namespace := range n.caches {
    delete(n.caches, namespace)
  }
}

// 从缓存dump到本地
func (n *namespaceCache) dump(name string) error {

  dumps := make(map[string]Configuration)
  for namespace, cache := range n.caches {
    dumps[namespace] = cache.dump()
  }
  f, err := os.OpenFile(name, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
  if err != nil {
    return err
  }
  defer f.Close()
  gob.Register(map[interface{}]interface{}{})
  gob.Register([]interface{}{})
  err = gob.NewEncoder(f).Encode(&dumps)
  return err
}

// 从本地load到缓存
func (n *namespaceCache) load(name string) error {
  n.drain()

  f, err := os.OpenFile(name, os.O_RDONLY, 0755)
  if err != nil {
    return err
  }
  defer f.Close()

  dumps := make(map[string]Configuration)

  gob.Register(map[interface{}]interface{}{})
  gob.Register([]interface{}{})
  if err := gob.NewDecoder(f).Decode(&dumps); err != nil {
    return err
  }

  for namespace, kv := range dumps {
    cache := n.mustGetCache(namespace)
    for k, v := range kv {
      cache.set(k, v)
    }
    cache.setSourceType(LOCAL)
  }
  return nil
}


type cache struct {
  kv sync.Map
  sourceType SourceType
}

func newCache() *cache {
  return &cache{
    kv: sync.Map{},
    sourceType:REMOTE,
  }
}

func (c *cache) set(key, val interface{}) {
  c.kv.Store(key, val)
}

func (c *cache) get(key interface{}) (interface{}, bool) {
  if val, ok := c.kv.Load(key); ok {
    return val, true
  }
  return nil, false
}

func (c *cache) setSourceType(sourceType SourceType){
  c.sourceType=sourceType
}

func (c *cache) getSourceType()SourceType{
  return c.sourceType
}

func (c *cache) delete(key interface{}) {
  c.kv.Delete(key)
}

func (c *cache) dump() Configuration {
  var ret = Configuration{}
  c.kv.Range(func(key, val interface{}) bool {
    k, ok := key.(string)
    if ok {
      ret[k] = val
    }
    return true
  })
  return ret
}
