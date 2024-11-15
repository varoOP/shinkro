package plex

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"

	"github.com/pkg/errors"
)

func (c *Client) GetShowID(ctx context.Context, key string) (*GUID, error) {
	baseUrl, err := url.Parse(c.config.Url)
	if err != nil {
		return nil, errors.Wrap(err, "plex url invalid")
	}

	baseUrl = baseUrl.JoinPath(key)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, baseUrl.String(), nil)
	if err != nil {
		return nil, errors.Errorf("%v, request=%v", err, *req)
	}

	var plexResp PlexResponse

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Add("ContainerStart", "X-Plex-Container-Start=0")
	req.Header.Add("ContainerSize", "Plex-Container-Size=100")

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "network error")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Errorf("%v, response status: %v, response body: %v", err, resp.StatusCode, string(body))
	}

	defer resp.Body.Close()
	err = json.Unmarshal(body, &plexResp)
	if err != nil {
		return nil, errors.Errorf("%v, response status: %v, response body: %v", err, resp.StatusCode, string(body))
	}

	if len(plexResp.MediaContainer.Metadata) == 1 {
		return &plexResp.MediaContainer.Metadata[0].GUID, nil
	}

	return nil, errors.Errorf("something went wrong in getting guid from plex:%v, response status: %v, response body: %v", err, resp.StatusCode, string(body))
}

func (c *Client) Test(ctx context.Context) error {
	resp, err := c.http.Get(c.config.Url)
	if err != nil {
		return errors.Wrap(err, "network error")
	}

	if resp.StatusCode == http.StatusUnauthorized {
		return errors.New("unauthorized: check plex token")
	}

	if resp.StatusCode != http.StatusOK {
		return errors.New("woops plex won't let us connect")
	}

	return nil
}
