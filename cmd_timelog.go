package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/rodaine/table"
)

func formatAsTime(decimal float64) string {
	hours := int(decimal)
	minutes := int((decimal - float64(hours)) * 60)
	return fmt.Sprintf("%d:%02d", hours, minutes)
}

func doListTimelog() {
	if (flagProjectId == 0 && flagIssueId == 0) || (flagProjectId != 0 && flagIssueId != 0) {
		fmt.Println("Error: must specify either -project or -issue")
		flagSetTimelog.PrintDefaults()
		os.Exit(1)
	}

	var (
		totalTime       float64 = 0
		totalTimeString string
	)

	params := remoteTimelogParams{}
	if flagProjectId > 0 {
		params.projectId = flagProjectId
		params.criteria = "project"
	} else {
		params.issueId = uint(flagIssueId)
		params.criteria = "issue"
	}
	entries, err := context.fetchTimeEntries(&params)
	if err != nil {
		panic(err)
	}

	tbl := table.New("ID", "USER", "ACTIVITY", "ISSUE", "COMMENT", "DATE", "TIME")
	for _, entry := range entries.Entries {
		totalTime += entry.Hours
		var hours string
		if flagDecimalTime {
			hours = fmt.Sprintf("%.02f", entry.Hours)
		} else {
			hours = formatAsTime(entry.Hours)
		}
		tbl.AddRow(
			entry.Id, entry.User.Name, entry.Activity.Name,
			entry.Issue.Name, entry.Comment, entry.SpentOn, hours,
		)
	}

	if flagDecimalTime {
		totalTimeString = fmt.Sprintf("%.02f", totalTime)
	} else {
		totalTimeString = formatAsTime(totalTime)
	}
	tbl.AddRow("", "", "", "", "", "----------", strings.Repeat("-", len(totalTimeString)))
	tbl.AddRow("", "", "", "", "", "TOTAL TIME", totalTimeString)

	tbl.Print()
}
