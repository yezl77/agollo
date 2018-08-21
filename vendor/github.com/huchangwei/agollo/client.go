package agollo

import (
  "context"
  "encoding/json"
  "net/http"
  "strings"
  "gopkg.in/yaml.v2"
  "reflect"
)

// Client for apollo
type Client struct {
  conf *Conf

  updateChan chan *ChangeEvent

  caches         *namespaceCache
  releaseKeyRepo *cache

  longPoller poller
  requester  requester

  ctx    context.Context
  cancel context.CancelFunc
}

// result of query config
type result struct {
  NamespaceName  string        `json:"namespaceName"`
  Configurations Configuration `json:"configurations"`
  ReleaseKey     string        `json:"releaseKey"`
}

type Configuration map[string]interface{}

// NewClient create client from conf
func NewClient(conf *Conf) *Client {
  client := &Client{
    conf:           checkConf(conf),
    caches:         newNamespaceCache(),
    releaseKeyRepo: newCache(),

    requester: newHTTPRequester(&http.Client{Timeout: queryTimeout}),
  }
  client.longPoller = newLongPoller(conf, longPoolInterval,
    client.handleNamespaceUpdate)
  client.ctx, client.cancel = context.WithCancel(context.Background())
  return client
}

func checkConf(ret *Conf) *Conf {
  if len(ret.IP) == 0 {
    ret.IP = defaultIP
  }
  if len(ret.Cluster) == 0 {
    ret.Cluster = defaultCluster
  }
  if len(ret.EnvLocalPath) == 0 {
    ret.EnvLocalPath = defaultEnvLocalPath
  }
  if len(ret.NameSpaceNames) == 0 {
    ret.NameSpaceNames = make([]string, 1)
    ret.NameSpaceNames[0] = defaultNameSpaceName
  }
  return ret
}

// Start sync config
func (c *Client) Start() error {

  // preload all config to local first
  if err := c.preload(); err != nil {
    return err
  }

  // start fetch update
  go c.longPoller.start()

  return nil
}

// handleNamespaceUpdate sync config for namespace, delivery
// changes to subscriber
func (c *Client) handleNamespaceUpdate(namespace string) error {
  change, err := c.sync(namespace)
  if err != nil || change == nil {
    return err
  }

  c.deliveryChangeEvent(change)
  return nil
}

// Stop sync config
func (c *Client) Stop() error {
  c.longPoller.stop()
  c.cancel()
  // close(c.updateChan)
  c.updateChan = nil
  return nil
}

// fetchAllConfig fetch from remote, if failed ,will load from local file
func (c *Client) preload() error {
  if err := c.longPoller.preload(); err != nil {
    if err2 := c.loadLocal(c.conf.EnvLocalPath); err2 != nil {
      return err2
    }
    return err
  }
  return nil
}

// loadLocal load caches from local file
func (c *Client) loadLocal(name string) error {
  if c.conf.EnvLocal {
    return c.caches.load(name)
  }
  return nil
}

// dump caches to file
func (c *Client) dump(name string) error {
  return c.caches.dump(name)
}

// WatchUpdate get all updates
func (c *Client) WatchUpdate() <-chan *ChangeEvent {
  if c.updateChan == nil {
    c.updateChan = make(chan *ChangeEvent)
  }
  return c.updateChan
}

func (c *Client) mustGetCache(namespace string) *cache {
  return c.caches.mustGetCache(namespace)
}

// GetStringValueWithNameSpace get value from given namespace
func (c *Client) GetStringValueWithNameSpace(namespace string, key,
defaultValue interface{}) (interface{}, SourceType, error) {
  cache := c.mustGetCache(namespace)
  ret, _ := cache.get(key)
  if ret != "" && ret != nil {
    return ret, cache.getSourceType(), nil
  }
  if err := c.loadLocal(c.conf.EnvLocalPath); err != nil {
    return defaultValue, DEFAULT, err
  }
  cache = c.mustGetCache(namespace)
  ret, _ = cache.get(key)
  if ret == "" || ret == nil {
    return defaultValue, DEFAULT, nil
  }
  return ret, cache.getSourceType() , nil
}

// GetStringValue from default namespace
func (c *Client) GetStringValue(key, defaultValue interface{}) (interface{},
  SourceType, error) {
  return c.GetStringValueWithNameSpace(defaultNamespace, key, defaultValue)
}

// sync namespace config
func (c *Client) sync(namespace string) (*ChangeEvent, error) {
  releaseKey := c.getReleaseKey(namespace)
  url := configURL(c.conf, namespace, releaseKey)
  bts, err := c.requester.request(url)
  if err != nil || len(bts) == 0 {
    return nil, err
  }
  r, err := c.parse(bts)
  if err != nil {
    return nil, err
  }
  return c.handleResult(r)

}

func (c *Client) parse(bts []byte) (*result, error) {
  var result result
  if err := json.Unmarshal(bts, &result); err != nil {

    return nil, err
  }

  m := make(Configuration)
  if content, ok := result.Configurations["content"]; ok {
    contentStr, _ := content.(string)
    if strings.HasSuffix(result.NamespaceName, ".yaml") {
      if err := yaml.Unmarshal([]byte(contentStr), &m); err != nil {
        return nil, err
      }
    } else {
      if err := json.Unmarshal([]byte(contentStr), &m); err != nil {
        return nil, err
      }
    }
    result.Configurations = m
  }

  return &result, nil

}

// deliveryChangeEvent push change to subscriber
func (c *Client) deliveryChangeEvent(change *ChangeEvent) {
  if c.updateChan == nil {
    return
  }
  select {
  case <-c.ctx.Done():
  case c.updateChan <- change:
  }
}

// TODO 需要赋值 sourceType
// handleResult generate changes from query result, and update local cache
func (c *Client) handleResult(result *result) (*ChangeEvent, error) {
  var ret = ChangeEvent{
    Namespace: result.NamespaceName,
    Changes:   map[string]*Change{},
  }

  cache:= c.mustGetCache(result.NamespaceName)
  cache.setSourceType(REMOTE)
  kv := cache.dump()

  for k, v := range kv {
    if _, ok := result.Configurations[k]; !ok {
      cache.delete(k)
      ret.Changes[k] = makeDeleteChange(k, v)
    }
  }

  for k, v := range result.Configurations {
    cache.set(k, v)
    old, ok := kv[k]
    if !ok {
      ret.Changes[k] = makeAddChange(k, v)
      continue
    }
    if ! reflect.DeepEqual(old, v) {
      ret.Changes[k] = makeModifyChange(k, old, v)
    }
  }
  c.setReleaseKey(result.NamespaceName, result.ReleaseKey)

  // dump caches to file
  err := c.dump(c.conf.EnvLocalPath)

  if len(ret.Changes) == 0 {
    return nil, err
  }
  return &ret, err
}

func (c *Client) getReleaseKey(namespace string) interface{} {
  releaseKey, _ := c.releaseKeyRepo.get(namespace)
  return releaseKey
}

func (c *Client) setReleaseKey(namespace, releaseKey string) {
  c.releaseKeyRepo.set(namespace, releaseKey)
}
