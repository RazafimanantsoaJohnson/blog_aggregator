package config

import (
	"encoding/json"
	"os"
)

const configFilename = "/.gatorconfig.json"

func Read() (Config, error) {
	var curConfig Config
	configFilename, err := getConfigFilename()
	if err != nil {
		return curConfig, err
	}
	fContent, err := os.ReadFile(configFilename)
	if err != nil {
		return curConfig, err
	}

	err = json.Unmarshal(fContent, &curConfig)
	if err != nil {
		return curConfig, err
	}
	return curConfig, nil
}

func getConfigFilename() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return homeDir + configFilename, nil
}
