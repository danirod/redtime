package main

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

type trackInteractiveOption int

const (
	TrackQuit trackInteractiveOption = iota
	TrackCommit
	TrackDiscard
)

func pushSpool(spool *spoolfile) {
	entry := spool.CreateTimeEntry(time.Now())
	if err := context.pushActivity(entry); err != nil {
		fmt.Println("")
		fmt.Println("Something bad has happened, maybe you need to do it yourself now")
		panic(err)
	}
}

func showConfirmPush() bool {
	fmt.Print(`Push this entry? Type "cancel" for cancel, otherwise, it will push: `)
	var option string
	if _, err := fmt.Scanln(&option); err != nil {
		// If presses enter without typing anything, the error is 'unexpected newline'
		return err.Error() == "unexpected newline"
	}
	return strings.ToLower(option) != "cancel"
}

func showSpoolMenu() trackInteractiveOption {
	fmt.Println(`What do you want to do?
	(C) Commit existing spool and start new
	(D) Discard existing spool and start new
	(Q) Quit and don't do anything`)
	fmt.Print("Choose an option: ")
	var option string
	if opt, err := fmt.Scanln(&option); opt != 1 || err != nil || len(option) > 1 {
		return TrackQuit
	}
	switch lower := strings.ToLower(option); lower[0] {
	case 'c':
		return TrackCommit
	case 'd':
		return TrackDiscard
	default:
		return TrackQuit
	}
}

func countElapsed(since int64) string {
	seconds := time.Now().Unix() - since
	hh, mmss := seconds/3600, seconds%3600
	mm, ss := mmss/60, mmss%60
	return fmt.Sprintf("%d:%02d:%02d", hh, mm, ss)
}

func doTrack() {
	if flagIssueId == 0 {
		fmt.Println("Missing issue ID")
		flagSetTrack.PrintDefaults()
		os.Exit(1)
	}

	if flagActivityId == 0 {
		// Check if there is a default activity
		act, err := context.fetchDefaultActivity()
		if err != nil {
			panic(err)
		}
		if act == nil {
			fmt.Println("Missing activity ID and no default activity!")
			flagSetTrack.PrintDefaults()
			os.Exit(1)
		}
		flagActivityId = act.Id
	}

	var spool spoolfile
	exists, err := spool.Exists()
	if err != nil {
		panic(err)
	}
	if exists {
		if err := spool.ReadToDefault(); err != nil {
			panic(err)
		}
		fmt.Println("An spool file already exists!")
		spool.Print()
		switch showSpoolMenu() {
		case TrackQuit:
			os.Exit(1)
		case TrackCommit:
			if err := spool.ReadToDefault(); err != nil {
				panic(err)
			}
			pushSpool(&spool)
		}
	}

	spool = spoolfile{
		ActivityId: flagActivityId,
		IssueId:    flagIssueId,
		Comment:    flagComments,
		StartDate:  time.Now().Unix(),
	}

	if err := spool.WriteToDefault(); err != nil {
		panic(err)
	}

	spool.Print()

	fmt.Println("Counting time. Press Ctrl-C to stop tracking time")

	elapsed := countElapsed(spool.StartDate)
	fmt.Printf("\rT: %d | E: %s", flagIssueId, elapsed)

	ticker := time.NewTicker(time.Second)

	ctrlChannel := make(chan os.Signal, 1)
	signal.Notify(ctrlChannel, syscall.SIGINT)

	keepRunning := true
	for keepRunning {
		select {
		case <-ticker.C:
			elapsed = countElapsed(spool.StartDate)
			fmt.Printf("\rT: %d | E: %s", flagIssueId, elapsed)
		case <-ctrlChannel:
			keepRunning = false
		}
	}

	fmt.Println("Stopped counting time")
	if showConfirmPush() {
		pushSpool(&spool)
		fmt.Println("Successfully registered")
	}

	if err := spool.Delete(); err != nil {
		panic(err)
	}
	os.Exit(0)
}
