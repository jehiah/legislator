package legistar

import (
	"net/url"
	"time"
)

type EventItem struct {
	ID                int                          `json:"EventItemId"`
	GUID              string                       `json:"EventItemGuid"`
	LastModified      Time                         `json:"EventItemLastModifiedUtc"`
	RowVersion        string                       `json:"EventItemRowVersion"`
	EventID           int                          `json:"EventItemEventId"`
	AgendaSequence    int                          `json:"EventItemAgendaSequence"`
	MinutesSequence   int                          `json:"EventItemMinutesSequence"`
	AgendaNumber      string                       `json:"EventItemAgendaNumber"`
	Video             int                          `json:"EventItemVideo"`
	VideoIndex        int                          `json:"EventItemVideoIndex"`
	Version           string                       `json:"EventItemVersion"`
	AgendaNote        string                       `json:"EventItemAgendaNote"`
	MinutesNote       string                       `json:"EventItemMinutesNote"`
	ActionID          int                          `json:"EventItemActionId"`
	ActionName        string                       `json:"EventItemActionName"`
	ActionText        string                       `json:"EventItemActionText"`
	PassedFlag        int                          `json:"EventItemPassedFlag"`
	PassedFlagName    string                       `json:"EventItemPassedFlagName"`
	RollCallFlag      int                          `json:"EventItemRollCallFlag"`
	FlagExtra         int                          `json:"EventItemFlagExtra"`
	Title             string                       `json:"EventItemTitle"`
	Tally             string                       `json:"EventItemTally"`
	AccelaRecordID    string                       `json:"EventItemAccelaRecordId"`
	Consent           int                          `json:"EventItemConsent"`
	MoverID           int                          `json:"EventItemMoverId"`
	Mover             string                       `json:"EventItemMover"`
	SeconderID        int                          `json:"EventItemSeconderId"`
	Seconder          string                       `json:"EventItemSeconder"`
	MatterID          int                          `json:"EventItemMatterId"`
	MatterGUID        string                       `json:"EventItemMatterGuid"`
	MatterFile        string                       `json:"EventItemMatterFile"`
	MatterName        string                       `json:"EventItemMatterName"`
	MatterType        string                       `json:"EventItemMatterType"`
	MatterStatus      string                       `json:"EventItemMatterStatus"`
	MatterAttachments []EventItemMatterAttachments `json:"EventItemMatterAttachments"`
}
type EventItems []EventItem

type EventItemMatterAttachments struct {
	ID                   int    `json:"MatterAttachmentId"`
	GUID                 string `json:"MatterAttachmentGuid"`
	LastModified         Time   `json:"MatterAttachmentLastModifiedUtc"`
	RowVersion           string `json:"MatterAttachmentRowVersion"`
	Name                 string `json:"MatterAttachmentName"`
	Hyperlink            string `json:"MatterAttachmentHyperlink"`
	FileName             string `json:"MatterAttachmentFileName"`
	MatterVersion        string `json:"MatterAttachmentMatterVersion"`
	IsHyperlink          bool   `json:"MatterAttachmentIsHyperlink"`
	Binary               string `json:"MatterAttachmentBinary"`
	IsSupportingDocument bool   `json:"MatterAttachmentIsSupportingDocument"`
	ShowOnInternetPage   bool   `json:"MatterAttachmentShowOnInternetPage"`
	IsMinuteOrder        bool   `json:"MatterAttachmentIsMinuteOrder"`
	IsBoardLetter        bool   `json:"MatterAttachmentIsBoardLetter"`
	AgiloftID            int    `json:"MatterAttachmentAgiloftId"`
	Description          string `json:"MatterAttachmentDescription"`
	PrintWithReports     bool   `json:"MatterAttachmentPrintWithReports"`
	Sort                 int    `json:"MatterAttachmentSort"`
}

// http://webapi.legistar.com/Help/Api/GET-v1-Client-Events
type Event struct {
	ID                   int        `json:"EventId"`
	GUID                 string     `json:"EventGuid"`
	LastModified         Time       `json:"EventLastModifiedUtc"`
	RowVersion           string     `json:"EventRowVersion"`
	BodyID               int        `json:"EventBodyId"`
	BodyName             string     `json:"EventBodyName"`
	Date                 Time       `json:"EventDate"`
	Time                 ShortTime  `json:"EventTime"`
	VideoStatus          string     `json:"EventVideoStatus"`
	AgendaStatusID       int        `json:"EventAgendaStatusId"`
	AgendaStatusName     string     `json:"EventAgendaStatusName"`
	MinutesStatusID      int        `json:"EventMinutesStatusId"`
	MinutesStatusName    string     `json:"EventMinutesStatusName"`
	Location             string     `json:"EventLocation"`
	AgendaFile           string     `json:"EventAgendaFile"`
	MinutesFile          string     `json:"EventMinutesFile"`
	AgendaLastPublished  Time       `json:"EventAgendaLastPublishedUTC"`
	MinutesLastPublished Time       `json:"EventMinutesLastPublishedUTC"`
	Comment              string     `json:"EventComment"`
	VideoPath            string     `json:"EventVideoPath"`
	Media                string     `json:"EventMedia"`
	InSiteURL            string     `json:"EventInSiteURL"`
	Items                EventItems `json:"EventItems"`
}
type Events []Event

type EventDateFilter struct {
	Direction string
	time.Time
}

func (p EventDateFilter) Paramters() url.Values {
	return DateTimeFilter("EventDate", p.Direction, p.Time)
}

type EventLastModifiedFilter time.Time

func (p EventLastModifiedFilter) Paramters() url.Values {
	return DateTimeFilter("EventLastModifiedUtc", "gt", time.Time(p))
}

// http://webapi.legistar.com/Help/Api/GET-v1-Client-EventItems-EventItemId-RollCalls
type RollCall struct {
	ID           int    `json:"RollCallId"`
	GUID         string `json:"RollCallGuid"`
	LastModified Time   `json:"RollCallLastModifiedUtc"`
	RowVersion   string `json:"RollCallRowVersion"`
	PersonID     int    `json:"RollCallPersonId"`
	PersonName   string `json:"RollCallPersonName"`
	ValueID      int    `json:"RollCallValueId"`
	ValueName    string `json:"RollCallValueName"`
	Sort         int    `json:"RollCallSort"`
	Result       int    `json:"RollCallResult"`
	EventItemID  int    `json:"RollCallEventItemId"`
}
type RollCalls []RollCall
