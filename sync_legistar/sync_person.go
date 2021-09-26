package main

import (
	"log"
	"path/filepath"
	"time"

	"github.com/jehiah/legislator/db"
	"github.com/jehiah/legislator/legistar"
)

func (s *SyncApp) SyncPersons() error {

	persons, err := s.legistar.Persons(legistar.PersonLastModifiedFilter(s.LastSync.Persons))
	slugs := make(map[string]bool)

	for _, p := range persons {
		officeRecords, err := s.legistar.PersonOfficeRecords(p.ID)
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

		if err := s.writeFile(filepath.Join("people", record.Slug+".json"), record); err != nil {
			return err
		}
		time.Sleep(50 * time.Millisecond)
	}
	s.LastSync.Persons = db.Max(persons, func(i int) time.Time { return persons[i].LastModified.Time })

	return err
}
