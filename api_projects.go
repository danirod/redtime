package main

import (
	"encoding/json"
	"fmt"
)

type remoteProject struct {
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

type remoteProjectData struct {
	Project remoteProject `json:"project"`
}

type remoteProjectList struct {
	Projects   []remoteProject `json:"projects"`
	TotalCount int             `json:"total_count"`
	Offset     int             `json:"offset"`
	Limit      int             `json:"limit"`
}

func (ctx *redmineClient) fetchProjects() (*remoteProjectList, error) {
	req, err := ctx.buildGetRequest("projects", nil)
	if err != nil {
		return nil, err
	}
	data, err := ctx.secureRequest(req)
	if err != nil {
		return nil, err
	}

	var list remoteProjectList
	if err := json.Unmarshal(data, &list); err != nil {
		return nil, err
	}
	return &list, nil
}

func (ctx *redmineClient) fetchProjectById(id int) (*remoteProject, error) {
	path := fmt.Sprintf("projects/%d", id)
	req, err := ctx.buildGetRequest(path, nil)
	if err != nil {
		return nil, err
	}
	body, err := ctx.secureRequest(req)
	if err != nil {
		return nil, err
	}

	var data remoteProjectData
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, err
	}
	return &data.Project, nil
}
