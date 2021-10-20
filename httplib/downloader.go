package httplib

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	log "github.com/sirupsen/logrus"
)

type Downloader interface {
	Download(url string) ([]byte, error)
}

type downloader struct {
	client         http.Client
	bytesSizeLimit int64
}

func (d *downloader) Download(url string) ([]byte, error) {
	resp, err := d.client.Get(url)
	if err != nil {
		return nil, err
	}
	if !(resp.StatusCode >= 200 && resp.StatusCode <= 299) {
		return nil, fmt.Errorf("response status code: %d", resp.StatusCode)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.WithField("error", err).Error("cannot close request body")
		}
	}()

	reader, _ := resp.Body.(io.Reader)
	if d.bytesSizeLimit > 0 {
		reader = &io.LimitedReader{R: resp.Body, N: d.bytesSizeLimit}
	}

	b, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	return b, nil
}

type Option func(d *downloader) error

func NewDownloader(opts ...Option) (Downloader, error) {
	d := &downloader{
		client: http.Client{},
	}

	for _, opt := range opts {
		if err := opt(d); err != nil {
			return nil, err
		}
	}

	return d, nil
}

// OptionBytesSizeLimit limits the downloaded file size to the provided number of bytes.
func OptionBytesSizeLimit(n int64) Option {
	return func(d *downloader) error {
		d.bytesSizeLimit = n
		return nil
	}
}

// OptionHttpClient sets a custom http client to perform the request.
func OptionHttpClient(client http.Client) Option {
	return func(d *downloader) error {
		d.client = client
		return nil
	}
}
