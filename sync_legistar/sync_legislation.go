package main

import (
	"context"
	"fmt"
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

type MatterDateYearFilter struct {
	t  time.Time
	eq string
}

func (y MatterDateYearFilter) Paramters() url.Values {
	return legistar.DateTimeFilter("MatterIntroDate", y.eq, y.t)
}

func LegislationFilename(m db.Legislation) string {
	fn := strings.Fields(strings.ReplaceAll(m.File, "-", " "))[1] + ".json"
	return filepath.Join("introduction", strconv.Itoa(m.IntroDate.Year()), fn)
}

func (s *SyncApp) SyncLegislation() error {
	ctx := context.Background()
	filter := legistar.AndFilters(
		legistar.MatterLastModifiedFilter(s.LastSync.Matters),
		legistar.MatterTypeFilter("Introduction"),
		MatterDateYearFilter{time.Date(2014, time.January, 1, 0, 0, 0, 0, time.UTC), "gt"},
		// MatterDateYearFilter{time.Date(2014, time.June, 1, 0, 0, 0, 0, time.UTC), "lt"},
	)

	matters, err := s.legistar.Matters(ctx, filter)
	if err != nil {
		return err
	}

	for _, m := range matters {
		// temporary items - these appear to be working drafts
		if strings.HasPrefix(m.File, "T") {
			continue
		}
		l := db.NewLegislation(m)
		fn := LegislationFilename(l)
		s.legislationLookup[fn] = true
		err = s.updateLegislation(ctx, l)
		if err != nil {
			return err
		}
	}

	// // check if any have new sponsors
	// if err = s.UpdateMatterSponsors(); err != nil {
	// 	return err
	// }

	if len(matters) > 0 {
		s.LastSync.Matters = db.Max(matters, func(i int) time.Time { return matters[i].LastModified.Time })
	}
	return nil
}

// UpdateMatterByFile expects format 1234-2020
func (s *SyncApp) UpdateLegislationByFile(q string) error {
	ctx := context.Background()
	file := fmt.Sprintf("Int %s", q)
	filter := legistar.AndFilters(
		legistar.MatterTypeFilter("Introduction"),
		legistar.MatterFileFilter(file),
	)

	matters, err := s.legistar.Matters(ctx, filter)
	if err != nil {
		return err
	}
	if len(matters) != 1 {
		return fmt.Errorf("expected 1 response got %d for %q", len(matters), q)
	}
	return s.UpdateLegislation(ctx, matters[0].ID)
}

func (s *SyncApp) UpdateAllLegislation() error {
	ctx := context.Background()
	for fn := range s.legislationLookup {
		var l *db.Legislation
		err := s.readFile(fn, &l)
		if err != nil {
			return err
		}

		// TEMP
		update := false
		cutoffLow := time.Date(2006, 1, 1, 0, 0, 0, 0, time.UTC)
		cutoffHigh := time.Date(2008, 1, 1, 0, 0, 0, 0, time.UTC)
		for _, h := range l.History {
			if h.PassedFlagName != "" {
				update = true
			}
		}
		if l.IntroDate.Before(cutoffLow) || l.IntroDate.After(cutoffHigh) {
			update = false
		}
		if l.File == "Int 0683-2015" {
			update = false
		}
		if !update {
			continue
		}

		err = s.UpdateLegislationWithRetry(ctx, l.ID)
		if err != nil {
			log.Printf("Got error after retry; sleeping 30s and skipping %s %s", l.File, err)
			time.Sleep(time.Second * 30)
			// return err
		}

	}
	return nil
}

func (s *SyncApp) UpdateLegislationWithRetry(ctx context.Context, ID int) error {
	m, err := s.legistar.Matter(ctx, ID)
	if err != nil {
		log.Print(err)
		time.Sleep(time.Second)
		m, err = s.legistar.Matter(ctx, ID)
	}
	if err != nil {
		return err
	}
	l := db.NewLegislation(m)
	err = s.updateLegislation(ctx, l)
	if err != nil {
		log.Print(err)
		time.Sleep(time.Second)
		err = s.updateLegislation(ctx, l)
	}
	return err
}

func (s *SyncApp) UpdateLegislation(ctx context.Context, ID int) error {
	m, err := s.legistar.Matter(ctx, ID)
	if err != nil {
		return err
	}
	l := db.NewLegislation(m)
	return s.updateLegislation(ctx, l)
}

func (s *SyncApp) updateLegislation(ctx context.Context, l db.Legislation) error {
	fn := LegislationFilename(l)
	return s.updateMatter(ctx, fn, l)
}

func (s *SyncApp) LoadLegislation() error {
	files, err := filepath.Glob(filepath.Join(s.targetDir, "introduction", "*", "*.json"))
	if err != nil {
		return err
	}
	for _, fn := range files {
		if filepath.Base(fn) == "index.json" {
			continue
		}
		fn = strings.TrimPrefix(fn, s.targetDir+"/")
		s.legislationLookup[fn] = true
	}

	log.Printf("loaded %d legislation files", len(s.legislationLookup))
	return nil
}

func (s SyncApp) UpdateLegislationSponsors() error {
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
		switch l.StatusName {
		case "Enacted":
			// once things are enacted there shouldn't be sponsor updates
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
			err = s.updateLegislation(ctx, *l)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
