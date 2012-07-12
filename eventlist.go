package main

import (
	"bufio"
	"bytes"
	"log"
	"net/http"
	"sort"
	"strings"
	"time"
)

const timeLayout = "20060102T150405Z"

type EventList []Event

func (l *EventList) Append(e Event) {
	*l = append(*l, e)
}

func (l EventList) Len() int           { return len(l) }
func (l EventList) Less(i, j int) bool { return l[i].Summary < l[j].Summary }
func (l EventList) Swap(i, j int)      { l[i], l[j] = l[j], l[i] }

func GetEvents(url string) (events EventList, err error) {
	var (
		eventChan = make(chan Event)
		lineChan  = make(chan string)
	)

	// Absorb events
	go func(eventChan <-chan Event) {
		for event := range eventChan {
			events.Append(event)
		}
	}(eventChan)

	// Build events
	go EventBuilder(lineChan, eventChan)

	response, err := http.Get(url)
	defer response.Body.Close()
	if err != nil {
		return
	}

	reader := bufio.NewReader(response.Body)
	buffer := bytes.NewBuffer(make([]byte, 1024))

	for {
		part, prefix, err := reader.ReadLine()
		if err != nil {
			break
		}
		buffer.Write(part)
		if !prefix {
			lineChan <- buffer.String()
			buffer.Reset()
		}
	}

	sort.Sort(events)

	return
}

func EventBuilder(in <-chan string, out chan<- Event) {
	var tmp Event
	for line := range in {
		switch {
		case line == "BEGIN:VEVENT":
			tmp = Event{}
		case line == "END:VEVENT":
			out <- tmp
		case strings.HasPrefix(line, "DTSTART:"):
			tmp.Start = ParseDate(line)
		case strings.HasPrefix(line, "DTEND:"):
			tmp.End = ParseDate(line)
		case strings.HasPrefix(line, "SUMMARY:"):
			var s string
			s = strings.SplitN(line, ":", 2)[1]
			s = strings.Replace(s, "\\,", ",", -1)
			tmp.Summary = s
		}
	}
}

func ParseDate(s string) (t time.Time) {
	var err error

	// Trim off DT(START|END): if it's there
	if strings.Contains(s, ":") {
		s = strings.Split(s, ":")[1]
	}

	// Try to prevent error
	// parsing time "19700308T020000" as "20060102T150405Z": cannot parse "" as "Z"
	if s[len(s)-1] != 'Z' {
		s = s + "Z"
	}
	if t, err = time.Parse(timeLayout, s); err != nil {
		log.Fatal(err)
	}
	return
}
