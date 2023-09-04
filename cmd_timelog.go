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

	params := remoteTimelogParams{}
	if flagProjectId > 0 {
		params.projectId = flagProjectId
		params.criteria = "project"
	} else {
		params.issueId = uint(flagIssueId)
		params.criteria = "issue"
	}
	if err := params.TimeFilter(flagTimeSince, flagTimeUntil); err != nil {
		panic(err)
	}
	entries, err := context.fetchTimeEntries(&params)
	if err != nil {
		panic(err)
	}

	switch flagPivot {
	case "issue":
		doIssuePivotTimelog(entries)
	case "":
		doRawTimelog(entries)
	}
}

func doRawTimelog(entries *remoteTimelogList) {
	var (
		totalTime       float64 = 0
		totalTimeString string
	)

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

func doIssuePivotTimelog(entries *remoteTimelogList) {
	var (
		totals            = map[int]float64{}
		issues            = map[int]namedEntry{}
		totalTime float64 = 0
	)

	// pivot by timelog
	for _, entry := range entries.Entries {
		if _, ok := totals[entry.Issue.Id]; !ok {
			totals[entry.Issue.Id] = 0
		}
		totals[entry.Issue.Id] += entry.Hours

		if _, ok := issues[entry.Issue.Id]; !ok {
			issues[entry.Issue.Id] = entry.Issue
		}
	}

	tbl := table.New("ISSUE", "SUBJECT", "TIME")
	for id, total := range totals {
		totalTime += total
		var hours string
		if flagDecimalTime {
			hours = fmt.Sprintf("%.02f", total)
		} else {
			hours = formatAsTime(total)
		}
		tbl.AddRow(id, issues[id].Name, hours)
	}

	var totalHours string
	if flagDecimalTime {
		totalHours = fmt.Sprintf("%.02f", totalTime)
	} else {
		totalHours = formatAsTime(totalTime)
	}
	tbl.AddRow("", "----------", strings.Repeat("-", len(totalHours)))
	tbl.AddRow("", "TOTAL TIME", totalHours)
	tbl.Print()
}
