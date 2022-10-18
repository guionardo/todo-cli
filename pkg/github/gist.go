package github

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/guionardo/todo-cli/pkg/logger"
)

// https://docs.github.com/pt/rest/gists/gists#list-gists-for-the-authenticated-user

const (
	CollectionFileName     = "todo-cli.yaml"
	DefaultGistDescription = "TODO CLI CONFIG"
)

type GitHubGistAPI struct {
	Authorization     string
	GistId            string
	GistDescription   string
	ConfigFileRawURL  string
	ConfigFileContent []byte
	UpdatedAt         time.Time
	client            *http.Client
	enabled           rune
	check             chan rune
	LastError         error
}

func NewGitHubGistAPI(authorization string) *GitHubGistAPI {
	api := &GitHubGistAPI{
		Authorization:   authorization,
		GistDescription: DefaultGistDescription,
		client:          &http.Client{},
		check:           make(chan rune),
		enabled:         ' ',
	}
	go api.checkAuth()
	return api
}

func (api *GitHubGistAPI) request(method string, url string, debugMsg string, body []byte) (content []byte, err error) {
	time_start := time.Now()

	defer func() {
		logger.Debugf("%s took %v", debugMsg, time.Since(time_start))
	}()

	logger.Debugf("GetToDoConfigFileGist")
	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	req.Header = http.Header{
		"Authorization": {fmt.Sprintf("Bearer %s", api.Authorization)},
		"Accept":        {"application/vnd.github.v3+json"},
	}

	res, err := api.client.Do(req)
	if err != nil {
		debugMsg = fmt.Sprintf("Error: %s - %s", err.Error(), debugMsg)
		return nil, err
	}
	if res.StatusCode >= 300 {
		debugMsg = fmt.Sprintf("Error: %s - %s", debugMsg, res.Status)
		return nil, errors.New(res.Status)
	}
	return ioutil.ReadAll(res.Body)
}

func (api *GitHubGistAPI) checkAuth() {
	if api.Authorization == "" {
		api.check <- 'N' // No authorization
		return
	}
	if err := api.ValidateAuth(api.Authorization); err != nil {
		api.LastError = err
		api.check <- 'X' // Invalid authorization
	} else {
		api.check <- 'Y' // Valid authorization
	}
}
func (api *GitHubGistAPI) Enabled() bool {
	if api.enabled == ' ' {
		api.enabled = <-api.check
	}
	return api.enabled == 'Y'
}

func (api *GitHubGistAPI) ValidateAuth(authorization string) error {
	if len(authorization) == 0 {
		return errors.New("Authorization is empty")
	}
	_, err := api.request("GET", "https://api.github.com/gists", "ValidateAuth", nil)
	if err != nil {
		err = fmt.Errorf("Invalid token: %s", err.Error())
	}
	// if len(api.ConfigFileRawURL) == 0 {
	// 	err = api.GetToDoConfigFileGist()
	// }
	return err
}

func (api *GitHubGistAPI) GetToDoConfigFileGist() error {
	if !api.Enabled() {
		return errors.New("")
	}
	if api.Authorization == "" {
		return fmt.Errorf("Authorization is empty")
	}
	logger.Debugf("GetToDoConfigFileGist")
	content, err := api.request("GET", "https://api.github.com/gists", "GetToDoConfigFileGist", nil)
	if err != nil {
		return err
	}
	var gistList []GistResponse
	err = json.Unmarshal(content, &gistList)
	if err != nil {
		return err
	}

	var configGist *GistResponse
	for _, gist := range gistList {
		if gist.Description == api.GistDescription && len(gist.Files) == 1 {
			if gistFile, ok := gist.Files[CollectionFileName]; ok {
				api.GistId = gist.Id
				api.ConfigFileRawURL = gistFile.RawURL
				api.UpdatedAt = gist.UpdatedAt
				configGist = &gist
				break
			}
		}
	}
	if configGist == nil {
		logger.Debugf("Config GIST not found")
		return fmt.Errorf("Config gist not found")
	}
	logger.Debugf("Config GIST found: %s @ %v", configGist.Description, configGist.UpdatedAt)
	return nil

}

func (api *GitHubGistAPI) GetConfigFileContent() error {
	if api.Authorization == "" {
		return fmt.Errorf("Authorization is empty")
	}
	if api.ConfigFileRawURL == "" {
		return fmt.Errorf("ConfigFileRawURL is empty")
	}
	logger.Debugf("GetConfigFileContent")
	req, err := http.NewRequest("GET", api.ConfigFileRawURL, nil)
	req.Header = http.Header{
		"Authorization": {fmt.Sprintf("Bearer %s", api.Authorization)},
		"Accept":        {"application/vnd.github.v3+json"},
	}
	if err != nil {
		return err
	}
	res, err := api.client.Do(req)
	if err != nil {
		return err
	}
	body, err := ioutil.ReadAll(res.Body)
	if err == nil {
		api.ConfigFileContent = body
		logger.Debugf("Config file content:\n%s", string(body))
	} else {
		api.ConfigFileContent = nil
	}

	return err
}

// curl \
//   -X PATCH \
//   -H "Accept: application/vnd.github+json" \
//   -H "Authorization: Bearer <YOUR-TOKEN>" \
//   https://api.github.com/gists/GIST_ID \
//   -d '{"description":"An updated gist description","files":{"README.md":{"content":"Hello World from GitHub"}}}'

func (api *GitHubGistAPI) SetConfigFileGist(configFileName string) error {
	if api.Authorization == "" {
		return fmt.Errorf("Authorization is empty")
	}
	content, err := os.ReadFile(configFileName)
	if err != nil {
		return err
	}
	gist := GistRequest{
		Description: api.GistDescription,
		Public:      false,
		Files: map[string]FileRequest{
			CollectionFileName: {
				Content: string(content),
			},
		},
	}
	body, err := json.Marshal(gist)
	if err != nil {
		return err
	}
	method := "POST"
	url := "https://api.github.com/gists"
	if api.GistId != "" {
		method = "PATCH"
		url = fmt.Sprintf("https://api.github.com/gists/%s", api.GistId)
	}
	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	logger.Debugf("SetConfigFileGist: %s %s", api.GistDescription, configFileName)
	req.Header = http.Header{
		"Authorization": {fmt.Sprintf("Bearer %s", api.Authorization)},
		"Accept":        {"application/vnd.github.v3+json"},
	}
	res, err := api.client.Do(req)
	if err != nil {
		return err
	}
	if res.StatusCode >= 400 {
		return fmt.Errorf("Invalid status code: %d", res.StatusCode)
	}
	var gistResponse GistResponse
	err = json.NewDecoder(res.Body).Decode(&gistResponse)
	if err != nil {
		return err
	}
	api.GistId = gistResponse.Id
	api.ConfigFileRawURL = gistResponse.Files[CollectionFileName].RawURL
	api.UpdatedAt = gistResponse.UpdatedAt
	logger.Debugf("Config GIST created: %s @ %v", gistResponse.Description, gistResponse.UpdatedAt)
	return nil
}

func (api *GitHubGistAPI) DeleteGist() error {
	if api.Authorization == "" {
		return errors.New("Authorization is empty")
	}
	if api.GistId == "" {
		return errors.New("GistId is empty")
	}
	req, err := http.NewRequest("DELETE", fmt.Sprintf("https://api.github.com/gists/%s", api.GistId), nil)
	if err != nil {
		return err
	}
	req.Header = http.Header{
		"Authorization": {fmt.Sprintf("Bearer %s", api.Authorization)},
		"Accept":        {"application/vnd.github.v3+json"},
	}
	res, err := api.client.Do(req)
	if res.StatusCode >= 400 {
		return fmt.Errorf("Invalid status code: %d %v", res.StatusCode, err)
	}
	api.GistId = ""
	return nil
}
