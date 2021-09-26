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
	return nil
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
	People time.Time

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
	os.MkdirAll(*targetDir, 0777)
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
