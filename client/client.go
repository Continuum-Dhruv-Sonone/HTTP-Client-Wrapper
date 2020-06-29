package client

import (
	"io"
	"net/http"
	"net/url"

	"github.com/afex/hystrix-go/hystrix"
	"github.com/pkg/errors"
)

var (
	errCb = errors.New("Server returned error")
)

type client struct {
	client *http.Client
}

// New creates a new instance of httpClient
func New() *client {
	return &client{
		client: &http.Client{},
	}
}

func (cl *client) Do(req *http.Request) (*http.Response, error) {

	var (
		response *http.Response
		err      error
	)

	commandName := req.URL.Host

	err = hystrix.Do(commandName, func() error {
		response, err = cl.client.Do(req)
		if err != nil {
			return err
		}

		if response.StatusCode >= http.StatusInternalServerError {
			return errCb
		}
		return nil
	}, nil)

	if err == errCb {
		return response, nil
	}

	return response, err
}

func (cl *client) Get(reqURL string) (*http.Response, error) {

	var (
		response *http.Response
		err      error
	)

	u, err := url.Parse(reqURL)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to parse url: %s", reqURL)
	}

	hostname := u.Host

	err = hystrix.Do(hostname, func() error {
		response, err = cl.client.Get(reqURL)
		if err != nil {
			return err
		}

		if response.StatusCode >= http.StatusInternalServerError {
			return errCb
		}
		return nil
	}, nil)

	if err == errCb {
		return response, nil
	}

	return response, err
}

func (cl *client) Post(reqURL, contentType string, body io.Reader) (*http.Response, error) {

	var (
		response *http.Response
		err      error
	)

	u, err := url.Parse(reqURL)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to parse url: %s", reqURL)
	}

	hostname := u.Host

	err = hystrix.Do(hostname, func() error {
		response, err = cl.client.Post(reqURL, contentType, body)
		if err != nil {
			return err
		}

		if response.StatusCode >= http.StatusInternalServerError {
			return errCb
		}
		return nil
	}, nil)

	if err == errCb {
		return response, nil
	}

	return response, err
}
