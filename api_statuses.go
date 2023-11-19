package main

import "encoding/json"

type remoteIssueStatusEntry struct {
	namedEntry
}

type remoteIssueStatusList struct {
	Statuses []remoteIssueStatusEntry `json:"issue_statuses"`
}

func (ctx *redmineClient) fetchStatuses() (*remoteIssueStatusList, error) {
	req, err := ctx.buildGetRequest("issue_statuses", nil)
	if err != nil {
		return nil, err
	}
	data, err := ctx.secureRequest(req)
	if err != nil {
		return nil, err
	}

	var list remoteIssueStatusList
	if err := json.Unmarshal(data, &list); err != nil {
		return nil, err
	}
	return &list, nil
}
