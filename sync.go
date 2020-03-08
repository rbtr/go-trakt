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
	syncBasePath   = "/sync"
	collectionPath = "/collection"
)

var _ Sync = (*Client)(nil)

type Sync interface {
	Collection(context.Context, *CollectionBody) (*CollectionResult, error)
}

type CollectionBody struct {
	Movies   []Movie   `json:"movies"`
	Shows    []Show    `json:"shows"`
	Seasons  []Season  `json:"season"`
	Episodes []Episode `json:"episode"`
}

type Collection struct {
	Movies   int `json:"movies"`
	Episodes int `json:"episodes"`
}

type CollectionResult struct {
	Added    Collection `json:"added"`
	Updated  Collection `json:"updated"`
	Existing Collection `json:"existing"`
	NotFound struct {
		Movies   []Movie   `json:"movies"`
		Shows    []Show    `json:"shows"`
		Seasons  []Season  `json:"seasons"`
		Episodes []Episode `json:"episode"`
	} `json:"not_found"`
}

func (c *Client) Collection(ctx context.Context, collectionBody *CollectionBody) (*CollectionResult, error) {
	postBody, err := json.Marshal(collectionBody)
	if err != nil {
		return nil, err
	}
	uri := path.Join(c.BaseURL.String(), syncBasePath, collectionPath)
	req, err := http.NewRequest(http.MethodPost, uri, bytes.NewReader(postBody))
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
		return nil, errors.Errorf("error updating collection: %d", resp.StatusCode)
	}
	result := &CollectionResult{}
	err = json.NewDecoder(resp.Body).Decode(result)
	if err != nil {
		return nil, err
	}
	return result, nil
}
