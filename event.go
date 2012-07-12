package main

import (
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
