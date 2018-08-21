package agollo

var (
  defaultClient *Client
)

// Start agollo
func Start() error {
  return StartWithConfFile(defaultConfName)
}

// StartWithConfFile run agollo with conf file
func StartWithConfFile(name string) error {
  conf, err := NewConf(name)
  if err != nil {
    return err
  }
  return StartWithConf(conf)
}

// StartWithConf run agollo with Conf
func StartWithConf(conf *Conf) error {
  defaultClient = NewClient(conf)

  return defaultClient.Start()
}

// Stop sync config
func Stop() error {
  return defaultClient.Stop()
}

// WatchUpdate get all updates
func WatchUpdate() <-chan *ChangeEvent {
  return defaultClient.WatchUpdate()
}

// GetStringValueWithNameSpace get value from given namespace
func GetStringValueWithNameSpace(namespace string, key,
defaultValue interface{}) (interface{}, SourceType, error) {
  return defaultClient.GetStringValueWithNameSpace(namespace, key, defaultValue)
}

// GetStringValue from default namespace
func GetStringValue(key, defaultValue interface{}) (interface{}, SourceType,
  error) {
  return GetStringValueWithNameSpace(defaultNamespace, key, defaultValue)
}
