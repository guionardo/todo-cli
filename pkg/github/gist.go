package github

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
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
	Debug             bool
}

func NewGitHubGistAPI(authorization string, debug bool) *GitHubGistAPI {
	return &GitHubGistAPI{
		Authorization:   authorization,
		GistDescription: DefaultGistDescription,
		client:          &http.Client{},
		Debug:           debug,
	}
}
func (api *GitHubGistAPI) Log(format string, v ...interface{}) {
	if api.Debug {
		log.Printf(format, v...)
	}
}

func (api *GitHubGistAPI) GetToDoConfigFileGist() error {
	if api.Authorization == "" {
		return fmt.Errorf("Authorization is empty")
	}
	api.Log("GetToDoConfigFileGist")
	req, err := http.NewRequest("GET", "https://api.github.com/gists", nil)
	if err != nil {
		return err
	}
	req.Header = http.Header{
		"Authorization": {fmt.Sprintf("Bearer %s", api.Authorization)},
		"Accept":        {"application/vnd.github.v3+json"},
	}

	res, err := api.client.Do(req)
	if err != nil {
		return err
	}
	if res.StatusCode == 403 {
		return errors.New("Invalid token")
	}
	if res.StatusCode != 200 {
		return fmt.Errorf("Invalid status code: %d", res.StatusCode)
	}
	var gistList []GistResponse
	err = json.NewDecoder(res.Body).Decode(&gistList)
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
		api.Log("Config GIST not found")
		return fmt.Errorf("Config gist not found")
	}
	api.Log("Config GIST found: %s @ %v", configGist.Description, configGist.UpdatedAt)
	return nil

}

func (api *GitHubGistAPI) GetConfigFileContent() error {
	if api.Authorization == "" {
		return fmt.Errorf("Authorization is empty")
	}
	if api.ConfigFileRawURL == "" {
		return fmt.Errorf("ConfigFileRawURL is empty")
	}
	api.Log("GetConfigFileContent")
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
		api.Log("Config file content:\n%s", string(body))
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
	api.Log("SetConfigFileGist: %s %s", api.GistDescription, configFileName)
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
	api.Log("Config GIST created: %s @ %v", gistResponse.Description, gistResponse.UpdatedAt)
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

// func GetToDoConfigFileGist(auth string) ([]byte, error) {
// 	client := &http.Client{}
// 	req, err := http.NewRequest("GET", "https://api.github.com/gists", nil)
// 	if err != nil {
// 		return nil, err
// 	}
// 	req.Header = http.Header{
// 		"Authorization": {fmt.Sprintf("Bearer %s", auth)},
// 		"Accept":        {"application/vnd.github.v3+json"},
// 	}

// 	res, err := client.Do(req)
// 	if err != nil {
// 		return nil, err
// 	}
// 	if res.StatusCode != 200 {
// 		return nil, fmt.Errorf("Invalid status code: %d", res.StatusCode)
// 	}
// 	var gistList []GistResponse
// 	err = json.NewDecoder(res.Body).Decode(&gistList)
// 	if err != nil {
// 		return nil, err
// 	}
// 	var configGist *GistResponse
// 	for _, gist := range gistList {
// 		if gist.Description == "todo config" && len(gist.Files) == 1 {
// 			if _, ok := gist.Files[CollectionFileName]; ok {
// 				configGist = &gist
// 				break
// 			}
// 		}
// 	}
// 	if configGist == nil {
// 		return nil, fmt.Errorf("Config gist not found")
// 	}

// 	req, err = http.NewRequest("GET", configGist.Files[CollectionFileName].RawURL, nil)
// 	req.Header = http.Header{
// 		"Authorization": {fmt.Sprintf("Bearer %s", auth)},
// 		"Accept":        {"application/vnd.github.v3+json"},
// 	}
// 	if err != nil {
// 		return nil, err
// 	}
// 	res, err = client.Do(req)
// 	if err != nil {
// 		return nil, err
// 	}

// 	var collection ToDoCollection
// 	err = yaml.NewDecoder(res.Body).Decode(&collection)
// 	if err != nil {
// 		return nil, err
// 	}
// 	//collection.Save()
// 	//TODO: Continuar implementação

// 	return nil, err
// }
