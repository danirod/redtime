package main

import (
	"io"
	"net/url"
	"strings"
	"testing"
)

func TestRequestBuilderCreatesGet(t *testing.T) {
	rb, err := NewRequestBuilder("http://localhost:3000/redmine", "api_key_1234")
	if err != nil {
		t.Error(err)
	}
	params := url.Values{
		"project_id": []string{"4"},
	}
	rq, err := rb.Get("issues", params)
	if err != nil {
		t.Error(err)
	}
	if url := rq.URL.String(); url != "http://localhost:3000/redmine/issues.json?project_id=4" {
		t.Errorf("Invalid URL. Was %s", url)
	}
	if header := rq.Header.Get("X-Redmine-Api-Key"); header != "api_key_1234" {
		t.Errorf("Invalid API Key. Was %s", header)
	}
}

func TestRequestBuilderCreatesPost(t *testing.T) {
	rb, err := NewRequestBuilder("http://localhost:3000/redmine", "api_key_1234")
	if err != nil {
		t.Error(err)
	}

	body := strings.NewReader("body-request")
	rq, err := rb.Post("issues", body)
	if err != nil {
		t.Error(err)
	}
	if url := rq.URL.String(); url != "http://localhost:3000/redmine/issues.json" {
		t.Errorf("Invalid URL. Was %s", url)
	}
	if header := rq.Header.Get("X-Redmine-Api-Key"); header != "api_key_1234" {
		t.Errorf("Invalid API Key. Was %s", header)
	}
	r, err := io.ReadAll(rq.Body)
	if err != nil {
		t.Error(err)
	}
	if string(r) != "body-request" {
		t.Errorf("Body was not correct")
	}
}
