package main

import (
	"context"
	"fmt"
	"log"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/jehiah/legislator/db"
	"github.com/jehiah/legislator/legistar"
)

func ResolutionFilename(m db.Legislation) string {
	fn := strings.Fields(strings.ReplaceAll(m.File, "-", " "))[1] + ".json"
	return filepath.Join("resolution", strconv.Itoa(m.IntroDate.Year()), fn)
}

func (s *SyncApp) SyncResolution() error {
	ctx := context.Background()
	filter := legistar.AndFilters(
		// legistar.MatterLastModifiedFilter(s.LastSync.Resolution),
		legistar.MatterTypeFilter("Land Use Application"),
		MatterDateYearFilter{time.Date(2022, time.January, 1, 0, 0, 0, 0, time.UTC), "gt"},
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
		fn := ResolutionFilename(l)
		s.legislationLookup[fn] = true
		err = s.updateResolution(ctx, l)
		if err != nil {
			return err
		}
	}

	// // check if any have new sponsors
	// if err = s.UpdateMatterSponsors(); err != nil {
	// 	return err
	// }

	if len(matters) > 0 {
		s.LastSync.Resolution = db.Max(matters, func(i int) time.Time { return matters[i].LastModified.Time })
	}
	return nil
}

// UpdateMatterByFile expects format 1234-2020
func (s *SyncApp) UpdateResolutionByFile(q string) error {
	ctx := context.Background()
	file := fmt.Sprintf("Int %s", q)
	filter := legistar.AndFilters(
		legistar.MatterTypeFilter("Land Use Application"),
		legistar.MatterFileFilter(file),
	)

	matters, err := s.legistar.Matters(ctx, filter)
	if err != nil {
		return err
	}
	if len(matters) != 1 {
		return fmt.Errorf("expected 1 response got %d for %q", len(matters), q)
	}
	return s.UpdateResolution(ctx, matters[0].ID)
}

func (s *SyncApp) UpdateAllResolution() error {
	ctx := context.Background()
	for fn := range s.resolutionLookkup {
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
		if !update {
			continue
		}

		err = s.UpdateResolutionWithRetry(ctx, l.ID)
		if err != nil {
			log.Printf("Got error after retry; sleeping 30s and skipping %s %s", l.File, err)
			time.Sleep(time.Second * 30)
			// return err
		}

	}
	return nil
}

func (s *SyncApp) UpdateResolutionWithRetry(ctx context.Context, ID int) error {
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
	err = s.updateResolution(ctx, l)
	if err != nil {
		log.Print(err)
		time.Sleep(time.Second)
		err = s.updateResolution(ctx, l)
	}
	return err
}

func (s *SyncApp) UpdateResolution(ctx context.Context, ID int) error {
	m, err := s.legistar.Matter(ctx, ID)
	if err != nil {
		return err
	}
	l := db.NewLegislation(m)
	return s.updateResolution(ctx, l)
}

func (s *SyncApp) updateResolution(ctx context.Context, l db.Legislation) error {
	fn := ResolutionFilename(l)
	return s.updateMatter(ctx, fn, l)
}

func (s *SyncApp) LoadResolution() error {
	files, err := filepath.Glob(filepath.Join(s.targetDir, "resolution", "*", "*.json"))
	if err != nil {
		return err
	}
	for _, fn := range files {
		if filepath.Base(fn) == "index.json" {
			continue
		}
		fn = strings.TrimPrefix(fn, s.targetDir+"/")
		s.resolutionLookkup[fn] = true
	}

	log.Printf("loaded %d resolution files", len(s.resolutionLookkup))
	return nil
}
