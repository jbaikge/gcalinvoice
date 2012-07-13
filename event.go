package main

import (
	"regexp"
	"time"
)

type Event struct {
	Start   time.Time
	End     time.Time
	Summary string
}

func (e Event) Duration() float64 {
	return e.End.Sub(e.Start).Hours()
}

func (e Event) PO() string {
	return regexp.MustCompile("^#(\\d+)").FindString(e.Summary)
}
