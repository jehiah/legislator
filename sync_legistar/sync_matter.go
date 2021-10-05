package main

import (
	"log"
	"net/url"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/jehiah/legislator/db"
	"github.com/jehiah/legislator/legistar"
)

type IntroDateYearFilter struct {
	t  time.Time
	eq string
}

func (y IntroDateYearFilter) Paramters() url.Values {
	return legistar.DateTimeFilter("MatterIntroDate", y.eq, y.t)
}

func LegislationFilename(m db.Legislation) string {
	fn := strings.Fields(strings.ReplaceAll(m.File, "-", " "))[1] + ".json"
	return filepath.Join("introduction", strconv.Itoa(m.IntroDate.Year()), fn)
}

func (s *SyncApp) SyncMatter() error {

	filter := legistar.AndFilters(
		legistar.MatterLastModifiedFilter(s.LastSync.Matters),
		legistar.MatterTypeFilter("Introduction"),
		IntroDateYearFilter{time.Date(2014, time.January, 1, 0, 0, 0, 0, time.UTC), "gt"},
		// IntroDateYearFilter{time.Date(2014, time.June, 1, 0, 0, 0, 0, time.UTC), "lt"},
	)

	matters, err := s.legistar.Matters(filter)
	if err != nil {
		return err
	}

	for _, m := range matters {

		l := db.NewLegislation(m)
		fn := LegislationFilename(l)
		if s.legislationLookup[fn] {
			continue
		}
		s.legislationLookup[fn] = true

		sponsors, err := s.legistar.MatterSponsors(m.ID)
		if err != nil {
			return err
		}
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

		if err = s.writeFile(fn, l); err != nil {
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
	files, err := filepath.Glob(filepath.Join(s.targetDir, "introduction", "*", "*.json"))
	if err != nil {
		return err
	}
	for _, file := range files {
		s.legislationLookup[strings.TrimPrefix(file, s.targetDir+"/")] = true
	}
	log.Printf("loaded %d legislation files", len(s.legislationLookup))
	return nil
}
