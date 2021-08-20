package logging

import (
	mate "github.com/heralight/logrus_mate"
	"github.com/sirupsen/logrus"
)

type TextFormatterConfig struct {
	ForceColors      bool   `json:"force_colors,string"`
	DisableColors    bool   `json:"disable_colors,string"`
	DisableTimestamp bool   `json:"disable_timestamp,string"`
	FullTimestamp    bool   `json:"full_timestamp,string"`
	TimestampFormat  string `json:"timestamp_format"`
	DisableSorting   bool   `json:"disable_sorting,string"`
}

func init() {
	mate.RegisterFormatter("strict_text", NewTextFormatter)
}

func NewTextFormatter(options mate.Options) (formatter logrus.Formatter, err error) {
	conf := TextFormatterConfig{}

	if err = options.ToObject(&conf); err != nil {
		return
	}

	formatter = &logrus.TextFormatter{
		ForceColors:      conf.ForceColors,
		DisableColors:    conf.DisableColors,
		DisableTimestamp: conf.DisableTimestamp,
		FullTimestamp:    conf.FullTimestamp,
		TimestampFormat:  conf.TimestampFormat,
		DisableSorting:   conf.DisableSorting,
	}
	return
}
