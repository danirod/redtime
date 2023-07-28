package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func doTrack() {
	if flagIssueId == 0 {
		fmt.Println("Missing issue ID")
		flagSetTrack.PrintDefaults()
		os.Exit(1)
	}
	fmt.Println("Counting time. Press Ctrl-C to stop tracking time")

	before := time.Now()
	ctrlChannel := make(chan os.Signal, 1)
	signal.Notify(ctrlChannel, syscall.SIGINT)
	<-ctrlChannel
	after := time.Now()

	seconds := after.Unix() - before.Unix()
	hours := float64(seconds) / 3600.0
	if err := context.pushActivity(&remoteTimeEntry{
		IssueId:    flagIssueId,
		SpentOn:    time.Now().Format("2006-01-02"),
		Hours:      hours,
		ActivityId: flagActivityId,
		Comments:   flagComments,
	}); err != nil {
		fmt.Println("")
		fmt.Println("Something bad has happened, maybe you need to do it yourself now")
		fmt.Println("Seconds spent here:", seconds)
		panic(err)
	}
	fmt.Println("Successfully registered")
	os.Exit(0)
}
