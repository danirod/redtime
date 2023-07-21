package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type NamedRelation struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type TimeEntryActivities struct {
	Activities []NamedRelation `json:"time_entry_activities"`
}

func (act *TimeEntryActivities) Map() map[int]string {
	activities := map[int]string{}
	for _, item := range act.Activities {
		activities[item.Id] = item.Name
	}
	return activities
}

type RemoteIssue struct {
	Id             int           `json:"id"`
	Subject        string        `json:"subject"`
	Description    string        `json:"description"`
	StartDate      string        `json:"start_date"`
	EndDate        string        `json:"end_date"`
	DoneRatio      int           `json:"done_ratio"`
	EstimatedHours float64       `json:"estiamted_hours"`
	SpentHours     float64       `json:"spent_hours"`
	CreatedOn      string        `json:"created_on"`
	UpdatedOn      string        `json:"updated_on"`
	ClosedOn       string        `json:"closed_on"`
	Project        NamedRelation `json:"project"`
	Tracker        NamedRelation `json:"tracker"`
	Status         NamedRelation `json:"status"`
	Priority       NamedRelation `json:"priority"`
	Author         NamedRelation `json:"author"`
	FixedVersion   NamedRelation `json:"fixed_version"`
}

type RemoteIssueData struct {
	Issue RemoteIssue `json:"issue"`
}

type RemoteIssueList struct {
	Issues     []RemoteIssue `json:"issues"`
	TotalCount int           `json:"total_count"`
	Offset     int           `json:"offset"`
	Limit      int           `json:"limit"`
}

type RemoteProject struct {
	Id             int    `json:"id"`
	Name           string `json:"name"`
	Identifier     string `json:"identifier"`
	Description    string `json:"description"`
	Status         int    `json:"status"`
	IsPublic       bool   `json:"is_public"`
	InheritMembers bool   `json:"inherit_members"`
	CreatedOn      string `json:"created_on"`
	UpdatedOn      string `json:"updated_on"`
}

type RemoteProjectData struct {
	Project RemoteProject `json:"project"`
}

type RemoteProjectList struct {
	Projects   []RemoteProject `json:"projects"`
	TotalCount int             `json:"total_count"`
	Offset     int             `json:"offset"`
	Limit      int             `json:"limit"`
}

func (ctx *Context) FetchProjects() (*RemoteProjectList, error) {
	req, err := ctx.buildGetRequest("projects", nil)
	if err != nil {
		return nil, err
	}
	data, err := ctx.secureRequest(req)
	if err != nil {
		return nil, err
	}

	var list RemoteProjectList
	if err := json.Unmarshal(data, &list); err != nil {
		return nil, err
	}
	return &list, nil
}

func (ctx *Context) FetchProject(id int) (*RemoteProject, error) {
	path := fmt.Sprintf("projects/%d", id)
	req, err := ctx.buildGetRequest(path, nil)
	if err != nil {
		return nil, err
	}
	body, err := ctx.secureRequest(req)
	if err != nil {
		return nil, err
	}

	var data RemoteProjectData
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, err
	}
	return &data.Project, nil
}

type IssueParams struct {
	ProjectId     int
	IncludeClosed bool
}

func (params *IssueParams) Encode() url.Values {
	values := url.Values{}
	values.Add("project_id", fmt.Sprintf("%d", params.ProjectId))
	if params.IncludeClosed {
		values.Add("status_id", "*")
	}
	return values
}

func (ctx *Context) FetchIssues(params *IssueParams) (*RemoteIssueList, error) {
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

	var issues RemoteIssueList
	if err := json.Unmarshal(data, &issues); err != nil {
		return nil, err
	}
	return &issues, nil
}

func (ctx *Context) FetchIssue(id int) (*RemoteIssue, error) {
	url := fmt.Sprintf("issues/%d", id)
	req, err := ctx.buildGetRequest(url, nil)
	if err != nil {
		return nil, err
	}
	data, err := ctx.secureRequest(req)
	if err != nil {
		return nil, err
	}

	var issue RemoteIssueData
	if err := json.Unmarshal(data, &issue); err != nil {
		return nil, err
	}
	return &issue.Issue, nil
}

func (ctx *Context) FetchActivities() (map[int]string, error) {
	var (
		req        *http.Request
		data       []byte
		activities TimeEntryActivities
		err        error
	)

	if req, err = ctx.buildGetRequest("enumerations/time_entry_activities", nil); err != nil {
		return nil, err
	}
	if data, err = ctx.secureRequest(req); err != nil {
		return nil, err
	}
	if err = json.Unmarshal(data, &activities); err != nil {
		return nil, err
	}
	return activities.Map(), nil
}

type RemoteTimeEntry struct {
	IssueId    int     `json:"issue_id"`
	SpentOn    string  `json:"spent_on"`
	Hours      float64 `json:"hours"`
	ActivityId int     `json:"activity_id"`
	Comments   string  `json:"comments"`
}

type RemoteTimeEntryPayload struct {
	TimeEntry RemoteTimeEntry `json:"time_entry"`
}

func (ctx *Context) PushActivity(entry *RemoteTimeEntry) error {
	payload := &RemoteTimeEntryPayload{TimeEntry: *entry}
	jsonBody, err := json.Marshal(&payload)
	if err != nil {
		return err
	}

	reader := bytes.NewReader(jsonBody)
	req, err := ctx.buildPostRequest("time_entries", reader)
	if err != nil {
		return err
	}
	_, err = ctx.secureCreate(req)
	return err
}
