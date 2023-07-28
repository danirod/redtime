package main

import (
	"time"

	"github.com/rodaine/table"
)

func doListProjects() {
	proj, err := context.fetchProjects()
	if err != nil {
		panic(err)
	}

	header := []interface{}{"ID", "NAME", "UPDATED ON"}
	if flagFull {
		header = []interface{}{"ID", "NAME", "SLUG", "PUBLIC", "UPDATED ON"}
	}
	tbl := table.New(header...)
	for _, p := range proj.Projects {
		updatedAt, _ := time.Parse(time.RFC3339, p.UpdatedOn)
		updatedAtString := updatedAt.Format("2006-01-02 15:04:05")
		row := []interface{}{p.Id, p.Name, updatedAtString}
		if flagFull {
			public := " "
			if p.IsPublic {
				public = "X"
			}
			row = []interface{}{p.Id, p.Name, p.Identifier, public, updatedAtString}
		}
		tbl.AddRow(row...)
	}
	tbl.Print()
}
