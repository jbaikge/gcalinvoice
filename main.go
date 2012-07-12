package main

import (
	"flag"
	"fmt"
	"log"
	"time"
)

type DateRange struct {
	Start, End time.Time
}

func main() {
	var start, url string
	flag.StringVar(&start, "start", "1970-01-01", "Start date of ")
	flag.StringVar(&url, "url", "", "Private URL of Google Calendar ICS to process")
	flag.Parse()

	rStart, err := time.Parse("2006-01-02", start)
	if err != nil {
		log.Fatal("Invalid time: " + start)
	}
	dateRange := DateRange{rStart, time.Now()}

	events, err := GetEvents(url, dateRange)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Event Count: %d\n", len(events))
	for _, event := range events {
		fmt.Printf("%-72s %5.2f\n", event.Summary, event.Duration())
	}
}
