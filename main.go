package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

const timeLayout = "20060102T150405Z"

type Event struct {
	Start   time.Time
	End     time.Time
	Summary string
}

func (e Event) Duration() float64 {
	return e.End.Sub(e.Start).Hours()
}

func GetEvents(url string) (events []Event, err error) {
	var (
		eventChan = make(chan Event)
		lineChan  = make(chan string)
	)

	// Absorb events
	go func(eventChan <-chan Event) {
		for event := range eventChan {
			events = append(events, event)
		}
	}(eventChan)

	// Build events
	go eventBuilder(lineChan, eventChan)

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

	return
}

func eventBuilder(lineChan <-chan string, eventChan chan<- Event) {
	var tmp Event
	for line := range lineChan {
		switch {
		case line == "BEGIN:VEVENT":
			tmp = Event{}
		case line == "END:VEVENT":
			eventChan <- tmp
		case strings.HasPrefix(line, "DTSTART:"):
			tmp.Start = parseDate(line)
		case strings.HasPrefix(line, "DTEND:"):
			tmp.End = parseDate(line)
		case strings.HasPrefix(line, "SUMMARY:"):
			var s string
			s = strings.SplitN(line, ":", 2)[1]
			s = strings.Replace(s, "\\,", ",", -1)
			tmp.Summary = s
		}
	}
}

func parseDate(s string) (t time.Time) {
	var err error
	datetime := strings.Split(s, ":")[1]
	if t, err = time.Parse(timeLayout, datetime); err != nil {
		log.Fatal(err)
	}
	return
}

func main() {
	var url string
	flag.StringVar(&url, "url", "", "Private URL of Google Calendar to process")
	flag.Parse()

	events, err := GetEvents(url)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Event Count: %d\n", len(events))
	for _, event := range events {
		fmt.Printf("%-64s %0.2f\n", event.Summary, event.Duration())
	}
}
