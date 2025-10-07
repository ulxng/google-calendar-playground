package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
)

func main() {
	keyFile := os.Getenv("SERVICE_ACCOUNT_KEY_FILE")
	calendarID := os.Getenv("CALENDAR_ID")
	if calendarID == "" {
		log.Fatal("CALENDAR_ID is required")
	}

	ctx := context.Background()
	s, err := calendar.NewService(ctx, option.WithCredentialsFile(keyFile))
	if err != nil {
		log.Fatalf("newService: %v", err)
	}

	start := time.Now().Add(time.Hour * 24)
	end := start.Add(time.Hour)
	e := &calendar.Event{
		Summary:     "Demo meeting",
		Description: "added via service account",
		Start:       &calendar.EventDateTime{DateTime: start.Format(time.RFC3339)},
		End:         &calendar.EventDateTime{DateTime: end.Format(time.RFC3339)},
	}
	created, err := s.Events.Insert(calendarID, e).Do()
	if err != nil {
		log.Fatalf("insert: %v", err)
	}
	log.Println("created:", created.Id)

	events, err := s.Events.List(calendarID).
		ShowDeleted(false).
		SingleEvents(true).
		OrderBy("startTime").
		TimeMin(time.Now().Format(time.RFC3339)).
		MaxResults(10).
		Do()
	if err != nil {
		log.Fatalf("list: %v", err)
	}

	if len(events.Items) == 0 {
		fmt.Println("no events found")
		return
	}
	for _, event := range events.Items {
		fmt.Printf("event: %v from %s to %s\n", event.Summary, event.Start.DateTime, event.End.DateTime)
	}
}
