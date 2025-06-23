package config

import (
	"encoding/json"
	"os"
)

func (conf *Config) SetUser(username string) error {
	conf.CurUser = username
	return write(conf)
}

func write(conf *Config) error {
	configFilename, err := getConfigFilename()
	if err != nil {
		return err
	}
	jsonConfig, err := json.Marshal(conf)
	if err != nil {
		return err
	}
	os.WriteFile(configFilename, []byte(jsonConfig), 0777)
	return nil
}
