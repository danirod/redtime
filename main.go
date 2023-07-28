package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/rodaine/table"
)

type subcommand struct {
	name    string
	fs      *flag.FlagSet
	handler func()
}
type subcommandMap []subcommand

var (
	flagSetTrack      = flag.NewFlagSet("track", flag.ExitOnError)
	flagSetActivities = flag.NewFlagSet("activities", flag.ExitOnError)
	flagSetIssues     = flag.NewFlagSet("issues", flag.ExitOnError)
	flagSetProjects   = flag.NewFlagSet("projects", flag.ExitOnError)

	flagFull          bool
	flagIncludeClosed bool
	flagProjectId     uint
	flagIssueId       int
	flagActivityId    int
	flagComments      string

	context *redmineClient

	subcommands = []subcommand{
		{name: "projects", fs: flagSetProjects, handler: doListProjects},
		{name: "issues", fs: flagSetIssues, handler: doListIssues},
		{name: "activities", fs: flagSetActivities, handler: doListActivities},
		{name: "track", fs: flagSetTrack, handler: doTrack},
	}
)

func init() {
	godotenv.Load()

	flagSetTrack.IntVar(&flagIssueId, "issue", 0, "The issue ID")
	flagSetTrack.IntVar(&flagActivityId, "activity", 0, "The activity ID")
	flagSetTrack.StringVar(&flagComments, "comment", "", "The comment to assign")

	flagSetIssues.UintVar(&flagProjectId, "project", 0, "The project ID")
	flagSetIssues.BoolVar(&flagIncludeClosed, "closed", false, "Include closed issues")

	flagSetProjects.BoolVar(&flagFull, "full", false, "show extra information per project")
}

func main() {
	var err error
	key := os.Getenv("REDMINE_API_KEY")

	if context, err = newContext("http://localhost:3000", key); err != nil {
		panic(err)
	}

	if len(os.Args) < 2 {
		printUsage()
		return
	}
	for _, subcommand := range subcommands {
		if subcommand.name == os.Args[1] {
			subcommand.fs.Parse(os.Args[2:])
			subcommand.handler()
			return
		}
	}
	printUsage()
}

func printUsage() {
	fmt.Print("Subcommands: ")
	for _, subcommand := range subcommands {
		fmt.Print(subcommand.name + " ")
	}
	fmt.Println()
}

func doTrack() {
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

func doListIssues() {
	if flagProjectId == 0 {
		fmt.Println("Missing project ID")
		flagSetIssues.PrintDefaults()
		return
	}

	params := remoteIssueParams{ProjectId: int(flagProjectId)}
	if flagIncludeClosed {
		params.IncludeClosed = true
	}
	issues, err := context.fetchIssues(&params)
	if err != nil {
		panic(err)
	}

	tbl := table.New("ID", "SUBJECT", "ASSIGNED", "TYPE", "STATUS", "PRIORITY", "LAST UPDATE")
	for _, issue := range issues.Issues {
		updatedAt, _ := time.Parse(time.RFC3339, issue.UpdatedOn)
		updatedAtString := updatedAt.Format("2006-01-02 15:04:05")
		tbl.AddRow(issue.Id, issue.Subject, issue.AssignedTo.Name, issue.Tracker.Name, issue.Status.Name, issue.Priority.Name, updatedAtString)
	}
	tbl.Print()
}
