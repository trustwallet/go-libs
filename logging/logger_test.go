package logging_test

import (
	"bytes"
	"strconv"
	"strings"
	"testing"

	"github.com/heirko/go-contrib/logrusHelper"
	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/spf13/viper"
	"github.com/trustwallet/go-libs/logging"
	"gotest.tools/assert"
)

func TestGetLogger(t *testing.T) {
	logger := logging.GetLogger()

	assert.Equal(t, logger.Level, logrus.InfoLevel, "default logger minimum level is Info")
}

func TestGetLoggerForComponent(t *testing.T) {
	logger1 := logging.GetLoggerForComponent("logger1")
	logger1 = logger1.WithField("custom", "logger1_only")

	logger2 := logging.GetLoggerForComponent("logger2")

	logAndAssertText(t, logger1, func(fields map[string]string) {
		assert.Equal(t, "logger1", fields[logging.FieldKeyComponent])
	})
	logAndAssertText(t, logger2, func(fields map[string]string) {
		assert.Equal(t, "logger2", fields[logging.FieldKeyComponent])

		_, ok := fields["custom"]
		assert.Assert(t, !ok, "custom field should exist on logger1 only")
	})
}

func TestParseConfigWithViper(t *testing.T) {
	yamlConfig := []byte(`
logging:
  out:
    name: stdout
  level: debug 
  formatter:
    name: text
    options:
      disable_colors: true
      full_timestamp: false
    hooks:
    - name: file
      options:
        filename: debug.log,
        maxsize: 5000,
        maxdays: 1,
        rotate: true,
        priority: LOG_INFO,
        tag: ""
`)

	viper.SetConfigType("yaml")
	err := viper.ReadConfig(bytes.NewBuffer(yamlConfig))
	assert.NilError(t, err)
	t.Logf("All keys: %#v", viper.AllSettings())

	logger := logging.GetLogger()
	// Unmarshal configuration from Viper
	var c = logrusHelper.UnmarshalConfiguration(viper.Sub("logging"))
	err = logrusHelper.SetConfig(logger, c)
	assert.NilError(t, err)

	assert.Equal(t, logger.Level, logrus.DebugLevel, "logging level set to debug via config")
}

func TestSetLoggerConfig(t *testing.T) {
	yamlConfig := []byte(`
logging:
  level: debug 
  formatter:
    name: text
    options:
      disable_timestamp: true
`)

	viper.SetConfigType("yaml")
	err := viper.ReadConfig(bytes.NewBuffer(yamlConfig))
	assert.NilError(t, err)

	var config logging.Config
	err = viper.UnmarshalKey("logging", &config)
	assert.NilError(t, err)

	err = logging.SetLoggerConfig(config)
	assert.NilError(t, err)

	logger := logging.GetLogger()
	assert.Equal(t, logger.Level, logrus.DebugLevel, "logging level set to debug via config")
	assert.Equal(t, logger.Formatter.(*logrus.TextFormatter).DisableTimestamp, true)
}

func TestOverrideBoolOptionAsString(t *testing.T) {
	yamlConfig := []byte(`
logging:
  level: debug 
  formatter:
    name: strict_text
    options:
      disable_timestamp: "true"
`)

	viper.SetConfigType("yaml")
	err := viper.ReadConfig(bytes.NewBuffer(yamlConfig))
	assert.NilError(t, err)

	var config logging.Config
	err = viper.UnmarshalKey("logging", &config)
	assert.NilError(t, err)

	err = logging.SetLoggerConfig(config)
	assert.NilError(t, err)

	logger := logging.GetLogger()
	assert.Equal(t, logger.Level, logrus.DebugLevel, "logging level set to debug via config")
	assert.Equal(t, logger.Formatter.(*logrus.TextFormatter).DisableTimestamp, true)
}

func TestSetLoggerConfigForStandardLogger(t *testing.T) {
	// Not every component would be able to use logging.GetLogger()
	// This test makes sure the config loaded with viper is also
	// applied globally to standard logger

	yamlConfig := []byte(`
logging:
  level: debug 
  formatter:
    name: text
    options:
      disable_timestamp: true
`)

	viper.SetConfigType("yaml")
	err := viper.ReadConfig(bytes.NewBuffer(yamlConfig))
	assert.NilError(t, err)

	var config logging.Config
	err = viper.UnmarshalKey("logging", &config)
	assert.NilError(t, err)

	err = logging.SetLoggerConfig(config)
	assert.NilError(t, err)

	std := logrus.StandardLogger()
	assert.Equal(t, std.Formatter.(*logrus.TextFormatter).DisableTimestamp, true)
}

func TestSetLogger(t *testing.T) {
	testLogger, hook := test.NewNullLogger()
	testLogger.SetLevel(logrus.WarnLevel)
	logging.SetLogger(testLogger)

	logger1 := logging.GetLogger()
	logger1.Info("you should not see me printed")
	logger1.Warn("you should see this printed")

	logger2 := logging.GetLoggerForComponent("testing")
	logger2.Debug("you should not see me too")
	logger1.Error("you should see this printed")

	for _, e := range hook.Entries {
		t.Log(e)
	}
	assert.Equal(t, len(hook.Entries), 2)
}

func logAndAssertText(t *testing.T, entry *logrus.Entry, assertions func(fields map[string]string)) {
	var buffer bytes.Buffer
	entry.Logger.Out = &buffer
	entry.Logger.Formatter.(*logrus.TextFormatter).DisableColors = true
	entry.Info()

	fields := make(map[string]string)
	for _, kv := range strings.Split(strings.TrimRight(buffer.String(), "\n"), " ") {
		if !strings.Contains(kv, "=") {
			continue
		}
		kvArr := strings.Split(kv, "=")
		key := strings.TrimSpace(kvArr[0])
		val := kvArr[1]
		if kvArr[1][0] == '"' {
			var err error
			val, err = strconv.Unquote(val)
			assert.NilError(t, err)
		}
		fields[key] = val
	}
	assertions(fields)
}
