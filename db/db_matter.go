package db

import (
	"time"

	"github.com/jehiah/legislator/legistar"
)

// NOTE: Struct ordering affects JSON format

type Legislation struct {
	ID    int
	GUID  string
	File  string
	Name  string
	Title string

	TypeID        int
	TypeName      string
	StatusID      int
	StatusName    string
	BodyID        int
	BodyName      string
	IntroDate     time.Time
	AgendaDate    time.Time
	PassedDate    time.Time
	EnactmentDate time.Time
	Version       string `json:",omitempty"`

	Sponsors []PersonReference `json:",omitempty"`

	Summary string // MatterEXText5

	TextID int    `json:",omitempty"`
	Text   string `json:",omitempty"` // current version
	RTF    string `json:",omitempty"` // current version

	LastModified time.Time
}

// Shallow returns a copy of l with details about the text and sponsors ommitted
func (l Legislation) Shallow() Legislation {
	l.Version = ""
	l.Text = ""
	l.TextID = 0
	l.RTF = ""
	l.Sponsors = nil
	return l
}

func NewLegislation(m legistar.Matter) Legislation {
	l := Legislation{
		ID:    m.ID,
		GUID:  m.GUID,
		File:  m.File,
		Name:  m.Name,
		Title: m.Title,

		TypeID:        m.TypeID,
		TypeName:      m.TypeName,
		StatusID:      m.StatusID,
		StatusName:    m.StatusName,
		BodyID:        m.BodyID,
		BodyName:      m.BodyName,
		IntroDate:     m.IntroDate.Time,
		AgendaDate:    m.AgendaDate.Time,
		PassedDate:    m.PassedDate.Time,
		EnactmentDate: m.EnactmentDate.Time,
		Version:       m.Version,

		Summary: m.EXText5,

		LastModified: m.LastModified.Time,
	}
	return l
}
