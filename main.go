package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

type subcommand struct {
	name        string
	description string
	fs          *flag.FlagSet
	handler     func()
}
type subcommandMap []subcommand

var (
	flagSetTrack      = flag.NewFlagSet("track", flag.ExitOnError)
	flagSetActivities = flag.NewFlagSet("activities", flag.ExitOnError)
	flagSetIssues     = flag.NewFlagSet("issues", flag.ExitOnError)
	flagSetProjects   = flag.NewFlagSet("projects", flag.ExitOnError)
	flagSetTimelog    = flag.NewFlagSet("timelog", flag.ExitOnError)

	flagFull          bool
	flagIncludeClosed bool
	flagProjectId     uint
	flagIssueId       int
	flagPivot         string
	flagDecimalTime   bool
	flagActivityId    uint
	flagComments      string
	flagTimeSince     string
	flagTimeUntil     string

	context *redmineClient

	subcommands = []subcommand{
		{
			name:        "activities",
			description: "list the available activities",
			fs:          flagSetActivities,
			handler:     doListActivities,
		},
		{
			name:        "issues",
			description: "list the issues in a project",
			fs:          flagSetIssues,
			handler:     doListIssues,
		},
		{
			name:        "projects",
			description: "list the projects in the instance",
			fs:          flagSetProjects,
			handler:     doListProjects,
		},
		{
			name:        "status",
			description: "show the status of the running task",
			fs:          nil,
			handler:     doStatus,
		},
		{
			name:        "timelog",
			description: "list the time entries in a project or an issue",
			fs:          flagSetTimelog,
			handler:     doListTimelog,
		},
		{
			name:        "track",
			description: "track time",
			fs:          flagSetTrack,
			handler:     doTrack,
		},
	}
)

func init() {
	godotenv.Load()

	flagSetTrack.IntVar(&flagIssueId, "issue", 0, "The issue ID")
	flagSetTrack.UintVar(&flagActivityId, "activity", 0, "The activity ID")
	flagSetTrack.StringVar(&flagComments, "comment", "", "The comment to assign")

	flagSetIssues.UintVar(&flagProjectId, "project", 0, "The project ID")
	flagSetIssues.BoolVar(&flagIncludeClosed, "closed", false, "Include closed issues")

	flagSetTimelog.UintVar(&flagProjectId, "project", 0, "The project to get entries for")
	flagSetTimelog.IntVar(&flagIssueId, "issue", 0, "The issue to get entries for")
	flagSetTimelog.BoolVar(&flagDecimalTime, "decimal", false, "Use decimal time when reporting")
	flagSetTimelog.StringVar(&flagTimeSince, "since", "", "Since when include entries")
	flagSetTimelog.StringVar(&flagTimeUntil, "until", "", "Until when include entries")
	flagSetTimelog.StringVar(&flagPivot, "pivot", "", "Group entries by a specific field")

	flagSetProjects.BoolVar(&flagFull, "full", false, "show extra information per project")
}

func main() {
	var err error
	url := os.Getenv("REDMINE_API_URL")
	key := os.Getenv("REDMINE_API_KEY")

	if context, err = newContext(url, key); err != nil {
		panic(err)
	}

	if len(os.Args) < 2 {
		printUsage()
		return
	}
	for _, subcommand := range subcommands {
		if subcommand.name == os.Args[1] {
			if subcommand.fs != nil {
				subcommand.fs.Parse(os.Args[2:])
			}
			subcommand.handler()
			return
		}
	}
	printUsage()
}

func printUsage() {
	fmt.Printf("%s <command> <flags>\n", filepath.Base(os.Args[0]))
	fmt.Println("Commands:")
	for _, cmd := range subcommands {
		fmt.Printf("   %-12s %s\n", cmd.name, cmd.description)
	}
	fmt.Println()
}
