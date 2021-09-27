package main

import (
	"encoding/json"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/jehiah/legislator/db"
	"github.com/jehiah/legislator/legistar"
)

type IntroDateYearFilter int

func (y IntroDateYearFilter) Paramters() url.Values {
	return legistar.DateTimeFilter("MatterIntroDate", time.Date(int(y), time.January, 1, 0, 0, 0, 0, time.UTC))
}

func (s *SyncApp) SyncMatter() error {

	filter := legistar.AndFilters(
		legistar.MatterLastModifiedFilter(s.LastSync.Matters),
		legistar.MatterTypeFilter("Introduction"),
		IntroDateYearFilter(2020),
	)

	matters, err := s.legistar.Matters(filter)
	if err != nil {
		return err
	}

	for _, m := range matters {
		sponsors, err := s.legistar.MatterSponsors(m.ID)
		if err != nil {
			return err
		}

		l := db.NewLegislation(m)
		for _, p := range sponsors {
			if p.MatterVersion != m.Version {
				continue
			}
			l.Sponsors = append(l.Sponsors, s.personLookup[p.NameID].Reference())
		}

		versions, err := s.legistar.MatterTextVersions(m.ID)
		if err != nil {
			return err
		}
		l.TextID = versions.LatestTextID()
		txt, err := s.legistar.MatterText(m.ID, l.TextID)
		l.Text = txt.Plain
		l.RTF = txt.RTF

		fn := strings.Fields(strings.ReplaceAll(m.File, "-", " "))[1] + ".json"
		if err = s.writeFile(filepath.Join("introduction", strconv.Itoa(m.IntroDate.Year()), fn), l); err != nil {
			return err
		}
		time.Sleep(250 * time.Millisecond)
	}

	if len(matters) > 0 {
		s.LastSync.Matters = db.Max(matters, func(i int) time.Time { return matters[i].LastModified.Time })
	}
	return nil
}

func (s *SyncApp) LoadMatter() error {
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
