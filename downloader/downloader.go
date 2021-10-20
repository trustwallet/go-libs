package downloader

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	log "github.com/sirupsen/logrus"
)

type Interface interface {
	Download(url string) ([]byte, error)
}

type dl struct {
	client     http.Client
	bytesLimit int64
}

func (d *dl) Download(url string) ([]byte, error) {
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
	if d.bytesLimit > 0 {
		reader = &io.LimitedReader{R: resp.Body, N: d.bytesLimit}
	}

	b, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	return b, nil
}

type Option func(d *dl) error

func New(opts ...Option) (Interface, error) {
	d := &dl{
		client: http.Client{},
	}

	for _, opt := range opts {
		if err := opt(d); err != nil {
			return nil, err
		}
	}

	return d, nil
}

func OptionBytesLimit(bytesLimit int64) Option {
	return func(d *dl) error {
		d.bytesLimit = bytesLimit
		return nil
	}
}

func OptionHttpClient(client http.Client) Option {
	return func(d *dl) error {
		d.client = client
		return nil
	}
}
