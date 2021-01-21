package middleware

import (
	"net/http"
	"strconv"

	"github.com/evalphobia/logrus_sentry"
	"github.com/getsentry/raven-go"
	log "github.com/sirupsen/logrus"
	"github.com/trustwallet/golibs/client"
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

var SentryErrorHandler = func(res *http.Response, uri string) error {
	statusCode := res.StatusCode
	//Improve ways to identify if worth logging the error
	if statusCode != http.StatusOK && statusCode != http.StatusNotFound {
		log.WithFields(log.Fields{
			"tags": raven.Tags{
				{Key: "status_code", Value: strconv.Itoa(res.StatusCode)},
				{Key: "host", Value: res.Request.URL.Host},
				{Key: "url", Value: uri},
			},
			"fingerprint": []string{"client_errors"},
		}).Error("Client Errors")
	}

	return nil
}

func InitClientWithSentry(url string) client.Request {
	return client.Request{
		Headers:      map[string]string{},
		HttpClient:   client.DefaultClient,
		ErrorHandler: SentryErrorHandler,
		BaseUrl:      url,
	}
}

func InitJSONClientWithSentry(baseUrl string) client.Request {
	headers := map[string]string{
		"Content-Type": "application/json",
		"Accept":       "application/json",
	}
	return client.Request{
		Headers:      headers,
		HttpClient:   client.DefaultClient,
		ErrorHandler: SentryErrorHandler,
		BaseUrl:      baseUrl,
	}
}
