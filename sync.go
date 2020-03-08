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
	Movies   []Movie   `json:"movies,omitempty"`
	Shows    []Show    `json:"shows,omitempty"`
	Seasons  []Season  `json:"season,omitempty"`
	Episodes []Episode `json:"episodes,omitempty"`
}

type Collection struct {
	Movies   int `json:"movies,omitempty"`
	Episodes int `json:"episodes,omitempty"`
}

type CollectionResult struct {
	Added    Collection `json:"added,omitempty"`
	Updated  Collection `json:"updated,omitempty"`
	Existing Collection `json:"existing,omitempty"`
	NotFound struct {
		Movies   []Movie   `json:"movies,omitempty"`
		Shows    []Show    `json:"shows,omitempty"`
		Seasons  []Season  `json:"seasons,omitempty"`
		Episodes []Episode `json:"episode,omitempty"`
	} `json:"not_found,omitempty"`
}

func (c *Client) Collection(ctx context.Context, collectionBody *CollectionBody) (*CollectionResult, error) {
	postBody, err := json.Marshal(collectionBody)
	if err != nil {
		return nil, err
	}
	uri := *c.BaseURL
	uri.Path = path.Join(uri.Path, syncBasePath, collectionPath)
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
	if resp.StatusCode != 201 {
		return nil, errors.Errorf("error updating collection: %d", resp.StatusCode)
	}
	result := &CollectionResult{}
	err = json.NewDecoder(resp.Body).Decode(result)
	if err != nil {
		return nil, err
	}
	return result, nil
}
