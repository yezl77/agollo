package main

import (
  "net/http"
  "strings"
  "fmt"
  "encoding/json"
  "gopkg.in/yaml.v2"
  "time"
)

type Response struct {
  NamespaceName  string                 `json:"namespaceName"`
  Configurations map[string]interface{} `json:"configurations"`
  ReleaseKey     string                 `json:"releaseKey"`
}
type ResponseYaml struct {
  Name    string `yaml:"name"`
  Timeout int    `yaml:"timeout"`
  Spouse  Spouse `yaml:"spouse"`
}

type Spouse struct {
  Name string        `yaml:"name"`
  Size int           `yaml:"size"`
  Info []interface{} `yaml:"info"`
}

type notification struct {
  NamespaceName  string `json:"namespaceName,omitempty"`
  NotificationID int    `json:"notificationId,omitempty"`
}

var jsonid = 126

func IndexHandler(w http.ResponseWriter, r *http.Request) {
  fmt.Println(r.URL.Path)

  if strings.Contains(r.URL.Path, "/configs/app-apollo-demo/default/testjson.json") {
    m := map[string]interface{}{"enabled": true, "path": "  ", "test": "", "number": 888}
    mjson, _ := json.Marshal(m)
    mString := string(mjson)
    resp := Response{
      NamespaceName:  "testjson.json",
      Configurations: map[string]interface{}{"content": mString},
      ReleaseKey:     "20180802113631-1d4f94f06a312155",
    }
    data, _ := json.Marshal(resp)
    fmt.Fprintln(w, string(data))

  } else if strings.Contains(r.URL.Path, "/configs/app-apollo-demo/default/application") {
    m := map[string]interface{}{"apollo": "admin"}
    mjson, _ := json.Marshal(m)
    mString := string(mjson)
    resp := Response{
      NamespaceName:  "application",
      Configurations: map[string]interface{}{"content": mString},
      ReleaseKey:     "20180802113631-1d4f94f06a312154",
    }
    data, _ := json.Marshal(resp)
    fmt.Fprintln(w, string(data))
  } else if strings.Contains(r.URL.Path, "/configs/app-apollo-demo/default/testyaml.yaml") {
    infoarray := make([]interface{}, 2)
    infoarray[0] = 2048
    infoarray[1] = "byte"

    testyaml := ResponseYaml{Name: "root",
      Timeout: jsonid,
      Spouse: Spouse{Name: "spouseadmin", Size: 1024, Info: infoarray},
    }
    myaml, _ := yaml.Marshal(testyaml)
    mString := string(myaml)
    resp := Response{
      NamespaceName:  "testyaml.yaml",
      Configurations: map[string]interface{}{"content": mString},
      ReleaseKey:     "20180802113631-1d4f94f06a312156",
    }
    data, _ := json.Marshal(resp)
    fmt.Fprintln(w, string(data))
  } else if strings.Contains(r.URL.Path, "/notifications/v2") {
    notif1 := notification{
      NamespaceName:  "testjson.json",
      NotificationID: 54,
    }
    notif2 := notification{
      NamespaceName:  "application",
      NotificationID: 105,
    }
    notif3 := notification{
      NamespaceName:  "testyaml.yaml",
      NotificationID: jsonid,
    }
    jsonid++

    var dd = make([]notification, 3)
    dd[0] = notif1
    dd[1] = notif2
    dd[2] = notif3
    data, _ := json.Marshal(dd)
    time.Sleep(3 * time.Second)
    fmt.Fprintln(w, string(data))
  }

}

func main() {
  var srv = &http.Server{Addr: ":8080"}
  http.HandleFunc("/", IndexHandler)
  srv.ListenAndServe()
}
