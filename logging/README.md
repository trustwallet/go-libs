# logging package

Add dependency to the project

```sh
go get github.com/trustwallet/golibs/logging
```

Features:

* [Logrus Wrapper](#logrus-wrapper) allows an easy acces to common logger instance as well as override it in tests
* [Logging Configuration](#logging-configuration) allows for the logging configuration to be loaded with viper


### Logrus Wrapper

By default `logrus` operates agains global instance (common for `go` ), but such approach doesn't allow to replace the `logger` instance during testing.
The `logging` package allows to get the current instance of the `logger`:

```go
log := logging.GetLogger() 
```

Also there is a helper method to get log Entry with `component` filed preset:

```go
log := logging.GetLoggerForComponent("market")
log.Info("some log entry")
// time="2021-08-19T12:33:21Z" level=info msg="some log entry" component="market"
```

For testing purposes the `logger` instance can be replaced:

```go
func TestMyService (t *testing.T) {
        testLogger, hook := test.NewNullLogger()
	testLogger.SetLevel(logrus.WarnLevel)
	logging.SetLogger(testLogger)

	// create instance of service which utilises logging.GetLogger() inside
	s := service.NewService() 
	s.DoSomWork()

	// all logged messages are available here
	for _, e := range hook.Entries {
		t.Log(e)
	}
}
```

### Logging Configuration

Utilizes [Logrus Mate](https://github.com/gogap/logrus_mate) and 
[Logrus Helper](https://github.com/heirko/go-contrib/tree/master/logrusHelper) to load configuration with [Viper](https://github.com/spf13/viper) 🐍  

Which means it can be easily specified via config file per environment (e.g. disable timestamps when deployed to Heroku)

Assuming the `config.yml`

```yaml
market:
  foo: bar

logging:
  level: debug 
  formatter:
    name: text
    options:
      disable_timestamp: true
```

And the corresponding go `struct`:

```go
type Configuration struct {
	Market struct {
		Foo string `mapstructure:"foo"`
	} `mapstructure:"market"`
        Logging logging.Config `mapstructure:"logging"`
}
```

Once viper has unmarshalled the configuration taken from all sources:

```go
err = logging.SetLoggerConfig(config.Logging)
// error handling

log := logging.GetLogger()
```

✨  It's fully backward compatible for code which uses `logrus` directly.

```go

import log "github.com/sirupsen/logrus"

func LogSomething() {
         log.Info("some log message") // will respect logging configuration
}
```
