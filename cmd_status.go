package main

import (
	"fmt"
	"os"
)

func doStatus() {
	var spool spoolfile
	exists, err := spool.Exists()
	if err != nil {
		panic(err)
	}
	if !exists {
		fmt.Println("No task is active")
		os.Exit(1)
	}

	if err := spool.ReadToDefault(); err != nil {
		panic(err)
	}

	spool.Print()
}
