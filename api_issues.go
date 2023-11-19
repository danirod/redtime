package main

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
)

type remoteNewIssuePayloadData struct {
	Issue struct {
		Subject       string `json:"subject"`
		ProjectId     uint   `json:"project_id"`
		TrackerId     uint   `json:"tracker_id"`
		ParentIssueId uint   `json:"parent_issue_id"`
	} `json:"issue"`
}

func (payloadData *remoteNewIssuePayloadData) Encode() (string, error) {
	data := make(map[string]interface{})
	data["subject"] = payloadData.Issue.Subject
	data["project_id"] = payloadData.Issue.ProjectId
	data["tracker_id"] = payloadData.Issue.TrackerId
	if pid := payloadData.Issue.ParentIssueId; pid != 0 {
		data["parent_issue_id"] = pid
	}

	wrapped := map[string]interface{}{
		"issue": data,
	}

	encoded, err := json.Marshal(&wrapped)
	if err != nil {
		return "", err
	}
	return string(encoded), nil
}

func (ctx *redmineClient) pushNewIssue(payload remoteNewIssuePayloadData) error {
	data, err := payload.Encode()
	if err != nil {
		return err
	}
	req, err := ctx.buildPostRequest("issues", strings.NewReader(data))
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")
	_, err = ctx.secureCreate(req)
	if err != nil {
		return err
	}
	return nil
}

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
	IssueIds      []int
}

func (params *remoteIssueParams) Encode() url.Values {
	values := url.Values{}
	if params.ProjectId > 0 {
		values.Add("project_id", fmt.Sprintf("%d", params.ProjectId))
	}
	if len(params.IssueIds) > 0 {
		issue_id := ""
		for _, id := range params.IssueIds {
			issue_id = fmt.Sprintf("%s,%d", issue_id, id)
		}
		values.Add("issue_id", issue_id[1:])
	}
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
