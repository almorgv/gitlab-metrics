package api

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const (
	apiPath = "/api/v4"
)

type Api interface {
	DoRequest(req *http.Request) (*http.Response, error)
	FetchData(url string, data interface{}) error
	GetBaseUrl() string
}

type api struct {
	url   string
	token string
}

func (a *api) GetBaseUrl() string {
	return a.url
}

func (a *api) DoRequest(req *http.Request) (*http.Response, error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}
	client := &http.Client{Transport: tr}
	req.Header.Add("PRIVATE-TOKEN", a.token)
	return client.Do(req)
}

func (a *api) FetchData(url string, data interface{}) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("make request %s: %v", url, err)
	}

	resp, err := a.DoRequest(req)
	if err != nil {
		return fmt.Errorf("do request to gitlab %s: %v", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("do request to gitlab %s code %d '%s'", url, resp.StatusCode, body)
	}

	err = json.NewDecoder(resp.Body).Decode(data)
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("decode response %s: %v", url, err)
	}

	return nil
}
