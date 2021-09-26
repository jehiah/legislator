package main

import (
	"encoding/json"
	"flag"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/jehiah/legislator/legistar"
)

type SyncApp struct {
	legistar  *legistar.Client
	targetDir string

	LastSync
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
	return json.Unmarshal(b, &s.LastSync)
}

func (s *SyncApp) Run() error {
	os.MkdirAll(s.targetDir, 0777)
	os.MkdirAll(filepath.Join(s.targetDir, "people"), 0777)
	s.LastRun = time.Now().UTC().Truncate(time.Second)
	err := s.SyncPersons()
	if err != nil {
		return err
	}
	return nil
}

func (s SyncApp) writeFile(fn string, o interface{}) error {
	b, err := json.Marshal(o)
	if err != nil {
		return err
	}
	fn = filepath.Join(s.targetDir, fn)
	log.Printf("creating %s", fn)
	f, err := os.Create(fn)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.Write(b)
	if err != nil {
		return err
	}
	return f.Close()
}

func (s SyncApp) Save() error {
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

type LastSync struct {
	// Matter time.Time
	Persons time.Time

	LastRun time.Time
}

func main() {
	targetDir := flag.String("target-dir", "", "Target Directory")
	flag.Parse()
	if *targetDir == "" {
		log.Fatal("set --target-dir")
	}

	s := &SyncApp{
		legistar: &legistar.Client{
			Client: "nyc",
			Token:  os.Getenv("NYC_LEGISLATOR_TOKEN"),
		},
		targetDir: *targetDir,
	}
	if err := s.Load(); err != nil {
		log.Fatal(err)
	}
	if err := s.Run(); err != nil {
		log.Fatal(err)
	}
	if err := s.Save(); err != nil {
		log.Fatal(err)
	}
}
