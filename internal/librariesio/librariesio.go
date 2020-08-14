package librariesio

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

func NewClient(apiKey string) *Client {
	return &Client{
		apiKey: apiKey,
	}
}

type LookUpResult struct {
	HomePage      string `json:"homepage"`
	RepositoryURL string `json:"repository_url"`
}

type Client struct {
	apiKey string
}

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
