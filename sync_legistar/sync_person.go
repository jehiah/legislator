package main

import (
	"encoding/json"
	"log"
	"os"
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
	for _, file := range files {
		b, err := os.ReadFile(file)
		if err != nil {
			return err
		}
		var p db.Person
		err = json.Unmarshal(b, &p)
		if err != nil {
			return err
		}
		s.personLookup[p.ID] = p
	}
	log.Printf("loaded %d people", len(s.personLookup))
	return nil
}
