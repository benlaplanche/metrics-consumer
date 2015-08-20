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

`OriginID` - this is the OriginID for which you want to filter the firehose. It currently only takes a single string, it could be adapted to take an array easily. If this is omitted or an empty string then the firehose is not filtered. 

## Building

`go build -o metrics-consumer main.go`

## Executing

With `config.json` in the `config` dir

`./metrics-consumer`

Specifying a config file elsehwere
`./metrics-consumer --config=run.json`

## Errors

If you get the following error `2015/08/20 15:01:16 Error getting oauth token: Received a status code 401 Unauthorized. Please check your username and password.`

The user you are connecting with may not have the `doppler.firehose` scope on their client user. 

You can add this as follows
```
uaac target https://uaa.10.244.0.34.xip.io --skip-ssl-validation
uaac token get admin
uaac client update admin --authorities "clients.read password.write clients.secret clients.write uaa.admin scim.write scim.read doppler.firehose"
```