package client

import (
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"ping-pong/internal/config"
)

type HTTPClient struct {
	client *http.Client
}

func NewHTTPClient(cfg *config.Config) *HTTPClient {
	return &HTTPClient{
		client: &http.Client{
			Timeout: cfg.ClientConfig.Timeout,
		},
	}
}

func (c *HTTPClient) GetString(url string) string {
	resp, err := c.client.Get(url)
	if err != nil {
		return ""
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			logrus.Error(err.Error())
		}
	}(resp.Body)

	if resp.StatusCode == http.StatusOK {
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			logrus.Errorf("error occured while reading response: %s", err.Error())
		}

		return string(bodyBytes)
	}
	return ""
}
