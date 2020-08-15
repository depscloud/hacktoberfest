package librariesio

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

// NewClient produces a client to libraries.io with the provided apiKey.
func NewClient(apiKey string) *Client {
	return &Client{
		apiKey: apiKey,
	}
}

// LookUpResult represents a response we would get from libraries.io
type LookUpResult struct {
	HomePage      string `json:"homepage"`
	RepositoryURL string `json:"repository_url"`
}

// Client intermediates communication with libraries.io
type Client struct {
	apiKey string
}

// LookUp attempts to find a look up information about an open source library.
func (c *Client) LookUp(platform, name string) (*LookUpResult, error) {
	uri := fmt.Sprintf("https://libraries.io/api/%s/%s?api_key=%s",
		url.QueryEscape(platform), url.QueryEscape(name), url.QueryEscape(c.apiKey))

	resp, err := http.Get(uri)
	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	result := &LookUpResult{}
	if err := json.Unmarshal(data, result); err != nil {
		return nil, err
	}

	return result, nil
}
