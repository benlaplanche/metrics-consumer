package nozzle

import (
	"crypto/tls"
	"fmt"
	"github.com/benlaplanche/metrics-consumer/config"
	"github.com/cloudfoundry/noaa"
	"github.com/cloudfoundry/sonde-go/events"
	"os"
)

type MetricsNozzle struct {
	config           *config.ConsumerConfig
	errs             chan error
	messages         chan *events.Envelope
	authTokenFetcher AuthTokenFetcher
	consumer         *noaa.Consumer
}

type AuthTokenFetcher interface {
	FetchAuthToken() string
}

func NewNozzle(config *config.ConsumerConfig, tokenFetcher AuthTokenFetcher) *MetricsNozzle {
	return &MetricsNozzle{
		config:           config,
		errs:             make(chan error),
		messages:         make(chan *events.Envelope),
		authTokenFetcher: tokenFetcher,
	}
}

func (m *MetricsNozzle) Start() {
	var authToken string

	authToken = m.authTokenFetcher.FetchAuthToken()

	m.consumeFirehose(authToken)

	m.processFirehose()
}

func (m *MetricsNozzle) consumeFirehose(authToken string) {
	m.consumer = noaa.NewConsumer(
		m.config.DopplerAddr,
		&tls.Config{InsecureSkipVerify: m.config.InsecureSSLSkipVerify},
		nil)

	go m.consumer.Firehose(m.config.FirehoseSubscriptionId, authToken, m.messages, m.errs)
}

func (m *MetricsNozzle) processFirehose() {
	for {
		select {
		case envelope := <-m.messages:
			m.handleMessage(envelope)
		case err := <-m.errs:
			m.handleError(err)
			return
		}
	}
}

func (m *MetricsNozzle) handleMessage(envelope *events.Envelope) {
	if envelope.GetOrigin() == m.config.OriginID {
		fmt.Printf("%v \n", envelope)
	}

}

func (m *MetricsNozzle) handleError(err error) {
	fmt.Fprintf(os.Stderr, "%v \n", err.Error())
	m.consumer.Close()
}
