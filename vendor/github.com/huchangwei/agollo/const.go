package agollo

import (
  "time"
)

const (
  defaultConfName      = "app.yaml"
  defaultNamespace     = "application"
  defaultCluster       = "default"
  defaultIP            = "127.0.0.1:8080"
  defaultNameSpaceName = "application"
  defaultEnvLocalPath  = "defaultFile"

  longPoolInterval      = time.Second * 2
  longPoolTimeout       = time.Second * 90
  queryTimeout          = time.Second * 2
  defaultNotificationID = -1
)
