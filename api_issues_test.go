package main

import "testing"

func TestEncodeRemoteIssueParamsWithProjectId(t *testing.T) {
	params := remoteIssueParams{
		ProjectId: 5,
	}
	query := params.Encode()
	if val := query.Get("project_id"); val != "5" {
		t.Errorf("Expected project_id to be 5, was %s", val)
	}
}

func TestEncodeRemoteIssueParamsWithoutIncludeClosed(t *testing.T) {
	params := remoteIssueParams{
		IncludeClosed: false,
	}
	query := params.Encode()
	if has := query.Has("status_id"); has {
		t.Errorf("expected status_id not to be present")
	}
}

func TestEncodeRemoteIssueParamsWithIncludeClosed(t *testing.T) {
	params := remoteIssueParams{
		IncludeClosed: true,
	}
	query := params.Encode()
	if val := query.Get("status_id"); val != "*" {
		t.Errorf("expected status_id to be *, was %s", val)
	}
}

func TestEncodeRemoteIssueParamsWithissueIds(t *testing.T) {
	params := remoteIssueParams{
		IssueIds: []int{2, 4, 6, 8},
	}
	query := params.Encode()
	if val := query.Get("issue_id"); val != "2,4,6,8" {
		t.Errorf("expected status_id to be `2,4,6,8`, was %s", val)
	}
}
