package middleware

import (
	"net/http"
	"strconv"

	"github.com/evalphobia/logrus_sentry"
	"github.com/getsentry/raven-go"
	log "github.com/sirupsen/logrus"
)

type SentryOption func(hook *logrus_sentry.SentryHook) error

func SetupSentry(dsn string, opts ...SentryOption) error {
	hook, err := logrus_sentry.NewSentryHook(dsn, []log.Level{
		log.PanicLevel,
		log.FatalLevel,
		log.ErrorLevel,
	})
	if err != nil {
		return err
	}
	hook.Timeout = 0
	hook.StacktraceConfiguration.Enable = true
	hook.StacktraceConfiguration.IncludeErrorBreadcrumb = true
	hook.StacktraceConfiguration.Context = 10
	hook.StacktraceConfiguration.SendExceptionType = true

	for _, o := range opts {
		err = o(hook)
		if err != nil {
			return err
		}
	}

	log.AddHook(hook)
	return nil
}

func WithDefaultLoggerName(name string) SentryOption {
	return func(hook *logrus_sentry.SentryHook) error {
		hook.SetDefaultLoggerName(name)
		return nil
	}
}

func WithEnvironment(env string) SentryOption {
	return func(hook *logrus_sentry.SentryHook) error {
		hook.SetEnvironment(env)
		return nil
	}
}

func WithHttpContext(h *raven.Http) SentryOption {
	return func(hook *logrus_sentry.SentryHook) error {
		hook.SetHttpContext(h)
		return nil
	}
}

func WithIgnoreErrors(errs ...string) SentryOption {
	return func(hook *logrus_sentry.SentryHook) error {
		return hook.SetIgnoreErrors(errs...)
	}
}

func WithIncludePaths(p []string) SentryOption {
	return func(hook *logrus_sentry.SentryHook) error {
		hook.SetIncludePaths(p)
		return nil
	}
}

func WithRelease(release string) SentryOption {
	return func(hook *logrus_sentry.SentryHook) error {
		hook.SetRelease(release)
		return nil
	}
}

func WithSampleRate(rate float32) SentryOption {
	return func(hook *logrus_sentry.SentryHook) error {
		return hook.SetSampleRate(rate)
	}
}

func WithTagsContext(t map[string]string) SentryOption {
	return func(hook *logrus_sentry.SentryHook) error {
		hook.SetTagsContext(t)
		return nil
	}
}

func WithUserContext(u *raven.User) SentryOption {
	return func(hook *logrus_sentry.SentryHook) error {
		hook.SetUserContext(u)
		return nil
	}
}

func WithServerName(serverName string) SentryOption {
	return func(hook *logrus_sentry.SentryHook) error {
		hook.SetServerName(serverName)
		return nil
	}
}

var SentryErrorHandler = func(res *http.Response, url string) error {
	statusCode := res.StatusCode
	//Improve ways to identify if worth logging the error
	if statusCode != http.StatusOK && statusCode != http.StatusNotFound {
		log.WithFields(log.Fields{
			"tags": raven.Tags{
				{Key: "status_code", Value: strconv.Itoa(res.StatusCode)},
				{Key: "host", Value: res.Request.URL.Host},
				{Key: "path", Value: res.Request.URL.Path},
			},
			"url":         url,
			"fingerprint": []string{"client_errors"},
		}).Error("Client Errors")
	}

	return nil
}
