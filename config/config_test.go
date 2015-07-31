package config_test

import (
	config "github.com/benlaplanche/metrics-consumer/config"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Config", func() {

	It("successfully parses a valid config file", func() {
		conf, err := config.Parse("config.json")

		Expect(err).ToNot(HaveOccurred())

		Expect(conf.UAAURL).To(Equal("https://uaa.10.244.0.34.xip.io"))
		Expect(conf.Username).To(Equal("admin"))
		Expect(conf.Password).To(Equal("admin-secret"))
		Expect(conf.DopplerAddr).To(Equal("wss://doppler.10.244.0.34.xip.io:4443"))
		Expect(conf.InsecureSSLSkipVerify).To(Equal(true))
		Expect(conf.FirehoseSubscriptionId).To(Equal("metrics-consumer-1"))

		Expect(conf.OriginID).To(Equal("service-metrics"))

	})

})
