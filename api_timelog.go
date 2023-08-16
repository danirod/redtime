package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"time"
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
	timeSince          time.Time
	timeUntil          time.Time
}

func (params *remoteTimelogParams) TimeFilter(since, until string) (err error) {
	if since != "" {
		if params.timeSince, err = time.Parse(time.DateOnly, since); err != nil {
			return
		}
	}
	if until != "" {
		if params.timeUntil, err = time.Parse(time.DateOnly, until); err != nil {
			return
		}
	}

	unixSince := params.timeSince.Unix()
	unixUntil := params.timeUntil.Unix()
	if unixSince > 0 && unixUntil > 0 && unixSince > unixUntil {
		err = errors.New("until is older than since")
	}
	return
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
	if params.timeSince.Unix() > 0 {
		values.Set("from", params.timeSince.Format("2006-01-02"))
	}
	if params.timeUntil.Unix() > 0 {
		values.Set("to", params.timeUntil.Format("2006-01-02"))
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
