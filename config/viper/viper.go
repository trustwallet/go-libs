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

	bindEnvs(receiver)
	if err := viper.Unmarshal(receiver); err != nil {
		log.Panic(err, "Error Unmarshal Viper Config File")
	}
}

func bindEnvs(iface interface{}, parts ...string) {
	ifv := reflect.ValueOf(iface)
	ift := reflect.TypeOf(iface)
	for i := 0; i < ift.NumField(); i++ {
		v := ifv.Field(i)
		t := ift.Field(i)
		tv, ok := t.Tag.Lookup("mapstructure")
		if !ok {
			continue
		}
		switch v.Kind() {
		case reflect.Struct:
			bindEnvs(v.Interface(), append(parts, tv)...)
		default:
			if err := viper.BindEnv(strings.Join(append(parts, tv), ".")); err != nil {
				log.Fatal(err)
			}
		}
	}
}
