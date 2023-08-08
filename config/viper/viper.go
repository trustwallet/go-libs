package viper

import (
	"reflect"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func Load(confPath string, receiver interface{}) {
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	viper.AddConfigPath(".")
	viper.SetConfigName("config")
	viper.SetConfigType("yml")

	configType := "default"
	if confPath != "" {
		viper.SetConfigFile(confPath)
		configType = "supplied"
	}

	err := viper.ReadInConfig()
	if err != nil {
		log.WithError(err).Fatalf("Read %s config", configType)
	}

	log.WithFields(log.Fields{"config": viper.ConfigFileUsed()}).Infof("Viper using %s config", configType)

	bindEnvs(reflect.ValueOf(receiver))
	if err := viper.Unmarshal(receiver); err != nil {
		log.Panic(err, "Error Unmarshal Viper Config File")
	}
}

// viper supports unmarshaling from env vars if the keys are known
//
// bindEnvs is a hack to let viper know in advance what keys exists
func bindEnvs(v reflect.Value, parts ...string) {
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return
		}

		bindEnvs(v.Elem(), parts...)
		return
	}

	ift := v.Type()
	for i := 0; i < ift.NumField(); i++ {
		val := v.Field(i)
		t := ift.Field(i)
		tv, ok := t.Tag.Lookup("mapstructure")
		if !ok {
			continue
		}
		switch val.Kind() {
		case reflect.Struct:
			bindEnvs(val, append(parts, tv)...)
		case reflect.Map:
			// bindEnvs hack doesn't work for maps, because we don't know all the possible
			// values for map keys. Therefore we do nothing to fallback to viper's default key detection.
			continue
		default:
			if err := viper.BindEnv(strings.Join(append(parts, tv), ".")); err != nil {
				log.Fatal(err)
			}
		}
	}
}
