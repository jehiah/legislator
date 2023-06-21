package main

import (
	"context"
	"log"
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

func (s *SyncApp) LoadEvents() error {
	files, err := filepath.Glob(filepath.Join(s.targetDir, "events", "*", "*.json"))
	if err != nil {
		return err
	}
	for _, fn := range files {
		fn = strings.TrimPrefix(fn, s.targetDir+"/")
		c := strings.FieldsFunc(fn, func(c rune) bool { return c == '_' || c == '.' })
		id, err := strconv.Atoi(c[len(c)-2])
		if err != nil {
			return err
		}
		s.eventsLookup[id] = append(s.eventsLookup[id], fn)
	}
	log.Printf("loaded %d events", len(s.eventsLookup))
	return nil
}

func (s *SyncApp) SyncRollCalls() error {
	ctx := context.Background()
	for ID, fns := range s.eventsLookup {
		var update bool
		fn := fns[0]

		var e *db.Event
		err := s.readFile(fn, &e)
		if err != nil {
			return err
		}
		if e.Date.Before(time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)) {
			continue
		}
		for _, i := range e.Items {
			if i.RollCallFlag != 0 {
				update = true
				break
			}
		}
		if update {
			err = s.SyncEvent(ctx, ID)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (s *SyncApp) SyncDuplicateEvents() error {
	ctx := context.Background()
	for ID, v := range s.eventsLookup {
		if len(v) != 1 {
			err := s.SyncEvent(ctx, ID)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (s *SyncApp) SyncAllEvent() error {
	for year := 2004; year <= 2004; year++ {
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

func (s *SyncApp) SyncEvent(ctx context.Context, ID int) error {
	event, err := s.legistar.Event(ctx, ID)
	if err != nil {
		return err
	}
	event.Items, err = s.legistar.EventItems(ctx, event.ID)
	if err != nil {
		return err
	}
	record := db.NewEvent(event, localTimezone)

	for i, v := range record.Items {
		if v.RollCallFlag == 0 {
			continue
		}
		rc, err := s.legistar.EventRollCalls(ctx, v.ID)
		if err != nil {
			return err
		}
		record.Items[i].RollCall = db.NewRollCalls(rc)
	}

	fn := EventFilename(record)
	for _, existingFile := range s.eventsLookup[event.ID] {
		if existingFile == fn {
			continue
		}
		// remove files where the name has changed
		err := s.removeFile(existingFile)
		if err != nil {
			return err
		}
	}
	if err := s.writeFile(fn, record); err != nil {
		return err
	}
	return nil
}

func (s *SyncApp) SyncEvents(filter legistar.Filters) error {
	if filter == nil {
		filter = legistar.EventLastModifiedFilter(s.LastSync.Events)
		// filter = legistar.EventLastModifiedFilter(time.Date(2022, 11, 1, 0, 0, 0, 0, time.UTC))
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

		for i, v := range record.Items {
			if v.RollCallFlag == 0 {
				continue
			}
			rc, err := s.legistar.EventRollCalls(ctx, v.ID)
			if err != nil {
				return err
			}
			record.Items[i].RollCall = db.NewRollCalls(rc)
		}

		fn := EventFilename(record)
		for _, existingFile := range s.eventsLookup[event.ID] {
			if existingFile == fn {
				continue
			}
			// remove files where the name has changed
			err := s.removeFile(existingFile)
			if err != nil {
				return err
			}
		}
		if err := s.writeFile(fn, record); err != nil {
			return err
		}
	}
	if len(events) > 0 {
		s.LastSync.Events = db.Max(events, func(i int) time.Time { return events[i].LastModified.Time })
	}

	return nil
}
