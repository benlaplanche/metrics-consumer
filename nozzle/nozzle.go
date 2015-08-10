package nozzle

import (
	"crypto/tls"
	"fmt"
	"github.com/benlaplanche/metrics-consumer/config"
	"github.com/cloudfoundry/noaa"
	"github.com/cloudfoundry/sonde-go/events"
	"io"
)

type MetricsNozzle struct {
	config           *config.ConsumerConfig
	errs             chan error
	messages         chan *events.Envelope
	authTokenFetcher AuthTokenFetcher
	consumer         *noaa.Consumer
	stdout           io.Writer
	stderr           io.Writer
}

type AuthTokenFetcher interface {
	FetchAuthToken() string
}

func NewNozzle(config *config.ConsumerConfig, tokenFetcher AuthTokenFetcher, stdout io.Writer, stderr io.Writer) *MetricsNozzle {
	return &MetricsNozzle{
		config:           config,
		errs:             make(chan error),
		messages:         make(chan *events.Envelope),
		authTokenFetcher: tokenFetcher,
		stdout:           stdout,
		stderr:           stderr,
	}
}

func (m *MetricsNozzle) Start() {
	var authToken string

	authToken = m.authTokenFetcher.FetchAuthToken()

	m.consumeFirehose(authToken)
	fmt.Println("**Started consuming the firehose**")
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
	if m.config.OriginID != "" {
		if envelope.GetOrigin() == m.config.OriginID {
			fmt.Fprintf(m.stdout, "%v \n", envelope)
		}
	} else {
		fmt.Fprintf(m.stdout, "%v \n", envelope)
	}

}

func (m *MetricsNozzle) handleError(err error) {
	fmt.Fprintf(m.stderr, "%v \n", err.Error())
	m.consumer.Close()
}
