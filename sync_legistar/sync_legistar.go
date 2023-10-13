package main

import (
	"context"
	"encoding/json"
	"flag"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/jehiah/legislator/db"
	"github.com/jehiah/legislator/legistar"
)

var localTimezone *time.Location

type SyncApp struct {
	legistar  *legistar.Client
	targetDir string

	personLookup      map[int]db.Person
	legislationLookup map[string]bool
	landUseLookup     map[string]bool
	resolutionLookkup map[string]bool
	eventsLookup      map[int][]string

	LastSync
}

type LastSync struct {
	Matters    time.Time
	Persons    time.Time
	Events     time.Time
	LandUse    time.Time
	Resolution time.Time

	LastRun time.Time
}

func (s *SyncApp) Load() error {
	fn := filepath.Join(s.targetDir, "last_sync.json")
	_, err := os.Stat(fn)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return err
	}
	b, err := os.ReadFile(fn)
	if err != nil {
		return err
	}
	err = json.Unmarshal(b, &s.LastSync)
	if err != nil {
		return err
	}
	err = s.LoadPersons()
	if err != nil {
		return err
	}
	err = s.LoadLegislation()
	if err != nil {
		return err
	}
	err = s.LoadEvents()
	if err != nil {
		return err
	}
	err = s.LoadLandUse()
	if err != nil {
		return err
	}
	err = s.LoadResolution()
	if err != nil {
		return err
	}
	return nil
}

func (s *SyncApp) Run() error {
	os.MkdirAll(s.targetDir, 0777)
	os.MkdirAll(filepath.Join(s.targetDir, "people"), 0777)
	os.MkdirAll(filepath.Join(s.targetDir, "introduction"), 0777)
	os.MkdirAll(filepath.Join(s.targetDir, "land_use"), 0777)
	os.MkdirAll(filepath.Join(s.targetDir, "resolution"), 0777)
	os.MkdirAll(filepath.Join(s.targetDir, "events"), 0777)
	s.LastRun = time.Now().UTC().Truncate(time.Second)
	err := s.SyncPersons()
	if err != nil {
		return err
	}
	err = s.SyncLegislation()
	if err != nil {
		return err
	}

	err = s.SyncEvents(nil)
	if err != nil {
		return err
	}
	err = s.SyncLandUse(nil)
	if err != nil {
		return err
	}
	err = s.SyncResolution(nil)
	if err != nil {
		return err
	}
	return nil
}

func (s SyncApp) openWriteFile(fn string) (*os.File, error) {
	fn = filepath.Join(s.targetDir, fn)
	err := os.MkdirAll(filepath.Dir(fn), 0777)
	if err != nil {
		return nil, err
	}
	log.Printf("creating %s", fn)
	return os.Create(fn)
}

func (s SyncApp) writeFile(fn string, o interface{}) error {
	f, err := s.openWriteFile(fn)
	if err != nil {
		return err
	}
	defer f.Close()

	e := json.NewEncoder(f)
	e.SetEscapeHTML(false)
	e.SetIndent("", "  ")
	err = e.Encode(o)
	if err != nil {
		return err
	}
	return f.Close()
}

func (s SyncApp) removeFile(fn string) error {
	fn = filepath.Join(s.targetDir, fn)
	log.Printf("removing %s", fn)
	err := os.Remove(fn)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}

func (s SyncApp) readFile(fn string, o interface{}) error {
	fn = filepath.Join(s.targetDir, fn)
	body, err := os.ReadFile(fn)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	return json.Unmarshal(body, &o)
}

func (s SyncApp) Save() error {
	fn := filepath.Join(s.targetDir, "last_sync.json")
	f, err := os.Create(fn)
	if err != nil {
		return err
	}
	defer f.Close()
	e := json.NewEncoder(f)
	e.SetEscapeHTML(false)
	e.SetIndent("", "  ")
	err = e.Encode(s.LastSync)
	if err != nil {
		return err
	}
	return f.Close()
}

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	targetDir := flag.String("target-dir", "", "Target Directory")
	updatePeople := flag.Bool("update-people", false, "update all people")
	updateLegislation := flag.String("update-legislation", "", "File of legislation to update i.e. 1234-2020")
	updateEvent := flag.String("update-event", "", "the ID of an event to update")
	timezone := flag.String("tz", "America/New_York", "timezone")
	updateAll := flag.Bool("update-all", false, "update all")
	skipIndexUpdate := flag.Bool("skip-index-update", false, "skip updating year index files and last_sync.json")
	flag.Parse()
	if *targetDir == "" {
		log.Fatal("set --target-dir")
	}

	s := &SyncApp{
		legistar:          legistar.NewClient("nyc", os.Getenv("NYC_LEGISLATOR_TOKEN")),
		personLookup:      make(map[int]db.Person),
		legislationLookup: make(map[string]bool),
		landUseLookup:     make(map[string]bool),
		resolutionLookkup: make(map[string]bool),
		eventsLookup:      make(map[int][]string),
		targetDir:         *targetDir,
	}
	ctx := context.Background()

	var err error
	localTimezone, err = time.LoadLocation(*timezone)
	if err != nil {
		log.Fatal(err)
	}

	if err = s.Load(); err != nil {
		log.Fatal(err)
	}
	switch {
	case *updateLegislation != "":
		err = s.UpdateLegislationByFile(*updateLegislation)
	case *updateEvent != "":
		id, err := strconv.Atoi(*updateEvent)
		if err != nil {
			log.Fatal(err)
		}
		err = s.SyncEvent(ctx, id)
	case *updateAll:
		// err = s.UpdateAllLegislation()
		// err = s.SyncAllEvent()
		// err = s.SyncDuplicateEvents()
		// err = s.SyncRollCalls()

		for year := 2021; year >= 2006; year-- {
			// filter := legistar.AndFilters(
			// 	MatterDateYearFilter{time.Date(year, time.January, 1, 0, 0, 0, 0, time.UTC), "gt"},
			// 	MatterDateYearFilter{time.Date(year+1, time.January, 1, 0, 0, 0, 0, time.UTC), "lt"},
			// )
			// err = s.SyncLandUse(filter)
			// if err != nil {
			// 	break
			// }
			// err = s.SyncResolution(filter)
			// if err != nil {
			// 	break
			// }

			filter := legistar.AndFilters(
				legistar.EventDateFilter{Time: time.Date(year, time.January, 1, 0, 0, 0, 0, time.UTC), Direction: "gt"},
				legistar.EventDateFilter{Time: time.Date(year+1, time.January, 1, 0, 0, 0, 0, time.UTC), Direction: "lt"},
			)
			err = s.SyncEvents(filter)
			if err != nil {
				break
			}

		}

	case *updatePeople:
		err = s.UpdateActive(ctx)
	default:
		err = s.Run()
	}
	if err != nil {
		log.Fatal(err)
	}
	if !*skipIndexUpdate {
		if err := s.Save(); err != nil {
			log.Fatal(err)
		}
	}
}
