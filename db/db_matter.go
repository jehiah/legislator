package db

import (
	"time"

	"github.com/jehiah/legislator/legistar"
)

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
	Version       string

	Sponsors []PersonReference

	Summary string // MatterEXText5

	Text string // current Version

	LastModified time.Time
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
