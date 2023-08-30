package main

import (
	"bytes"
	"encoding/json"
)

type remoteTimeEntry struct {
	IssueId    int     `json:"issue_id"`
	SpentOn    string  `json:"spent_on"`
	Hours      float64 `json:"hours"`
	ActivityId uint    `json:"activity_id"`
	Comments   string  `json:"comments"`
}

type remoteTimeEntryPayload struct {
	TimeEntry remoteTimeEntry `json:"time_entry"`
}

func (ctx *redmineClient) pushActivity(entry *remoteTimeEntry) error {
	payload := &remoteTimeEntryPayload{TimeEntry: *entry}
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
