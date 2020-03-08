package trakt

import (
	"fmt"
	"net/http"
	"net/url"
)

const (
	defaultBaseURL             = "https://api.trakt.tv"
	contentTypeApplicationJSON = "application/json"
)

type Client struct {
	BaseURL      *url.URL
	ClientID     string
	ClientSecret string

	accessToken string
	client      *http.Client
	headers     http.Header
}

func NewClient(httpClient *http.Client, clientID, clientSecret string) (*Client, error) {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	baseURL, _ := url.Parse(defaultBaseURL)
	c := &Client{
		BaseURL:      baseURL,
		ClientID:     clientID,
		ClientSecret: clientSecret,
		client:       httpClient,
	}
	return c, nil
}

func (c *Client) SetAuthorization(token string) {
	c.accessToken = token
}

func (c *Client) SetHeaders(req *http.Request) {
	req.Header.Set("Content-Type", contentTypeApplicationJSON)
	req.Header.Set("trakt-api-key", c.ClientID)
	req.Header.Set("trakt-api-version", "2")
	if c.accessToken != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.accessToken))
	}
}
