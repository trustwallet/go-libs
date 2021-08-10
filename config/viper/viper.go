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
		default:
			if err := viper.BindEnv(strings.Join(append(parts, tv), ".")); err != nil {
				log.Fatal(err)
			}
		}
	}
}
