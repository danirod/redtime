package main

import (
	"io"
	"net/http"
	"net/url"
	"strings"
)

type RequestBuilder struct {
	api  string
	root *url.URL
}

func NewRequestBuilder(root, api string) (*RequestBuilder, error) {
	url_root, err := url.Parse(root)
	if err != nil {
		return nil, err
	}
	return &RequestBuilder{root: url_root, api: api}, nil
}

func (rb *RequestBuilder) Get(uri string, params url.Values) (*http.Request, error) {
	return rb.request(http.MethodGet, uri, params, nil)
}

func (rb *RequestBuilder) Post(uri string, body io.Reader) (*http.Request, error) {
	return rb.request(http.MethodPost, uri, url.Values{}, body)
}

func (rb *RequestBuilder) request(method, uri string, params url.Values, body io.Reader) (*http.Request, error) {
	if !strings.HasSuffix(uri, ".json") {
		uri = uri + ".json"
	}
	path := rb.root.JoinPath(uri)
	path.RawQuery = params.Encode()
	rq, err := http.NewRequest(method, path.String(), body)
	if err == nil {
		rq.Header.Add("X-Redmine-Api-Key", rb.api)
	}
	return rq, err
}
