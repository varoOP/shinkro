package domain

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"runtime"
)

type ShinkroClient struct {
	Client http.Client
}

func (c *ShinkroClient) NewRequestWithContext(ctx context.Context, method, url string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", fmt.Sprintf("shinkro/%v (%v;%v)", runtime.Version(), runtime.GOOS, runtime.GOARCH))

	return req, nil
}

func GetWithContext(ctx context.Context, url string) (*http.Response, error) {
	client := ShinkroClient{
		Client: *http.DefaultClient,
	}

	req, err := client.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	return client.Client.Do(req)
} 