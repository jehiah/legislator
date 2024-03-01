package main

import (
	"fmt"
	"log"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/jehiah/legislator/db"
)

type ResubmitFile struct {
	Resubmitted []db.ResubmitLegislation
}

func (f ResubmitFile) Has(file string) bool {
	for _, r := range f.Resubmitted {
		if r.ToFile == file {
			return true
		}
	}
	return false
}
func (f *ResubmitFile) Add(r db.ResubmitLegislation) bool {
	if f.Has(r.ToFile) {
		return false
	}
	f.Resubmitted = append(f.Resubmitted, r)
	return true
}

func (s *SyncApp) loadResubmit() (map[string]*ResubmitFile, error) {
	data := make(map[string]*ResubmitFile)
	files, err := filepath.Glob(filepath.Join(s.targetDir, "resubmit", "*.json"))
	if err != nil {
		return nil, err
	}
	for _, fn := range files {
		fn = strings.TrimPrefix(fn, s.targetDir+"/")
		// log.Printf("loading %s", fn)

		var f ResubmitFile
		err := s.readFile(fn, &f)
		if err != nil {
			return nil, err
		}

		data[fn] = &f
	}

	log.Printf("loaded %d resubmit files", len(data))
	return data, nil
}

type Session struct {
	StartYear, EndYear int // inclusive
}

func (s Session) ContainsYear(n int) bool {
	return n >= s.StartYear && n <= s.EndYear
}

func (s SyncApp) Sessions() []Session {

	var years []int
	for fn := range s.legislationLookup {
		if !strings.HasSuffix(fn, "/0001.json") {
			continue
		}
		years = append(years, yearFromFile(fn))
	}
	sort.Ints(years)
	var sessions []Session
	for i, y := range years {
		if i+1 > len(years) {
			sessions = append(sessions, Session{y, years[i+1]})
		} else {
			sessions = append(sessions, Session{y, y + 1}) // could be 4 years 🤷‍♂️
		}
	}
	return sessions
}

func yearFromFile(s string) int {
	year, _ := strconv.Atoi(strings.SplitN(s, "/", 3)[1])
	return year
}

// SyncResubmit detects legislation that's resubmitted based on the Name or Title
func (s SyncApp) SyncResubmit() error {
	resubmitFilename := func(n int) string {
		return fmt.Sprintf("resubmit/%d.json", n)
	}
	data, err := s.loadResubmit()
	if err != nil {
		return err
	}

	resubmitName := make(map[string]string, len(s.legislationLookup))
	resubmitTitle := make(map[string]string, len(s.legislationLookup))
	resbumitSummary := make(map[string]string, len(s.legislationLookup))

	sessions := s.Sessions()
	numberAdded := 0

	for i, session := range sessions {
		if i == 0 {
			continue
		}
		lastSession := sessions[i-1]

		// build the lookups
		for fn := range s.legislationLookup {
			if !lastSession.ContainsYear(yearFromFile(fn)) {
				continue
			}
			var l *db.Legislation
			err = s.readFile(fn, &l)
			if err != nil {
				return err
			}
			switch l.StatusName {
			case "Filed", "Filed (End of Session)":
			default:
				continue
			}

			resubmitName[l.Name] = l.File
			resubmitTitle[l.Title] = l.File
			if l.Summary != "" {
				resbumitSummary[l.Summary] = l.File
			}
		}

		for fn := range s.legislationLookup {
			year := yearFromFile(fn)
			if !session.ContainsYear(year) {
				continue
			}
			if _, ok := data[resubmitFilename(year)]; !ok {
				data[resubmitFilename(year)] = &ResubmitFile{}
			}
			var l *db.Legislation
			err := s.readFile(fn, &l)
			if err != nil {
				return err
			}
			switch l.StatusName {
			case "Withdrawn":
				continue
			}
			old := resubmitName[l.Name]
			if old == "" {
				old = resubmitTitle[l.Title]
			}
			if old == "" && l.Summary != "" {
				old = resbumitSummary[l.Summary]
			}
			if old != "" {
				// log.Printf("resubmit %q => %q %s", old, l.File, resubmitFilename(l.IntroDate.Year()))
				added := data[resubmitFilename(l.IntroDate.Year())].Add(db.ResubmitLegislation{
					FromFile: old,
					ToFile:   l.File,
				})
				if added {
					numberAdded++
				}
			}
		}
	}

	log.Printf("detected %d new resubmissions", numberAdded)

	for fn, f := range data {
		if len(f.Resubmitted) == 0 {
			continue
		}
		sort.Slice(f.Resubmitted, func(i, j int) bool { return strings.Compare(f.Resubmitted[i].ToFile, f.Resubmitted[j].ToFile) == -1 })
		err = s.writeFile(fn, f)
		if err != nil {
			return err
		}
	}

	return nil
}
