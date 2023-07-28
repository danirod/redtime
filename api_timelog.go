package main

import (
	"encoding/json"
	"fmt"
	"net/url"
)

type remoteTimelogEntry struct {
	Id       int        `json:"id"`
	Project  namedEntry `json:"project"`
	Issue    namedEntry `json:"issue"`
	User     namedEntry `json:"user"`
	Activity namedEntry `json:"activity"`
	Hours    float64    `json:"hours"`
	Comment  string     `json:"comments"`
	SpentOn  string     `json:"spent_on"`
}

type remoteTimelogList struct {
	Entries    []remoteTimelogEntry `json:"time_entries"`
	TotalCount int                  `json:"total_count"`
	Offset     int                  `json:"offset"`
	Limit      int                  `json:"limit"`
}

type remoteTimelogParams struct {
	criteria           string // "project" or "issue"
	projectId          uint
	excludeSubprojects bool
	issueId            uint
}

func (params *remoteTimelogParams) Encode() url.Values {
	values := url.Values{}
	switch params.criteria {
	case "project":
		values.Add("project_id", fmt.Sprintf("%d", params.projectId))
		if params.excludeSubprojects {
			values.Add("subproject_id", "!*")
		}
	case "issue":
		values.Add("issue_id", fmt.Sprintf("%d", params.issueId))
	}
	return values
}

func (ctx *redmineClient) fetchTimeEntries(param *remoteTimelogParams) (*remoteTimelogList, error) {
	querystring := url.Values{}
	if param != nil {
		querystring = param.Encode()
	}

	req, err := ctx.buildGetRequest("time_entries", &querystring)
	if err != nil {
		return nil, err
	}

	data, err := ctx.secureRequest(req)
	if err != nil {
		return nil, err
	}

	var page remoteTimelogList
	if err := json.Unmarshal(data, &page); err != nil {
		return nil, err
	}
	if err := ctx.enhanceTimelogEntries(&page); err != nil {
		return nil, err
	}
	return &page, nil
}

func (ctx *redmineClient) enhanceTimelogEntries(list *remoteTimelogList) error {
	issueIds := []int{}
	issueSet := make(map[int]string)
	for _, entry := range list.Entries {
		if _, ok := issueSet[entry.Issue.Id]; !ok {
			issueSet[entry.Issue.Id] = ""
			issueIds = append(issueIds, entry.Issue.Id)
		}
	}

	params := remoteIssueParams{IssueIds: issueIds, IncludeClosed: true}
	issues, err := ctx.fetchIssues(&params)
	if err != nil {
		return err
	}
	for _, issue := range issues.Issues {
		issueSet[issue.Id] = issue.Subject
	}
	for idx, entry := range list.Entries {
		if val, ok := issueSet[entry.Issue.Id]; ok {
			list.Entries[idx].Issue.Name = val
		}
	}
	return nil
}
