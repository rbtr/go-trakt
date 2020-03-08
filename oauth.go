package trakt

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"path"

	"github.com/pkg/errors"
)

const (
	deviceCodePath  = "/oauth/device/code"
	deviceTokenPath = "/oauth/device/token"
)

var _ Authentication = (*Client)(nil)

type Authentication interface {
	DeviceCode(context.Context) (*DeviceCodeResult, error)
	Token(context.Context, string) (*AuthResult, error)
}

type DeviceCodeBody struct {
	ClientID string `json:"client_id"`
}

type DeviceCodeResult struct {
	DeviceCode      string `json:"device_code"`
	UserCode        string `json:"user_code"`
	VerificationURL string `json:"verification_url"`
	ExpiresIn       int    `json:"expires_in"`
	Interval        int    `json:"interval"`
}

func (c *Client) DeviceCode(ctx context.Context) (*DeviceCodeResult, error) {
	postBody, err := json.Marshal(DeviceCodeBody{ClientID: c.ClientID})
	if err != nil {
		return nil, err
	}
	uri := *c.BaseURL
	uri.Path = path.Join(uri.Path, deviceCodePath)
	req, err := http.NewRequest(http.MethodPost, uri.String(), bytes.NewReader(postBody))
	if err != nil {
		return nil, err
	}
	c.SetHeaders(req)
	req = req.WithContext(ctx)
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.Errorf("error getting device code: %d", resp.StatusCode)
	}
	result := &DeviceCodeResult{}
	err = json.NewDecoder(resp.Body).Decode(result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

type AuthBody struct {
	Code         string `json:"code"`
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}

type AuthResult struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
	CreatedAt    int    `json:"created_at"`
}

func (c *Client) Token(ctx context.Context, code string) (*AuthResult, error) {
	postBody, err := json.Marshal(AuthBody{
		Code:         code,
		ClientID:     c.ClientID,
		ClientSecret: c.ClientSecret,
	})
	if err != nil {
		return nil, err
	}
	uri := *c.BaseURL
	uri.Path = path.Join(uri.Path, deviceTokenPath)
	req, err := http.NewRequest(http.MethodPost, uri.String(), bytes.NewReader(postBody))
	if err != nil {
		return nil, err
	}
	c.SetHeaders(req)
	req = req.WithContext(ctx)
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.Errorf("error getting device code: %d", resp.StatusCode)
	}
	result := &AuthResult{}
	err = json.NewDecoder(resp.Body).Decode(result)
	if err != nil {
		return nil, err
	}
	return result, nil
}
