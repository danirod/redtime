package main

import (
	"encoding/json"
	"net/http"
)

type remoteActivityType struct {
	Id        uint   `json:"id"`
	Name      string `json:"name"`
	IsDefault bool   `json:"is_default"`
	Active    bool   `json:"active"`
}

type remoteTimeEntryActivities struct {
	Activities []remoteActivityType `json:"time_entry_activities"`
}

func (act *remoteTimeEntryActivities) Map() map[uint]remoteActivityType {
	activities := map[uint]remoteActivityType{}
	for _, item := range act.Activities {
		activities[item.Id] = item
	}
	return activities
}

func (ctx *redmineClient) fetchActivities() (map[uint]remoteActivityType, error) {
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
