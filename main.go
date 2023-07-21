package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
)

var (
	flagSetTrack      *flag.FlagSet
	flagSetActivities *flag.FlagSet

	flagIssueId    int
	flagActivityId int
	flagComments   string

	context *Context
)

func init() {
	godotenv.Load()

	flagSetTrack = flag.NewFlagSet("track", flag.ExitOnError)
	flagSetTrack.IntVar(&flagIssueId, "issue", 0, "The issue ID")
	flagSetTrack.IntVar(&flagActivityId, "activity", 0, "The activity ID")
	flagSetTrack.StringVar(&flagComments, "comment", "", "The comment to assign")

	flagSetActivities = flag.NewFlagSet("activities", flag.ExitOnError)
}

func main() {
	var err error
	key := os.Getenv("REDMINE_API_KEY")

	if context, err = NewContext("http://localhost:3000", key); err != nil {
		panic(err)
	}

	if len(os.Args) < 2 {
		panic("ay que no dijiste nada")
	}

	switch os.Args[1] {
	case "track":
		flagSetTrack.Parse(os.Args[2:])
		doTrack()
	case "activities":
		flagSetActivities.Parse(os.Args[2:])
		doActivities()
	}
}

func doTrack() {
	fmt.Println("Counting time. Press Ctrl-C to stop tracking time")

	before := time.Now()
	ctrlChannel := make(chan os.Signal)
	signal.Notify(ctrlChannel, syscall.SIGINT)
	<-ctrlChannel
	after := time.Now()

	seconds := after.Unix() - before.Unix()
	hours := float64(seconds) / 3600.0
	if err := context.PushActivity(&RemoteTimeEntry{
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

func doActivities() {
	acts, err := context.FetchActivities()
	if err != nil {
		panic(err)
	}
	fmt.Println("Activities:")
	for id, act := range acts {
		fmt.Printf("  %4d %s\n", id, act)
	}
}
