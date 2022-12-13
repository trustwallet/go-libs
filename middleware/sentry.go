package middleware

import (
	"bytes"
	"io"
	"net/http"
	"strconv"

	"github.com/evalphobia/logrus_sentry"
	"github.com/getsentry/raven-go"
	log "github.com/sirupsen/logrus"
)

type SentryOption func(hook *logrus_sentry.SentryHook) error
type SentryCondition func(res *http.Response, url string) bool

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
	hook.StacktraceConfiguration.SwitchExceptionTypeAndMessage = true

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
	// Improve ways to identify if worth logging the error
	if statusCode != http.StatusOK && statusCode != http.StatusNotFound {
		log.WithFields(log.Fields{
			"tags": raven.Tags{
				{Key: "status_code", Value: strconv.Itoa(res.StatusCode)},
				{Key: "host", Value: res.Request.URL.Host},
				{Key: "path", Value: res.Request.URL.Path},
				{Key: "body", Value: getBody(res)},
			},
			"url":         url,
			"fingerprint": []string{"client_errors"},
		}).Error("Client Errors")
	}

	return nil
}

// GetSentryErrorHandler initializes sentry logger for http response errors
// Responses to be logged are defined via passed conditions
func GetSentryErrorHandler(conditions ...SentryCondition) func(res *http.Response, url string) error {
	return func(res *http.Response, url string) error {
		for _, condition := range conditions {
			if condition(res, url) {
				log.WithFields(log.Fields{
					"tags": raven.Tags{
						{Key: "status_code", Value: strconv.Itoa(res.StatusCode)},
						{Key: "host", Value: res.Request.URL.Host},
						{Key: "path", Value: res.Request.URL.Path},
						{Key: "body", Value: getBody(res)},
					},
					"url":         url,
					"fingerprint": []string{"client_errors"},
				}).Error("Client Errors")

				break
			}
		}

		return nil
	}
}

func getBody(res *http.Response) string {
	bodyBytes, _ := io.ReadAll(res.Body)
	_ = res.Body.Close() //  must close
	res.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	return string(bodyBytes)
}

var (
	// SentryConditionAnd returns true only when all conditions are satisfied
	SentryConditionAnd = func(conditions ...SentryCondition) SentryCondition {
		return func(res *http.Response, url string) bool {
			result := true
			for _, condition := range conditions {
				if !condition(res, url) {
					result = false
					break
				}
			}

			return result
		}
	}

	// SentryConditionOr return true when any of conditions is satisfied
	SentryConditionOr = func(conditions ...SentryCondition) SentryCondition {
		return func(res *http.Response, url string) bool {
			for _, condition := range conditions {
				if condition(res, url) {
					return true
				}
			}

			return false
		}
	}

	SentryConditionNotStatusOk = func(res *http.Response, _ string) bool {
		return res.StatusCode < 200 || res.StatusCode > 299
	}

	SentryConditionNotStatusBadRequest = func(res *http.Response, _ string) bool {
		return res.StatusCode != http.StatusBadRequest
	}

	SentryConditionNotStatusNotFound = func(res *http.Response, _ string) bool {
		return res.StatusCode != http.StatusNotFound
	}
)
