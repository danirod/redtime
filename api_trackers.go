package main

import "encoding/json"

type remoteTrackerEntry struct {
	namedEntry
}

type remoteTrackerList struct {
	Trackers []remoteTrackerEntry `json:"trackers"`
}

func (ctx *redmineClient) fetchTrackers() (*remoteTrackerList, error) {
	req, err := ctx.buildGetRequest("trackers", nil)
	if err != nil {
		return nil, err
	}
	data, err := ctx.secureRequest(req)
	if err != nil {
		return nil, err
	}

	var list remoteTrackerList
	if err := json.Unmarshal(data, &list); err != nil {
		return nil, err
	}
	return &list, nil
}
