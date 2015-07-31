package main

import (
	"flag"
	"github.com/benlaplanche/metrics-consumer/config"
	"github.com/benlaplanche/metrics-consumer/nozzle"
	"github.com/cloudfoundry-incubator/datadog-firehose-nozzle/uaatokenfetcher"
	"log"
	"os"
)

func main() {
	configFilePath := flag.String("config", "config/config.json", "Path to config params for the consumer")
	flag.Parse()

	config, err := config.Parse(*configFilePath)
	if err != nil {
		log.Fatalf("Error parsing config: %s", err.Error())
	}

	tokenFetcher := &uaatokenfetcher.UAATokenFetcher{
		UaaUrl:                config.UAAURL,
		Username:              config.Username,
		Password:              config.Password,
		InsecureSSLSkipVerify: config.InsecureSSLSkipVerify,
	}

	metrics_nozzle := nozzle.NewNozzle(config, tokenFetcher, os.Stdout, os.Stderr)
	metrics_nozzle.Start()

}
