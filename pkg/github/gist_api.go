package github

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/guionardo/todo-cli/pkg/logger"
)

type GistAPI struct {
	Config         *GistConfig
	GistContent    []byte
	UpdatedAt      time.Time
	DisabledStatus string	
}

func (api *GistAPI) checkStatus(err error) error {
	if err == nil && len(api.Config.Authorization) == 0 {
		err = errors.New("Authorization is empty")
	}
	if err != nil {
		if strings.Contains(err.Error(), "401") {
			err = errors.New(err.Error() + " - Invalid token")
		}
		if strings.Contains(err.Error(), "404") {
			err = errors.New(err.Error() + " - Gist not found")
		}
		api.DisabledStatus = err.Error()
	}
	return err
}

func NewGistAPI(config *GistConfig) (api *GistAPI, err error) {
	if config.GistDescription == "" {
		config.GistDescription = DefaultGistDescription
	}
	api = &GistAPI{Config: config}
	if err = api.checkStatus(nil); err != nil {
		return
	}

	if len(config.RawURL) > 0 {
		body, err := downloadGistFile(config.RawURL, config.Authorization)
		if err == nil {
			api.GistContent = body
			api.UpdatedAt = time.Now()
			logger.Debugf("Gist downloaded from rawurl: %s", config.RawURL)
			return api, nil
		}
		err = api.checkStatus(err)
		logger.Warnf("Error downloading gist from rawurl: %s - %v", config.RawURL, err)
	}
	gists, err := getGists(config.Authorization)

	if err != nil {
		err = api.checkStatus(err)
		return
	}

	for _, gist := range gists {
		if gist.Description == config.GistDescription {
			if file, ok := gist.Files[CollectionFileName]; ok {
				api = &GistAPI{
					Config: config,
				}
				api.Config.RawURL = file.RawURL
				api.Config.GistId = gist.Id
				if api.GistContent, err = downloadGistFile(file.RawURL, config.Authorization); err == nil {
					api.UpdatedAt = gist.UpdatedAt
					logger.Debugf("Gist downloaded from gist: %s", gist.Id)
					return
				}
			}
		}
	}
	api = &GistAPI{
		Config: &GistConfig{
			Authorization:   config.Authorization,
			GistDescription: config.GistDescription,
			GistId:          "",
			RawURL:          "",
		},
		GistContent: []byte(""),
	}
	return
}

func downloadGistFile(url string, auth string) (body []byte, err error) {
	return request("GET", url, auth, nil, "DownloadGistFile")
}

func getGists(auth string) (gists []GistResponse, err error) {
	var content []byte
	if content, err = request("GET", "https://api.github.com/gists", auth, nil, "GetGists"); err == nil {
		err = json.Unmarshal(content, &gists)
	}
	return
}

func (api *GistAPI) Save(gistContent []byte) (err error) {
	if err = api.checkStatus(nil); err != nil {
		return
	}

	gist := GistRequest{
		Description: api.Config.GistDescription,
		Public:      false,
		Files: map[string]FileRequest{
			CollectionFileName: {
				Content: string(gistContent),
			},
		},
	}
	body, err := json.Marshal(gist)
	if err != nil {
		return err
	}
	method, action := "POST", "Created"

	url := "https://api.github.com/gists"
	if api.Config.GistId != "" {
		method, action = "PATCH", "Updated"
		url = fmt.Sprintf("https://api.github.com/gists/%s", api.Config.GistId)
	}
	var content []byte
	if content, err = request(method, url, api.Config.Authorization, body, "SaveGist"); err != nil {
		return
	}
	var gistResponse GistResponse
	if err = json.Unmarshal(content, &gistResponse); err != nil {
		return
	}
	api.Config.GistId = gistResponse.Id
	api.UpdatedAt = gistResponse.UpdatedAt
	api.GistContent = gistContent
	api.Config.RawURL = gistResponse.Files[CollectionFileName].RawURL
	logger.Debugf("Gist %s: %s", action, gistResponse.Id)
	return
}

func (api *GistAPI) Delete() (err error) {
	if err = api.checkStatus(nil); err != nil {
		return
	}
	if len(api.Config.GistId) == 0 {
		return errors.New("GistId is empty")
	}

	if _, err = request("DELETE", fmt.Sprintf("https://api.github.com/gists/%s", api.Config.GistId), api.Config.Authorization, nil, "DeleteGist"); err != nil {
		return
	}
	logger.Debugf("Gist Deleted: %s", api.Config.GistId)
	api.Config.GistId = ""
	api.UpdatedAt = time.Time{}
	api.GistContent = nil
	api.Config.RawURL = ""
	return
}


// logger.Debugf("GetToDoConfigFileGist")
// content, err := api.request("GET", "https://api.github.com/gists", "GetToDoConfigFileGist", nil)
// if err != nil {
// 	return err
// }
// var gistList []GistResponse
// err = json.Unmarshal(content, &gistList)
// if err != nil {
// 	return err
// }

// var configGist *GistResponse
// for _, gist := range gistList {
// 	if gist.Description == api.GistDescription && len(gist.Files) == 1 {
// 		if gistFile, ok := gist.Files[CollectionFileName]; ok {
// 			api.GistId = gist.Id
// 			api.ConfigFileRawURL = gistFile.RawURL
// 			api.UpdatedAt = gist.UpdatedAt
// 			configGist = &gist
// 			break
// 		}
// 	}
// }
// if configGist == nil {
// 	logger.Debugf("Config GIST not found")
// 	return fmt.Errorf("Config gist not found")
// }
// logger.Debugf("Config GIST found: %s @ %v", configGist.Description, configGist.UpdatedAt)
// return nil
