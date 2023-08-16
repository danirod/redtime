package main

import "testing"

func TestAssignNothing(t *testing.T) {
	var p remoteTimelogParams
	if err := p.TimeFilter("", ""); err != nil {
		t.Error(err)
	}
	params := p.Encode()
	if params.Has("from") {
		t.Errorf("from parameter is set, was not expected")
	}
	if params.Has("to") {
		t.Errorf("to parameter is set, was not expected")
	}
}

func TestAssignSinceParameter(t *testing.T) {
	var p remoteTimelogParams
	if err := p.TimeFilter("2023-07-01", ""); err != nil {
		t.Error(err)
	}
	params := p.Encode()
	if got := params.Get("from"); got != "2023-07-01" {
		t.Errorf("from parameter not set, expected 2023-07-01, got %s", got)
	}
	if params.Has("to") {
		t.Errorf("to parameter is set, was not expected")
	}
}

func TestAssignUntilParameter(t *testing.T) {
	var p remoteTimelogParams
	if err := p.TimeFilter("", "2023-07-31"); err != nil {
		t.Error(err)
	}
	params := p.Encode()
	if params.Has("from") {
		t.Errorf("from parameter is set, was not expected")
	}
	if got := params.Get("to"); got != "2023-07-31" {
		t.Errorf("to parameter not set, expected 2023-07-31, got %s", got)
	}
}

func TestAssignSinceAndUntilParameter(t *testing.T) {
	var p remoteTimelogParams
	if err := p.TimeFilter("2023-07-01", "2023-07-31"); err != nil {
		t.Error(err)
	}
	params := p.Encode()
	if got := params.Get("from"); got != "2023-07-01" {
		t.Errorf("from parameter not set, expected 2023-07-01, got %s", got)
	}
	if got := params.Get("to"); got != "2023-07-31" {
		t.Errorf("to parameter not set, expected 2023-07-31, got %s", got)
	}
}

func TestAssignInvalidDate(t *testing.T) {
	cases := []struct {
		from string
		to   string
	}{
		{"invalid", ""},
		{"", "invalid"},
		{"hello", "world"},
	}
	var p remoteTimelogParams
	for _, tt := range cases {
		if p.TimeFilter(tt.from, tt.to) == nil {
			t.Errorf("expected an error after setting %s,%s, did not get one", tt.from, tt.to)
		}
	}
}

func TestAssignUntilOlderThanSince(t *testing.T) {
	var p remoteTimelogParams
	if p.TimeFilter("2023-08-01", "2023-07-01") == nil {
		t.Errorf("expected an error because start date is more recent than end date")
	}
}
