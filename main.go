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
	end := start.Add(time.Minute * 30)
	e := &calendar.Event{
		Summary:      "empty slot",
		Description:  "added via service account",
		Start:        &calendar.EventDateTime{DateTime: start.Format(time.RFC3339)},
		End:          &calendar.EventDateTime{DateTime: end.Format(time.RFC3339)},
		Transparency: "transparent",
		Status:       "tentative",
	}
	created, err := s.Events.Insert(calendarID, e).Do()
	if err != nil {
		log.Fatalf("insert: %v", err)
	}
	log.Println("created:", created.Id)

	evt, err := s.Events.Get(calendarID, created.Id).Do()
	if err != nil {
		log.Fatalf("get: %v", err)
	}
	fmt.Printf("event: %s\n", e.Summary)
	// simulate event confirmation
	evt.Description = fmt.Sprintf("confirmation time: %d", time.Now().Unix())
	evt.Transparency = "opaque"
	evt.Status = "confirmed"
	updated, err := s.Events.Update(calendarID, created.Id, evt).Do()
	if err != nil {
		log.Fatalf("update: %v", err)
	}
	fmt.Printf("updated: %s\n", updated.Description)

	//find outdated empty slots
	events, err := s.Events.List(calendarID).
		Q("empty").
		ShowDeleted(false).
		SingleEvents(true).
		OrderBy("startTime").
		TimeMax(time.Now().Format(time.RFC3339)).
		MaxResults(100).
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
		if err := s.Events.Delete(calendarID, event.Id).Do(); err != nil {
			log.Printf("delete: %v", err)
		}
	}
}
