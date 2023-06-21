package db

import (
	"strings"
	"time"

	"github.com/jehiah/legislator/legistar"
)

// NOTE: Struct ordering affects JSON format

type Legislation struct {
	ID       int
	GUID     string
	File     string
	LocalLaw string `json:",omitempty"` // format: $year/$number
	Name     string
	Title    string

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

	Sponsors    []PersonReference `json:",omitempty"`
	History     []History         `json:",omitempty"`
	Attachments []Attachment      `json:",omitempty"`

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
		ID:       m.ID,
		GUID:     m.GUID,
		File:     m.File,
		LocalLaw: m.EnactmentNumber,
		Name:     m.Name,
		Title:    m.Title,

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

type History struct {
	ID int
	// GUID string

	Date        time.Time
	ActionID    int
	Action      string
	Description string
	BodyID      int
	BodyName    string

	EventID         int    `json:",omitempty"`
	AgendaSequence  int    `json:",omitempty"`
	MinutesSequence int    `json:",omitempty"`
	AgendaNumber    string `json:",omitempty"`
	Video           int    `json:",omitempty"`
	VideoIndex      int    `json:",omitempty"`
	Version         string
	AgendaNote      string `json:",omitempty"`
	MinutesNote     string `json:",omitempty"`

	PassedFlag     int    `json:",omitempty"` // Note: 0 may indicate a vote that failed
	PassedFlagName string `json:",omitempty"` // Pass, Fail
	RollCallFlag   int    `json:",omitempty"`
	FlagExtra      int    `json:",omitempty"`
	Tally          string `json:",omitempty"`
	AccelaRecordID string `json:",omitempty"`
	Consent        int    `json:",omitempty"`
	MoverID        int    `json:",omitempty"`
	MoverName      string `json:",omitempty"`
	SeconderID     int    `json:",omitempty"`
	SeconderName   string `json:",omitempty"`
	MatterStatusID int

	Votes []Vote `json:",omitempty"`

	LastModified time.Time
}

func NewHistory(h legistar.MatterHistory) History {
	return History{
		ID: h.ID,
		// GUID: h.GUID,

		Date:        h.ActionDate.Time,
		ActionID:    h.ActionID,
		Action:      h.ActionName,
		Description: h.ActionText,
		BodyID:      h.ActionBodyID,
		BodyName:    h.ActionBodyName,

		EventID:         h.EventID,
		AgendaSequence:  h.AgendaSequence,
		MinutesSequence: h.MinutesSequence,
		AgendaNumber:    h.AgendaNumber,
		Video:           h.Video,
		VideoIndex:      h.VideoIndex,
		Version:         h.Version,
		AgendaNote:      h.AgendaNote,
		MinutesNote:     h.MinutesNote,

		PassedFlag:     h.PassedFlag,
		PassedFlagName: h.PassedFlagName,
		RollCallFlag:   h.RollCallFlag,
		FlagExtra:      h.FlagExtra,
		Tally:          h.Tally,
		AccelaRecordID: h.AccelaRecordID,
		Consent:        h.Consent,
		MoverID:        h.MoverID,
		MoverName:      h.MoverName,
		SeconderID:     h.SeconderID,
		SeconderName:   h.SeconderName,
		MatterStatusID: h.MatterStatusID,

		LastModified: h.LastModified.Time,
	}
}

type Attachment struct {
	ID           int
	LastModified time.Time
	Name         string
	Link         string
	// MatterVersion        string
	// IsSupportingDocument bool
	// ShowOnInternetPage   bool
	// IsMinuteOrder        bool
	// IsBoardLetter        bool
	// AgiloftID            int
	// Description          string
	// PrintWithReports     bool
	Sort int
}

func NewAttachment(h legistar.MatterAttachment) Attachment {
	return Attachment{
		ID:           h.ID,
		LastModified: h.LastModified.Time,
		Name:         h.Name,
		Link:         h.Link,
		Sort:         h.Sort,
	}
}

type Vote struct {
	PersonReference
	VoteID int
	Vote   string // Afffirmative, Negative, Absent
	Result int    `json:",omitempty"` // 1 = Afirmative, 2 = Negative
	Sort   int
}
type Votes []Vote

func NewVotes(v legistar.Votes) Votes {
	var o Votes
	for _, vv := range v {
		o = append(o, NewVote(vv))
	}
	return o
}

func NewVote(v legistar.Vote) Vote {
	return Vote{
		PersonReference: PersonReference{
			FullName: strings.TrimSpace(v.PersonName),
			ID:       v.PersonID,
			Slug:     v.Slug(),
		},
		VoteID: v.ValueID,
		Vote:   v.ValueName,
		Result: v.Result,
		Sort:   v.Sort,
	}
}
