package main

import "fmt"

func doNewIssue() {
	if flagTitle == "" {
		fmt.Println("Title cannot be blank")
		flagSetNewIssue.PrintDefaults()
		return
	}
	if flagProjectId == 0 {
		fmt.Println("Project ID cannot be empty")
		flagSetNewIssue.PrintDefaults()
		return
	}
	if flagTrackerId == 0 {
		fmt.Println("Tracker ID cannot be empty")
		flagSetNewIssue.PrintDefaults()
		return
	}

	payload := remoteNewIssuePayloadData{}
	payload.Issue.Subject = flagTitle
	payload.Issue.ProjectId = flagProjectId
	payload.Issue.ParentIssueId = flagParentIssueId
	payload.Issue.TrackerId = flagTrackerId

	if err := context.pushNewIssue(payload); err != nil {
		panic(err)
	}

	fmt.Println("Created")
}
