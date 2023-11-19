package main

import "github.com/rodaine/table"

func doTrackers() {
	trackers, err := context.fetchTrackers()
	if err != nil {
		panic(err)
	}

	tbl := table.New("ID", "NAME")
	for _, tracker := range trackers.Trackers {
		tbl.AddRow(tracker.Id, tracker.Name)
	}
	tbl.Print()
}
