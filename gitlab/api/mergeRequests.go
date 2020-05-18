package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"path"

	"gitlab-metrics/gitlab/dto"
	"gitlab-metrics/log"
)

const (
	projectMergeRequestsPath = "/projects/%d/merge_requests"
	mergeRequestsPath        = "/merge_requests"
)

type MergeRequestsApi interface {
	GetMergeRequestsWithOpts(opts ProjectMergeRequestOpts) ([]dto.MergeRequest, error)
	GetProjectMergeRequestsWithOpts(projectId uint32, opts ProjectMergeRequestOpts) ([]dto.MergeRequest, error)
}

type mergeRequestsApi struct {
	log.Loggable
	api
}

func (m *mergeRequestsApi) GetMergeRequestsWithOpts(opts ProjectMergeRequestOpts) ([]dto.MergeRequest, error) {
	urlPath := path.Join(apiPath, mergeRequestsPath)
	urlValues := opts.ToValues()
	reqUrl := fmt.Sprintf("%s%s?%s", m.GetBaseUrl(), urlPath, urlValues.Encode())

	return m.fetchMergeRequests(reqUrl)
}

func (m *mergeRequestsApi) GetProjectMergeRequestsWithOpts(projectId uint32, opts ProjectMergeRequestOpts) ([]dto.MergeRequest, error) {
	urlPath := fmt.Sprintf(path.Join(apiPath, projectMergeRequestsPath), projectId)
	urlValues := opts.ToValues()
	reqUrl := fmt.Sprintf("%s%s?%s", m.GetBaseUrl(), urlPath, urlValues.Encode())

	return m.fetchMergeRequests(reqUrl)
}

func (m *mergeRequestsApi) fetchMergeRequests(url string) ([]dto.MergeRequest, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("make request: %v", err)
	}

	resp, err := m.DoRequest(req)
	if err != nil {
		return nil, fmt.Errorf("do request to gitlab: %v", err)
	}
	defer resp.Body.Close()

	var mergeRequests []dto.MergeRequest

	if resp.StatusCode == 403 {
		body, _ := ioutil.ReadAll(resp.Body)
		m.Log().Warnf("access denied: accessing %s code %d '%s'", url, resp.StatusCode, body)
		return mergeRequests, nil
	}

	if resp.StatusCode != 200 {
		body, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("do request to gitlab: code %d '%s'", resp.StatusCode, body)
	}

	err = json.NewDecoder(resp.Body).Decode(&mergeRequests)
	if err != nil {
		fmt.Println(err)
		return nil, fmt.Errorf("decode response: %v", err)
	}

	return mergeRequests, nil
}
