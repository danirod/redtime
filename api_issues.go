package main

import (
	"encoding/json"
	"fmt"
	"net/url"
)

type remoteIssue struct {
	Id             int        `json:"id"`
	Subject        string     `json:"subject"`
	Description    string     `json:"description"`
	StartDate      string     `json:"start_date"`
	EndDate        string     `json:"end_date"`
	DoneRatio      int        `json:"done_ratio"`
	EstimatedHours float64    `json:"estiamted_hours"`
	SpentHours     float64    `json:"spent_hours"`
	CreatedOn      string     `json:"created_on"`
	UpdatedOn      string     `json:"updated_on"`
	ClosedOn       string     `json:"closed_on"`
	Project        namedEntry `json:"project"`
	Tracker        namedEntry `json:"tracker"`
	Status         namedEntry `json:"status"`
	Priority       namedEntry `json:"priority"`
	Author         namedEntry `json:"author"`
	AssignedTo     namedEntry `json:"assigned_to"`
	FixedVersion   namedEntry `json:"fixed_version"`
}

type remoteIssueData struct {
	Issue remoteIssue `json:"issue"`
}

type remoteIssueList struct {
	Issues     []remoteIssue `json:"issues"`
	TotalCount int           `json:"total_count"`
	Offset     int           `json:"offset"`
	Limit      int           `json:"limit"`
}

type remoteIssueParams struct {
	ProjectId     int
	IncludeClosed bool
}

func (params *remoteIssueParams) Encode() url.Values {
	values := url.Values{}
	values.Add("project_id", fmt.Sprintf("%d", params.ProjectId))
	if params.IncludeClosed {
		values.Add("status_id", "*")
	}
	return values
}

func (ctx *redmineClient) fetchIssues(params *remoteIssueParams) (*remoteIssueList, error) {
	values := url.Values{}
	if params != nil {
		values = params.Encode()
	}

	req, err := ctx.buildGetRequest("issues", &values)
	if err != nil {
		return nil, err
	}
	data, err := ctx.secureRequest(req)
	if err != nil {
		return nil, err
	}

	var issues remoteIssueList
	if err := json.Unmarshal(data, &issues); err != nil {
		return nil, err
	}
	return &issues, nil
}

func (ctx *redmineClient) fetchIssue(id int) (*remoteIssue, error) {
	url := fmt.Sprintf("issues/%d", id)
	req, err := ctx.buildGetRequest(url, nil)
	if err != nil {
		return nil, err
	}
	data, err := ctx.secureRequest(req)
	if err != nil {
		return nil, err
	}

	var issue remoteIssueData
	if err := json.Unmarshal(data, &issue); err != nil {
		return nil, err
	}
	return &issue.Issue, nil
}
