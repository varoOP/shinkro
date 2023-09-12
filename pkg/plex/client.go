package plex

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"runtime"

	"github.com/pkg/errors"
)

func NewPlexClient(url string, token string) *PlexClient {
	return &PlexClient{
		Url:   url,
		Token: token,
	}
}

func (p *PlexClient) GetShowID(key string) (*GUID, error) {
	baseUrl, err := url.Parse(p.Url)
	if err != nil {
		return nil, errors.Wrap(err, "plex url invalid")
	}

	baseUrl = baseUrl.JoinPath(key)
	params := url.Values{}
	params.Add("X-Plex-Token", p.Token)
	baseUrl.RawQuery = params.Encode()
	req, err := http.NewRequest(http.MethodGet, baseUrl.String(), nil)
	if err != nil {
		return nil, errors.Errorf("%v, request=%v", err, *req)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Add("ContainerStart", "X-Plex-Container-Start=0")
	req.Header.Add("ContainerSize", "Plex-Container-Size=100")
	req.Header.Set("User-Agent", fmt.Sprintf("shinkro/%v (%v;%v)", runtime.Version(), runtime.GOOS, runtime.GOARCH))

	resp, err := p.Client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "network error")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Errorf("%v, response status: %v, response body: %v", err, resp.StatusCode, string(body))
	}

	defer resp.Body.Close()
	err = json.Unmarshal(body, &p.Resp)
	if err != nil {
		return nil, errors.Errorf("%v, response status: %v, response body: %v", err, resp.StatusCode, string(body))
	}

	if len(p.Resp.MediaContainer.Metadata) == 1 {
		return &p.Resp.MediaContainer.Metadata[0].GUID, nil
	}

	return nil, errors.Errorf("something went wrong in getting guid from plex:%v, response status: %v, response body: %v", err, resp.StatusCode, string(body))
}
