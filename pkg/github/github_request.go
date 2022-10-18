package github

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/guionardo/todo-cli/pkg/logger"
)

var client = &http.Client{}

func request(method string, url string, auth string, body []byte, debugMsg string) (content []byte, err error) {
	time_start := time.Now()
	defer func() {
		logger.Debugf("%s took %v", debugMsg, time.Since(time_start))
	}()
	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	req.Header = http.Header{
		"Authorization": {fmt.Sprintf("Bearer %s", auth)},
		"Accept":        {"application/vnd.github.v3+json"},
	}
	res, err := client.Do(req)
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
