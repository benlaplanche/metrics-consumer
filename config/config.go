package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type ConsumerConfig struct {
	UAAURL                 string `json:"UAAURL"`
	Username               string `json:"Username"`
	Password               string `json:"Password"`
	InsecureSSLSkipVerify  bool   `json:"InsecureSSLSkipVerify"`
	DopplerAddr            string `json:"DopplerAddr"`
	FirehoseSubscriptionId string `json:"FirehoseSubscriptionId"`
	OriginID               string `json:"OriginID,omitempty"`
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
