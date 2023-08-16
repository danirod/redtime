package main

import (
	"testing"
	"time"
)

func TestSpoolCreatesTimeEntry(t *testing.T) {
	endTimestamp := int64(1692213795)
	end := time.Unix(endTimestamp, 0)
	startTimestamp := endTimestamp - 7200

	spool := spoolfile{
		IssueId:    5,
		ActivityId: 8,
		Comment:    "doing something",
		StartDate:  startTimestamp,
	}
	entry := spool.CreateTimeEntry(end)

	if entry.ActivityId != 8 {
		t.Errorf("activityId is wrong")
	}
	if entry.Comments != "doing something" {
		t.Errorf("comments are wrong")
	}
	if entry.Hours != 2 {
		t.Errorf("hours are wrong")
	}
	if entry.IssueId != 5 {
		t.Errorf("issue id is wrong")
	}
	if entry.SpentOn != "2023-08-16" {
		t.Errorf("spentOn is wrong")
	}
}
