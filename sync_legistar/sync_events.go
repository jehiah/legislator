package main

import (
	"context"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gosimple/slug"
	"github.com/jehiah/legislator/db"
	"github.com/jehiah/legislator/legistar"
)

func EventFilename(event db.Event) string {
	body := slug.MakeLang(strings.TrimPrefix(event.BodyName, "Committee on "), "en")
	fn := event.Date.Format("2006-01-02_15_04") + "_" + body + "_" + strconv.Itoa(event.ID) + ".json"
	return filepath.Join("events", strconv.Itoa(event.Date.Year()), fn)
}

func (s *SyncApp) SyncAllEvent() error {
	for year := 2018; year <= 2023; year++ {
		for month := time.January; month <= time.December; month++ {
			start := time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
			end := start.AddDate(0, 1, 1)
			filter := legistar.AndFilters(
				legistar.EventDateFilter{"ge", start},
				legistar.EventDateFilter{"lt", end},
			)
			err := s.SyncEvents(filter)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (s *SyncApp) SyncEvents(filter legistar.Filters) error {
	if filter == nil {
		// filter = legistar.EventLastModifiedFilter(s.LastSync.Events)
		filter = legistar.EventLastModifiedFilter(time.Date(2022, 11, 1, 0, 0, 0, 0, time.UTC))
	}

	ctx := context.Background()
	events, err := s.legistar.Events(ctx, filter)

	if err != nil {
		return err
	}

	for _, event := range events {
		event.Items, err = s.legistar.EventItems(ctx, event.ID)
		if err != nil {
			return err
		}
		record := db.NewEvent(event, localTimezone)
		if err := s.writeFile(EventFilename(record), record); err != nil {
			return err
		}
	}
	if len(events) > 0 {
		s.LastSync.Events = db.Max(events, func(i int) time.Time { return events[i].LastModified.Time })
	}

	return nil
}

func (s *SyncApp) LoadEvents() error {
	return nil
}
