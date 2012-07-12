package main

import (
	"flag"
	"fmt"
	"log"
)

const timeLayout = "20060102T150405Z"

func main() {
	var url string
	flag.StringVar(&url, "url", "", "Private URL of Google Calendar ICS to process")
	flag.Parse()

	events, err := GetEvents(url)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Event Count: %d\n", len(events))
	for _, event := range events {
		fmt.Printf("%-72s %5.2f\n", event.Summary, event.Duration())
	}
}
