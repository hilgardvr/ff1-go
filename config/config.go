package config

import (
	"encoding/json"
	"io/ioutil"
)

type Config struct {
    AppPort string `json:"APP_PORT"`
    Neo4jUri string `json:"NEO4J_URI"`
    Neo4jUsername string `json:"NEO4J_USERNAME"`
    Neo4jPassword string `json:"NEO4J_PASSWORD"`
    EmailPassword string `json:"APP_PASSWORD"`
    EmailFrom string `json:"EMAIL_FROM"`
	SendEmails bool `json:"SEND_EMAILS"`
}

func ReadConfig(path string) (*Config, error) {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	config := Config{}
	err = json.Unmarshal(file, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}