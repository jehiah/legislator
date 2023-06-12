package db

import (
	"strings"
	"time"

	"github.com/jehiah/legislator/legistar"
)

type Event struct {
	ID   int
	GUID string

	BodyID   int
	BodyName string
	Date     time.Time

	Location          string
	VideoStatus       string
	AgendaStatusID    int
	AgendaStatusName  string
	MinutesStatusID   int
	MinutesStatusName string
	AgendaFile        string
	MinutesFile       string `json:",omitempty"`
	Comment           string `json:",omitempty"`
	VideoPath         string `json:",omitempty"`
	Media             string `json:",omitempty"`
	InSiteURL         string
	Items             []EventItem `json:",omitempty"`

	AgendaLastPublished  time.Time
	MinutesLastPublished time.Time
	LastModified         time.Time
}

func NewEvent(m legistar.Event, tz *time.Location) Event {
	e := Event{
		ID:   m.ID,
		GUID: m.GUID,

		BodyID:   m.BodyID,
		BodyName: m.BodyName,
		Date:     m.Time.Set(m.Date.Time, tz),

		Location:          m.Location,
		VideoStatus:       m.VideoStatus,
		AgendaStatusID:    m.AgendaStatusID,
		AgendaStatusName:  m.AgendaStatusName,
		MinutesStatusID:   m.MinutesStatusID,
		MinutesStatusName: m.MinutesStatusName,
		AgendaFile:        m.AgendaFile,
		MinutesFile:       m.MinutesFile,
		Comment:           m.Comment,
		VideoPath:         m.VideoPath,
		Media:             m.Media,
		InSiteURL:         m.InSiteURL,

		AgendaLastPublished:  m.AgendaLastPublished.Time,
		MinutesLastPublished: m.MinutesLastPublished.Time,
		LastModified:         m.LastModified.Time,
		Items:                NewEventItems(m.Items),
	}
	return e
}

type EventItem struct {
	ID   int
	GUID string

	Title string

	// EventID           int
	AgendaSequence    int
	MinutesSequence   int
	AgendaNumber      string `json:",omitempty"`
	Video             int    `json:",omitempty"`
	VideoIndex        int    `json:",omitempty"`
	Version           string
	AgendaNote        string                       `json:",omitempty"`
	MinutesNote       string                       `json:",omitempty"`
	ActionID          int                          `json:",omitempty"`
	ActionName        string                       `json:",omitempty"`
	ActionText        string                       `json:",omitempty"`
	PassedFlag        int                          `json:",omitempty"`
	PassedFlagName    string                       `json:",omitempty"`
	RollCallFlag      int                          `json:",omitempty"`
	FlagExtra         int                          `json:",omitempty"`
	Tally             string                       `json:",omitempty"`
	AccelaRecordID    string                       `json:",omitempty"`
	Consent           int                          `json:",omitempty"`
	MoverID           int                          `json:",omitempty"`
	Mover             string                       `json:",omitempty"`
	SeconderID        int                          `json:",omitempty"`
	Seconder          string                       `json:",omitempty"`
	MatterID          int                          `json:",omitempty"`
	MatterFile        string                       `json:",omitempty"`
	MatterName        string                       `json:",omitempty"`
	MatterType        string                       `json:",omitempty"`
	MatterStatus      string                       `json:",omitempty"`
	MatterAttachments []EventItemMatterAttachments `json:",omitempty"`

	LastModified time.Time
}
type EventItems []EventItem

func NewEventItems(e legistar.EventItems) EventItems {
	var o EventItems
	for _, ee := range e {
		o = append(o, NewEventItem(ee))
	}
	return o
}

func NewEventItem(e legistar.EventItem) EventItem {

	return EventItem{
		ID:   e.ID,
		GUID: e.GUID,
		// EventID        : e.// EventID        ,
		AgendaSequence:  e.AgendaSequence,
		MinutesSequence: e.MinutesSequence,
		AgendaNumber:    e.AgendaNumber,
		Video:           e.Video,
		VideoIndex:      e.VideoIndex,
		Version:         e.Version,
		AgendaNote:      e.AgendaNote,
		MinutesNote:     e.MinutesNote,
		ActionID:        e.ActionID,
		ActionName:      e.ActionName,
		ActionText:      e.ActionText,
		PassedFlag:      e.PassedFlag,
		PassedFlagName:  e.PassedFlagName,
		RollCallFlag:    e.RollCallFlag,
		FlagExtra:       e.FlagExtra,
		Title:           e.Title,
		Tally:           e.Tally,
		AccelaRecordID:  e.AccelaRecordID,
		Consent:         e.Consent,
		MoverID:         e.MoverID,
		Mover:           e.Mover,
		SeconderID:      e.SeconderID,
		Seconder:        e.Seconder,
		MatterID:        e.MatterID,
		MatterFile:      e.MatterFile,
		MatterName:      e.MatterName,
		MatterType:      e.MatterType,
		MatterStatus:    e.MatterStatus,
		// MatterAttachments : e.MatterAttachments ,

		LastModified: e.LastModified.Time,
	}
}

type EventItemMatterAttachments struct {
	ID                   int
	GUID                 string
	Name                 string
	Hyperlink            string
	FileName             string
	MatterVersion        string
	IsHyperlink          bool
	Binary               string
	IsSupportingDocument bool
	ShowOnInternetPage   bool
	IsMinuteOrder        bool
	IsBoardLetter        bool
	AgiloftID            int
	Description          string
	PrintWithReports     bool
	Sort                 int
	LastModified         time.Time
}

type RollCall struct {
	ID           int
	GUID         string
	LastModified time.Time
	Person       PersonReference

	ValueID   int    `json:",omitempty"`
	ValueName string `json:",omitempty"`
	Sort      int    `json:",omitempty"`
	Result    int    `json:",omitempty"`
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
