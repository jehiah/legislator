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

func LandUseFilename(m db.Legislation) string {
	fn := strings.Fields(strings.ReplaceAll(m.File, "-", " "))[1] + ".json"
	return filepath.Join("land_use", strconv.Itoa(m.IntroDate.Year()), fn)
}

func (s *SyncApp) SyncLandUse(filter legistar.Filters) error {
	ctx := context.Background()
	if filter == nil {
		filter = legistar.AndFilters(
			legistar.MatterLastModifiedFilter(s.LastSync.LandUse),
			legistar.MatterTypeFilter("Land Use Application"),
		)
	} else {
		filter = legistar.AndFilters(
			filter,
			legistar.MatterTypeFilter("Land Use Application"),
		)
	}

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
		fn := LandUseFilename(l)
		s.legislationLookup[fn] = true
		err = s.updateLandUse(ctx, l)
		if err != nil {
			return err
		}
	}

	// // check if any have new sponsors
	// if err = s.UpdateMatterSponsors(); err != nil {
	// 	return err
	// }

	if len(matters) > 0 {
		s.LastSync.LandUse = db.Max(matters, func(i int) time.Time { return matters[i].LastModified.Time })
	}
	return nil
}

// UpdateMatterByFile expects format 1234-2020
func (s *SyncApp) UpdateLandUseByFile(q string) error {
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
	return s.UpdateLandUse(ctx, matters[0].ID)
}

func (s *SyncApp) UpdateAllLandUse() error {
	ctx := context.Background()
	for fn := range s.landUseLookup {
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

		err = s.UpdateLandUseWithRetry(ctx, l.ID)
		if err != nil {
			log.Printf("Got error after retry; sleeping 30s and skipping %s %s", l.File, err)
			time.Sleep(time.Second * 30)
			// return err
		}

	}
	return nil
}

func (s *SyncApp) UpdateLandUseWithRetry(ctx context.Context, ID int) error {
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
	err = s.updateLandUse(ctx, l)
	if err != nil {
		log.Print(err)
		time.Sleep(time.Second)
		err = s.updateLandUse(ctx, l)
	}
	return err
}

func (s *SyncApp) UpdateLandUse(ctx context.Context, ID int) error {
	m, err := s.legistar.Matter(ctx, ID)
	if err != nil {
		return err
	}
	l := db.NewLegislation(m)
	return s.updateLandUse(ctx, l)
}

func (s *SyncApp) updateLandUse(ctx context.Context, l db.Legislation) error {
	fn := LandUseFilename(l)
	return s.updateMatter(ctx, fn, l)
}

func (s *SyncApp) LoadLandUse() error {
	files, err := filepath.Glob(filepath.Join(s.targetDir, "land_use", "*", "*.json"))
	if err != nil {
		return err
	}
	for _, fn := range files {
		if filepath.Base(fn) == "index.json" {
			continue
		}
		fn = strings.TrimPrefix(fn, s.targetDir+"/")
		s.landUseLookup[fn] = true
	}

	log.Printf("loaded %d land use files", len(s.landUseLookup))
	return nil
}
