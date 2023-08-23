package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/adrg/xdg"
)

type spoolfile struct {
	IssueId    int    `json:"issue_id"`
	ActivityId int    `json:"activity_id"`
	Comment    string `json:"comment"`
	StartDate  int64  `json:"start_date"`
}

func (s *spoolfile) CreateTimeEntry(end time.Time) *remoteTimeEntry {
	timestamp := time.Unix(s.StartDate, 0)
	date := timestamp.Format(time.DateOnly)
	seconds := end.Unix() - s.StartDate
	hours := float64(seconds) / 3600.0

	return &remoteTimeEntry{
		IssueId:    s.IssueId,
		SpentOn:    date,
		Hours:      hours,
		ActivityId: s.ActivityId,
		Comments:   s.Comment,
	}
}

func (s *spoolfile) FilePath() (string, error) {
	return xdg.DataFile("redtime/spoolfile.json")
}

func (s *spoolfile) Delete() error {
	path, err := s.FilePath()
	if err != nil {
		return err
	}
	return os.Remove(path)
}

func (s *spoolfile) Exists() (bool, error) {
	name, err := s.FilePath()
	if err != nil {
		return false, err
	}

	_, err = os.Stat(name)
	if err == nil {
		return true, nil
	}
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}
	return false, err
}

func (s *spoolfile) Read(r io.Reader) error {
	decoder := json.NewDecoder(r)
	decoder.DisallowUnknownFields()
	return decoder.Decode(s)
}

func (s *spoolfile) Write(w io.Writer) error {
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	return encoder.Encode(s)
}

func (s *spoolfile) ReadToDefault() error {
	path, err := s.FilePath()
	if err != nil {
		return err
	}

	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	return s.Read(file)
}

func (s *spoolfile) WriteToDefault() error {
	path, err := s.FilePath()
	if err != nil {
		return err
	}

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	return s.Write(file)
}

func (s *spoolfile) Print() {
	acts, err := context.fetchActivities()
	if err != nil {
		panic(err)
	}

	issue, err := context.fetchIssue(s.IssueId)
	if err != nil {
		panic(err)
	}

	startedAt := time.Unix(s.StartDate, 0)
	fmt.Printf("Comment: %s\n", s.Comment)
	fmt.Printf("Activity: %s\n", acts[s.ActivityId])
	fmt.Printf("Issue: %s\n", issue.Subject)
	fmt.Printf("Started on: %s\n", startedAt.Format(time.RFC822))
}
