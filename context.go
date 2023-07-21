package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type Context struct {
	apiRoot    string
	apiRootURL *url.URL
	apiToken   string
}

func NewContext(root, token string) (*Context, error) {
	apiRootURL, err := url.Parse(root)
	if err != nil {
		return nil, err
	}
	return &Context{apiRoot: root, apiToken: token, apiRootURL: apiRootURL}, nil
}

func (ctx *Context) buildGetRequest(urlpath string, params *url.Values) (*http.Request, error) {
	if !strings.HasSuffix(urlpath, ".json") {
		urlpath = urlpath + ".json"
	}
	finalURL := ctx.apiRootURL.JoinPath(urlpath)
	if params != nil {
		finalURL.RawQuery = params.Encode()
	}
	req, err := http.NewRequest(http.MethodGet, finalURL.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("X-Redmine-Api-Key", ctx.apiToken)
	return req, nil
}

func (ctx *Context) buildPostRequest(urlpath string, body io.Reader) (*http.Request, error) {
	// Build the payload URL.
	if !strings.HasSuffix(urlpath, ".json") {
		urlpath = urlpath + ".json"
	}
	finalURL := ctx.apiRootURL.JoinPath(urlpath)

	req, err := http.NewRequest(http.MethodPost, finalURL.String(), body)
	if err != nil {
		return nil, err
	}
	req.Header.Add("X-Redmine-Api-Key", ctx.apiToken)
	req.Header.Add("Content-Type", "application/json")
	return req, nil
}

func (ctx *Context) secureRequest(req *http.Request) ([]byte, error) {
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("unknown HTTP status code: %d", resp.StatusCode)
	}
	return io.ReadAll(resp.Body)
}

func (ctx *Context) secureCreate(req *http.Request) ([]byte, error) {
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 201 {
		return nil, fmt.Errorf("status code is not create: %d", resp.StatusCode)
	}
	return io.ReadAll(resp.Body)
}
