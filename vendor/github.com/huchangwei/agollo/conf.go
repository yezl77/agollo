package agollo

import (
  "os"
  "gopkg.in/yaml.v2"
)

// Conf ...
type Conf struct {
  AppID          string   `yaml:"appId,omitempty"`
  Cluster        string   `yaml:"cluster,omitempty"`
  NameSpaceNames []string `yaml:"namespaceNames,omitempty"`
  IP             string   `yaml:"ip,omitempty"`
  EnvLocal       bool     `yaml:"env_local,omitempty"`
  EnvLocalPath   string   `yaml:"env_local_path,omitempty"`
}

// NewConf create Conf from file
func NewConf(name string) (*Conf, error) {
  f, err := os.Open(name)
  if err != nil {
    return nil, err
  }
  defer f.Close()

  var ret Conf
  if err := yaml.NewDecoder(f).Decode(&ret); err != nil {
    return nil, err
  }
  return &ret, nil
}
