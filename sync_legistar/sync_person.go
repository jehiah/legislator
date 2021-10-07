package main

import (
	"context"
	"log"
	"path/filepath"
	"strings"
	"time"

	"github.com/jehiah/legislator/db"
	"github.com/jehiah/legislator/legistar"
)

func (s *SyncApp) SyncPersons() error {
	ctx := context.Background()
	persons, err := s.legistar.Persons(ctx,
		legistar.PersonLastModifiedFilter(s.LastSync.Persons),
	)
	slugs := make(map[string]bool)

	for _, p := range persons {
		officeRecords, err := s.legistar.PersonOfficeRecords(ctx, p.ID)
		if err != nil {
			return err
		}
		if len(officeRecords) == 0 {
			// skip individuals with no office
			continue
		}

		// ensure slug is unique
		if slugs[p.Slug()] {
			log.Fatalf("slug %s is a duplciate %#v", p.Slug(), p)
		}
		slugs[p.Slug()] = true

		record := db.NewPerson(p, officeRecords)
		s.personLookup[record.ID] = record

		if err := s.writeFile(filepath.Join("people", record.Slug+".json"), record); err != nil {
			return err
		}
		time.Sleep(50 * time.Millisecond)
	}
	if len(persons) > 0 {
		s.LastSync.Persons = db.Max(persons, func(i int) time.Time { return persons[i].LastModified.Time })
	}

	return err
}

func (s *SyncApp) LoadPersons() error {
	// load persons
	files, err := filepath.Glob(filepath.Join(s.targetDir, "people", "*.json"))
	if err != nil {
		return err
	}
	for _, fn := range files {
		fn = strings.TrimPrefix(fn, s.targetDir+"/")
		var p db.Person
		err = s.readFile(fn, &p)
		if err != nil {
			return err
		}
		s.personLookup[p.ID] = p
	}
	log.Printf("loaded %d people", len(s.personLookup))
	return nil
}
