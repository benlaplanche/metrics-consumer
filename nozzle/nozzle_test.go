package nozzle_test

import (
	config "github.com/benlaplanche/metrics-consumer/config"
	nozzle "github.com/benlaplanche/metrics-consumer/nozzle"
	. "github.com/cloudfoundry-incubator/datadog-firehose-nozzle/testhelpers"

	"fmt"
	"github.com/cloudfoundry-incubator/datadog-firehose-nozzle/uaatokenfetcher"
	"github.com/cloudfoundry/sonde-go/events"
	"github.com/gogo/protobuf/proto"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"log"
	"strings"
)

var _ = Describe("Nozzle", func() {
	var fakeUAA *FakeUAA
	var fakeFirehose *FakeFirehose
	var configuration *config.ConsumerConfig
	var metrics_nozzle *nozzle.MetricsNozzle
	var logOutput *gbytes.Buffer

	BeforeEach(func() {
		fakeUAA = NewFakeUAA("bearer", "123456789")
		fakeToken := fakeUAA.AuthToken()
		fakeFirehose = NewFakeFirehose(fakeToken)

		fakeUAA.Start()
		fakeFirehose.Start()

		tokenFetcher := &uaatokenfetcher.UAATokenFetcher{
			UaaUrl: fakeUAA.URL(),
		}

		configuration = &config.ConsumerConfig{
			UAAURL:                 fakeUAA.URL(),
			Username:               "admin",
			Password:               "admin",
			InsecureSSLSkipVerify:  true,
			DopplerAddr:            strings.Replace(fakeFirehose.URL(), "http:", "ws:", 1),
			FirehoseSubscriptionId: "metrics-consumer-tests",
			OriginID:               "service-metrics-test",
		}

		logOutput = gbytes.NewBuffer()
		log.SetOutput(logOutput)

		metrics_nozzle = nozzle.NewNozzle(configuration, tokenFetcher, GinkgoWriter, GinkgoWriter)
	})

	JustBeforeEach(func() {
		// adding data with origin "service-metrics-test"
		for i := 0; i < 10; i++ {
			envelope := events.Envelope{
				Origin:    proto.String("service-metrics-test"),
				Timestamp: proto.Int64(1000000000),
				EventType: events.Envelope_ValueMetric.Enum(),
				ValueMetric: &events.ValueMetric{
					Name:  proto.String(fmt.Sprintf("filteredMetric-%d", i)),
					Value: proto.Float64(float64(i)),
					Unit:  proto.String("gauge"),
				},
				Deployment: proto.String("service-metrics-test"),
				Job:        proto.String("metrics"),
			}
			fakeFirehose.AddEvent(envelope)
		}

		// adding data with origin "other"
		for i := 0; i < 10; i++ {
			envelope := events.Envelope{
				Origin:    proto.String("other"),
				Timestamp: proto.Int64(1000000000),
				EventType: events.Envelope_ValueMetric.Enum(),
				ValueMetric: &events.ValueMetric{
					Name:  proto.String(fmt.Sprintf("unfilteredMetric-%d", i)),
					Value: proto.Float64(float64(i)),
					Unit:  proto.String("gauge"),
				},
				Deployment: proto.String("cloudfoundry"),
				Job:        proto.String("doppler"),
			}
			fakeFirehose.AddEvent(envelope)
		}
	})

	It("receives filtered data from the firehose", func(done Done) {
		defer close(done)

		go metrics_nozzle.Start()
		Eventually(GinkgoWriter).Should(ContainSubstring("filteredMetric-1"))

	}, 2)

	AfterEach(func() {
		fakeUAA.Close()
		fakeFirehose.Close()
	})

})
