package main

import (
	"encoding/json"
	"flag"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/jehiah/legislator/db"
	"github.com/jehiah/legislator/legistar"
)

type SyncApp struct {
	legistar  *legistar.Client
	targetDir string

	personLookup      map[int]db.Person
	legislationLookup map[string]bool

	LastSync
}

type LastSync struct {
	Matters time.Time
	Persons time.Time

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
	err = s.LoadMatter()
	if err != nil {
		return err
	}
	return nil
}

func (s *SyncApp) Run() error {
	os.MkdirAll(s.targetDir, 0777)
	os.MkdirAll(filepath.Join(s.targetDir, "people"), 0777)
	os.MkdirAll(filepath.Join(s.targetDir, "introduction"), 0777)
	s.LastRun = time.Now().UTC().Truncate(time.Second)
	err := s.SyncPersons()
	if err != nil {
		return err
	}
	err = s.SyncMatter()
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
	err := s.CreateIndexes()
	if err != nil {
		return err
	}

	fn := filepath.Join(s.targetDir, "last_sync.json")
	f, err := os.Create(fn)
	if err != nil {
		return err
	}
	defer f.Close()
	b, err := json.Marshal(s.LastSync)
	if err != nil {
		return err
	}
	f.Write(b)
	return f.Close()
}

func main() {
	targetDir := flag.String("target-dir", "", "Target Directory")
	update := flag.String("update-one", "", "update one")
	flag.Parse()
	if *targetDir == "" {
		log.Fatal("set --target-dir")
	}

	s := &SyncApp{
		legistar:          legistar.NewClient("nyc", os.Getenv("NYC_LEGISLATOR_TOKEN")),
		personLookup:      make(map[int]db.Person),
		legislationLookup: make(map[string]bool),
		targetDir:         *targetDir,
	}

	if err := s.Load(); err != nil {
		log.Fatal(err)
	}
	if *update != "" {
		if err := s.UpdateOne(*update); err != nil {
			log.Fatal(err)
		}
	} else {
		if err := s.Run(); err != nil {
			log.Fatal(err)
		}
	}
	if err := s.Save(); err != nil {
		log.Fatal(err)
	}
}
