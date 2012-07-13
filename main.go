package main

import (
	"flag"
	"fmt"
	"log"
	"strings"
	"time"
)

type DateRange struct {
	Start, End time.Time
}

func prettyPrint(events EventList) {
	var (
		linef      = "%-73s%6.2f\n"
		hr         = strings.Repeat("-", 73+6)
		monthTotal = make(map[string]float64)
		poTotal    = make(map[string]float64)
		total      = 0.00
	)

	for _, e := range events {
		d := e.Duration()
		total += d
		monthTotal[e.Start.Format("January 2006")] += d
		poTotal[e.PO()] += d
		fmt.Printf(linef, e.Start.Format("Jan _2, 2006 - ")+e.Summary, d)
	}

	fmt.Println(hr)
	for m, d := range monthTotal {
		fmt.Printf(linef, m, d)
	}
	fmt.Println(hr)
	for p, d := range poTotal {
		fmt.Printf(linef, string(p), d)
	}
	fmt.Println(hr)
	fmt.Printf(linef, "Overall Total", total)
	fmt.Printf(linef, "Entry Count", float64(len(events)))
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

	prettyPrint(events)
}
