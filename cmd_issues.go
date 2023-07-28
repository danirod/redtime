package main

import (
	"fmt"
	"time"

	"github.com/rodaine/table"
)

func doListIssues() {
	if flagProjectId == 0 {
		fmt.Println("Missing project ID")
		flagSetIssues.PrintDefaults()
		return
	}

	params := remoteIssueParams{ProjectId: int(flagProjectId)}
	if flagIncludeClosed {
		params.IncludeClosed = true
	}
	issues, err := context.fetchIssues(&params)
	if err != nil {
		panic(err)
	}

	tbl := table.New("ID", "SUBJECT", "ASSIGNED", "TYPE", "STATUS", "PRIORITY", "LAST UPDATE")
	for _, issue := range issues.Issues {
		updatedAt, _ := time.Parse(time.RFC3339, issue.UpdatedOn)
		updatedAtString := updatedAt.Format("2006-01-02 15:04:05")
		tbl.AddRow(issue.Id, issue.Subject, issue.AssignedTo.Name, issue.Tracker.Name, issue.Status.Name, issue.Priority.Name, updatedAtString)
	}
	tbl.Print()
}
