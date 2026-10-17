package nidhogg

import (
	"fmt"
	"os"

	"sigs.k8s.io/yaml"
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

	if len(handlerConf.Daemonsets) == 0 {
		return HandlerConfig{}, fmt.Errorf("config needs at least one daemonset")
	}
	switch handlerConf.TaintEffect {
	case "", "NoSchedule", "PreferNoSchedule", "NoExecute":
	default:
		return HandlerConfig{}, fmt.Errorf("unknown taintEffect %q, use NoSchedule, PreferNoSchedule or NoExecute", handlerConf.TaintEffect)
	}

	if err := handlerConf.BuildSelectors(); err != nil {
		return HandlerConfig{}, err
	}

	return handlerConf, nil

}
