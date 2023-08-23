package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type redmineClient struct {
	builder *RequestBuilder
}

func newContext(root, token string) (*redmineClient, error) {
	builder, err := NewRequestBuilder(root, token)
	if err != nil {
		return nil, err
	}
	return &redmineClient{builder: builder}, nil
}

func (ctx *redmineClient) buildGetRequest(urlpath string, params *url.Values) (*http.Request, error) {
	if params != nil {
		return ctx.builder.Get(urlpath, *params)
	}
	return ctx.builder.Get(urlpath, url.Values{})
}

func (ctx *redmineClient) buildPostRequest(urlpath string, body io.Reader) (*http.Request, error) {
	return ctx.builder.Post(urlpath, body)
}

func (ctx *redmineClient) secureRequest(req *http.Request) ([]byte, error) {
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("unknown HTTP status code: %d", resp.StatusCode)
	}
	return io.ReadAll(resp.Body)
}

func (ctx *redmineClient) secureCreate(req *http.Request) ([]byte, error) {
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 201 {
		return nil, fmt.Errorf("status code is not create: %d", resp.StatusCode)
	}
	return io.ReadAll(resp.Body)
}
