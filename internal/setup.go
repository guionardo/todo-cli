package internal

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"log"
)

type Config struct {
	ToDoListName  string
	GistId        string
	Authorization string
}

func (c *Config) Load(configFile string) error {
	file, err := os.Open(configFile)
	if err != nil {
		return err
	}
	defer file.Close()
	body, _ := ioutil.ReadAll(file)
	return json.Unmarshal(body, c)
}

func (c *Config) Save(configFile string) error {
	body, err := json.Marshal(c)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(configFile, body, 0644)
}

func ConfigFile() (string, error) {
	if todo_config := os.Getenv("TODO_CONFIG"); todo_config != "" {
		if _, err := os.Stat(todo_config); err == nil {
			_, err := ParseConfigFile(todo_config)
			if err != nil {
				log.Printf("Error parsing config file %s: %s", todo_config, err)
			} else {
				return todo_config, nil
			}
		}
	}
	return "", nil
}

func ParseConfigFile(configFile string) (string, error) {
	return "", nil
}
