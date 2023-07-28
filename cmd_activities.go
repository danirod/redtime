package main

import "github.com/rodaine/table"

func doListActivities() {
	acts, err := context.fetchActivities()
	if err != nil {
		panic(err)
	}

	tbl := table.New("ID", "ACTIVITY")
	for id, act := range acts {
		tbl.AddRow(id, act)
	}
	tbl.Print()
}
