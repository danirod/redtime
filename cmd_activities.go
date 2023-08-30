package main

import "github.com/rodaine/table"

func doListActivities() {
	acts, err := context.fetchActivities()
	if err != nil {
		panic(err)
	}

	tbl := table.New("ID", "ACTIVITY", "DEF", "ACT")
	for id, act := range acts {
		isDefault := "   "
		isActive := "   "
		if act.IsDefault {
			isDefault = "[X]"
		}
		if act.Active {
			isActive = "[X]"
		}
		tbl.AddRow(id, act.Name, isDefault, isActive)
	}
	tbl.Print()
}
