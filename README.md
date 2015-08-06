# Metrics Consumer

Binary to consume the Cloud Foundry Metrics Firehose on a filtered origin. 

Primarily used for acceptance testing, but could be easily adapted to suit other needs.

Inspiration taken from the [DataDog Firehose Nozzle](https://github.com/cloudfoundry-incubator/datadog-firehose-nozzle)

*Tests don't fully pass yet - in process of writing a [ginkgo matcher](https://github.com/benlaplanche/metrics-matcher)*

## Configuration

Update the `config/config.json` file with the correct values for your environment

```
{
  "UAAURL": "https://uaa.10.244.0.34.xip.io",
  "Username": "admin",
  "Password": "admin",
  "DopplerAddr": "wss://doppler.10.244.0.34.xip.io:4443"
  "InsecureSSLSkipVerify": true
  "FirehoseSubscriptionID": "metrics-consumer-1"
  "OriginID": "service-metrics"
}
```

`UAAURL` - this is the UAA endpoint for your Cloud Foundry, typically `UAA.xxx`
`InsecureSSLSkipVerify` - set to true if you are using a self signed certificate

`Username` - this user needs to have the `doppler.firehose` permission. Typically the `admin` user has this as default
`Password` - password for the above user

`DopplerAddr` - this is typically `doppler.xxx`
`FirehoseSubscriptionID` - this should be unique, if you are running this binary lots of times in concurrent

`OriginID` - this is the OriginID for which you want to filter the firehose. It currently only takes a single string, it could be adapted to take an array easily. 

## Building

`go build -o metrics-consumer main.go`

## Executing

With `config.json` in the `config` dir

`./metrics-consumer`

Specifying a config file elsehwere
`./metrics-consumer --config=run.json`