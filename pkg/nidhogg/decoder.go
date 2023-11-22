package nidhogg

import (
	"fmt"
	"os"

	yaml "gopkg.in/yaml.v1"
)

// GetConfig reads the config file, parses it whether it be in json or yaml and returns a handler config
func GetConfig(config string) (HandlerConfig, error) {

	var handlerConf HandlerConfig
	bytes, err := os.ReadFile(config)
	if err != nil {
		return HandlerConfig{}, fmt.Errorf("unable to read config file: %v", err)
	}

	err = yaml.Unmarshal(bytes, &handlerConf)
	if err != nil {
		return HandlerConfig{}, fmt.Errorf("error parsing file: %v", err)
	}

	if err := handlerConf.BuildSelectors(); err != nil {
		return HandlerConfig{}, err
	}

	return handlerConf, nil

}
