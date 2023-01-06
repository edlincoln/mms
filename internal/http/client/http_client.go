package client

import (
	"context"
	"net/http"
	"net/url"
)

//go:generate mockery --name=HttpClient --output=../../client --filename=http_client_mock.go --structname=HttpClientMock --inpackage=true
type HttpClient interface {
	Get(ctx context.Context, destination string, headers map[string]string, params map[string]string) (*http.Response, error)
}

type HttpClientImpl struct {
	httpClient *http.Client
}

func NewHttpClient(httpClient *http.Client) HttpClient {
	return &HttpClientImpl{httpClient: httpClient}
}

func (h *HttpClientImpl) Get(ctx context.Context, destination string, headers map[string]string, params map[string]string) (*http.Response, error) {
	u, err := url.Parse(destination)
	if err != nil {
		return nil, err
	}

	if params != nil {
		values := url.Values{}
		for k, v := range params {
			values.Add(k, v)
		}
		u.RawQuery = values.Encode()
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	return h.httpClient.Do(req)
}
