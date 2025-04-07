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
		return nil, errors.Wrap(err, "plex request invalid")
	}

	var plexResp PlexResponse
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
	baseUrl, err := url.Parse(c.config.Url)
	if err != nil {
		return errors.Wrap(err, "plex url invalid")
	}

	baseUrl = baseUrl.JoinPath("/myplex/account")
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, baseUrl.String(), nil)
	if err != nil {
		return errors.Wrap(err, "plex request invalid")
	}

	resp, err := c.http.Do(req)
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

func (c *Client) GetServerList(ctx context.Context) (*ServerResponse, error) {
	plexUrl := "https://plex.tv/api/v2/resources"
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, plexUrl, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "error reading response body")
	}

	var servers []Server
	if err := json.Unmarshal(body, &servers); err == nil && len(servers) > 0 {
		return &ServerResponse{Servers: servers}, nil
	}

	if resp.StatusCode == http.StatusUnauthorized {
		return nil, errors.New("unauthorized: check plex token")
	}

	if resp.StatusCode == http.StatusBadRequest {
		return nil, errors.New("bad request: check plex client identifier")
	}

	return nil, errors.New("unknown or invalid Plex response")
}

func (c *Client) GetLibraries(ctx context.Context, plexUrl string) (*LibraryResponse, error) {
	libUrl, err := url.Parse(plexUrl)
	if err != nil {
		return nil, errors.Wrap(err, "plex url invalid")
	}

	libUrl = libUrl.JoinPath("/library/sections")

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, libUrl.String(), nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var libResp LibraryResponse
	if err := json.NewDecoder(resp.Body).Decode(&libResp); err != nil {
		return nil, errors.Wrap(err, "error decoding library response body")
	}

	if len(libResp.MediaContainer.Directory) > 0 {
		return &libResp, nil
	}

	if resp.StatusCode == http.StatusUnauthorized {
		return nil, errors.New("unauthorized: check plex token")
	}

	if resp.StatusCode == http.StatusBadRequest {
		return nil, errors.New("bad request: check plex client identifier")
	}

	return nil, errors.New("unknown or invalid Plex response")
}
