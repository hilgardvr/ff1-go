package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

type Config struct {
    AppPort string `json:"APP_PORT"`
    Neo4jUri string `json:"NEO4J_URI"`
    Neo4jUsername string `json:"NEO4J_USERNAME"`
    Neo4jPassword string `json:"NEO4J_PASSWORD"`
    EmailPassword string `json:"APP_PASSWORD"`
    EmailFrom string `json:"EMAIL_FROM"`
	SendEmails bool `json:"SEND_EMAILS"`
	UpdateMode bool `json:"UPDATE_MODE"`
}

var AppConfig Config

func ReadConfig(path string) (*Config, error) {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return readEnvConfig()
	} else {
		config := Config{}
		err = json.Unmarshal(file, &config)
		if err != nil {
			return nil, err
		}
		AppConfig = config
		return &config, nil
	}
}

func readEnvConfig() (*Config, error) {
	config := Config{}
	for _, e := range os.Environ() {
		pair := strings.SplitN(e, "=", 2)
		switch pair[0] {
			case "APP_PORT": config.AppPort = pair[1]
			case "NEO4J_URI": config.Neo4jUri = pair[1]
			case "NEO4J_USERNAME": config.Neo4jUsername = pair[1]
			case "NEO4J_PASSWORD": config.Neo4jPassword = pair[1]
			case "APP_PASSWORD": config.EmailPassword = pair[1]
			case "EMAIL_FROM": config.EmailFrom = pair[1]
			case "SEND_EMAILS": {
				toggle, err := strconv.ParseBool(pair[1])
				if err != nil {
					return &config, err
				}
				config.SendEmails = toggle
			}
		default: 
		}
	}
	return &config, nil
}