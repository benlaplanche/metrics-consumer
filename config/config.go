package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type ConsumerConfig struct {
	UAAURL                 string
	Username               string
	Password               string
	InsecureSSLSkipVerify  bool
	DopplerAddr            string
	FirehoseSubscriptionId string
	OriginID               string
}

func Parse(configPath string) (*ConsumerConfig, error) {
	configBytes, err := ioutil.ReadFile(configPath)
	var config ConsumerConfig

	if err != nil {
		return nil, fmt.Errorf("Can not read config file [%s]: %s", configPath, err)
	}

	err = json.Unmarshal(configBytes, &config)
	if err != nil {
		return nil, fmt.Errorf("Can not parse config file [%s]: %s", configPath, err)
	}

	return &config, nil
}
