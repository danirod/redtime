package main

import (
	"encoding/json"
	"net/http"
)

type remoteTimeEntryActivities struct {
	Activities []namedEntry `json:"time_entry_activities"`
}

func (act *remoteTimeEntryActivities) Map() map[int]string {
	activities := map[int]string{}
	for _, item := range act.Activities {
		activities[item.Id] = item.Name
	}
	return activities
}

func (ctx *redmineClient) fetchActivities() (map[int]string, error) {
	var (
		req        *http.Request
		data       []byte
		activities remoteTimeEntryActivities
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
