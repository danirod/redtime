package main

import "github.com/rodaine/table"

func doIssueStatuses() {
	statuses, err := context.fetchStatuses()
	if err != nil {
		panic(err)
	}

	tbl := table.New("ID", "NAME")
	for _, tracker := range statuses.Statuses {
		tbl.AddRow(tracker.Id, tracker.Name)
	}
	tbl.Print()
}
