package middleware

import (
	"github.com/evalphobia/logrus_sentry"
	log "github.com/sirupsen/logrus"
)

func SetupSentry(dsn string) error {
	hook, err := logrus_sentry.NewSentryHook(dsn, []log.Level{
		log.PanicLevel,
		log.FatalLevel,
		log.ErrorLevel,
	})
	if err != nil {
		return err
	}
	hook.StacktraceConfiguration.Enable = true
	hook.StacktraceConfiguration.IncludeErrorBreadcrumb = true
	hook.StacktraceConfiguration.Context = 10
	hook.StacktraceConfiguration.SendExceptionType = true
	log.AddHook(hook)
	return nil
}
