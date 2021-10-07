package main

import (
	"context"
	"log"
	"net/url"
	"path/filepath"
	"reflect"
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
	ctx := context.Background()
	filter := legistar.AndFilters(
		legistar.MatterLastModifiedFilter(s.LastSync.Matters),
		legistar.MatterTypeFilter("Introduction"),
		IntroDateYearFilter{time.Date(2014, time.January, 1, 0, 0, 0, 0, time.UTC), "gt"},
		// IntroDateYearFilter{time.Date(2014, time.June, 1, 0, 0, 0, 0, time.UTC), "lt"},
	)

	matters, err := s.legistar.Matters(ctx, filter)
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
		err = s.updateMatter(ctx, l)
		if err != nil {
			return err
		}
	}

	// check if any have new sponsors
	if err = s.UpdateMatterSponsors(); err != nil {
		return err
	}

	if len(matters) > 0 {
		s.LastSync.Matters = db.Max(matters, func(i int) time.Time { return matters[i].LastModified.Time })
	}
	return nil
}

func (s *SyncApp) UpdateMatter(ID int) error {
	ctx := context.Background()
	m, err := s.legistar.Matter(ctx, ID)
	if err != nil {
		return err
	}
	l := db.NewLegislation(m)
	return s.updateMatter(ctx, l)
}

func (s *SyncApp) updateMatter(ctx context.Context, l db.Legislation) error {
	fn := LegislationFilename(l)

	sponsors, err := s.legistar.MatterSponsors(ctx, l.ID)
	if err != nil {
		return err
	}
	l.Sponsors = []db.PersonReference{}
	for _, p := range sponsors {
		if p.MatterVersion != l.Version {
			continue
		}
		l.Sponsors = append(l.Sponsors, s.personLookup[p.NameID].Reference())
	}

	versions, err := s.legistar.MatterTextVersions(ctx, l.ID)
	if err != nil {
		return err
	}
	l.TextID = versions.LatestTextID()
	txt, err := s.legistar.MatterText(ctx, l.ID, l.TextID)
	l.Text = txt.Plain
	l.RTF = txt.RTF

	return s.writeFile(fn, l)
}

func (s *SyncApp) LoadMatter() error {
	files, err := filepath.Glob(filepath.Join(s.targetDir, "introduction", "*", "*.json"))
	if err != nil {
		return err
	}
	for _, fn := range files {
		fn = strings.TrimPrefix(fn, s.targetDir+"/")
		s.legislationLookup[fn] = true
	}

	log.Printf("loaded %d legislation files", len(s.legislationLookup))
	return nil
}

func (s SyncApp) UpdateMatterSponsors() error {
	ctx := context.Background()
	currentSessionStart := time.Date(2014, time.January, 1, 0, 0, 0, 0, time.UTC)
	for fn := range s.legislationLookup {
		var l *db.Legislation
		err := s.readFile(fn, &l)
		if err != nil {
			return err
		}
		if l.IntroDate.Before(currentSessionStart) {
			continue
		}

		// check for new sponsors
		matterSponsors, err := s.legistar.MatterSponsors(ctx, l.ID)
		if err != nil {
			return err
		}
		var sponsors []db.PersonReference
		for _, p := range matterSponsors {
			if p.MatterVersion != l.Version {
				continue
			}
			sponsors = append(sponsors, s.personLookup[p.NameID].Reference())
		}
		if !reflect.DeepEqual(sponsors, l.Sponsors) {
			err = s.updateMatter(ctx, *l)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
