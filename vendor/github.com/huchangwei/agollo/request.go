package agollo

import (
  "io"
  "io/ioutil"
  "net/http"
)

// this is a static check
var _ requester = (*httprequester)(nil)

type requester interface {
  request(url string) ([]byte, error)
}

type httprequester struct {
  client *http.Client
}

func newHTTPRequester(client *http.Client) requester {
  return &httprequester{
    client: client,
  }
}

func (h *httprequester) request(url string) ([]byte, error) {
  resp, err := h.client.Get(url)
  if nil != resp {
    defer resp.Body.Close()
  }
  if err != nil {
    return nil, err
  }
  if resp.StatusCode == http.StatusOK {
    return ioutil.ReadAll(resp.Body)
  }

  // Discard all body if status code is not 200
  io.Copy(ioutil.Discard, resp.Body)
  return nil, nil
}
