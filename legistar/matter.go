package legistar

import (
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/gosimple/slug"
)

// Matter
// http://webapi.legistar.com/Help/Api/GET-v1-Client-Matters
type Matter struct {
	ID                 int    `json:"MatterId"`
	GUID               string `json:"MatterGuid"`
	LastModified       Time   `json:"MatterLastModifiedUtc"`
	RowVersion         string `json:"MatterRowVersion"`
	File               string `json:"MatterFile"`
	Name               string `json:"MatterName"`
	Title              string `json:"MatterTitle"`
	TypeID             int    `json:"MatterTypeId"`
	TypeName           string `json:"MatterTypeName"`
	StatusID           int    `json:"MatterStatusId"`
	StatusName         string `json:"MatterStatusName"`
	BodyID             int    `json:"MatterBodyId"`
	BodyName           string `json:"MatterBodyName"`
	IntroDate          Time   `json:"MatterIntroDate"`
	AgendaDate         Time   `json:"MatterAgendaDate"`
	PassedDate         Time   `json:"MatterPassedDate"`
	EnactmentDate      Time   `json:"MatterEnactmentDate"`
	EnactmentNumber    string `json:"MatterEnactmentNumber"` // aka Local Law
	Requester          string `json:"MatterRequester"`
	Notes              string `json:"MatterNotes"`
	Version            string `json:"MatterVersion"`
	Text1              string `json:"MatterText1"`
	Text2              string `json:"MatterText2"`
	Text3              string `json:"MatterText3"`
	Text4              string `json:"MatterText4"`
	Text5              string `json:"MatterText5"`
	Date1              Time   `json:"MatterDate1"`
	Date2              Time   `json:"MatterDate2"`
	EXText1            string `json:"MatterEXText1"`
	EXText2            string `json:"MatterEXText2"`
	EXText3            string `json:"MatterEXText3"`
	EXText4            string `json:"MatterEXText4"`
	EXText5            string `json:"MatterEXText5"`
	EXText6            string `json:"MatterEXText6"`
	EXText7            string `json:"MatterEXText7"`
	EXText8            string `json:"MatterEXText8"`
	EXText9            string `json:"MatterEXText9"`
	EXText10           string `json:"MatterEXText10"`
	EXText11           string `json:"MatterEXText11"`
	EXDate1            Time   `json:"MatterEXDate1"`
	EXDate2            Time   `json:"MatterEXDate2"`
	EXDate3            Time   `json:"MatterEXDate3"`
	EXDate4            Time   `json:"MatterEXDate4"`
	EXDate5            Time   `json:"MatterEXDate5"`
	EXDate6            Time   `json:"MatterEXDate6"`
	EXDate7            Time   `json:"MatterEXDate7"`
	EXDate8            Time   `json:"MatterEXDate8"`
	EXDate9            Time   `json:"MatterEXDate9"`
	EXDate10           Time   `json:"MatterEXDate10"`
	AgiloftID          int    `json:"MatterAgiloftId"`
	Reference          string `json:"MatterReference"`
	RestrictViewViaWeb bool   `json:"MatterRestrictViewViaWeb"`
	Reports            []struct {
		ReportName string `json:"ReportName"`
		ReportURL  string `json:"ReportURL"`
		ReportType string `json:"ReportType"`
	} `json:"MatterReports"`
}

type Matters []Matter

type MatterLastModifiedFilter time.Time

func (p MatterLastModifiedFilter) Paramters() url.Values {
	return DateTimeFilter("MatterLastModifiedUtc", "gt", time.Time(p))
}

type MatterTypeFilter string

func (p MatterTypeFilter) Paramters() url.Values {
	return StringFilter("MatterTypeName", string(p))
}

type MatterFileFilter string

func (p MatterFileFilter) Paramters() url.Values {
	return StringFilter("MatterFile", string(p))
}

type MatterEnactmentNumberFilter string

func (p MatterEnactmentNumberFilter) Paramters() url.Values {
	return StringFilter("MatterEnactmentNumber", string(p))
}

// MatterSponsor
// http://webapi.legistar.com/Help/Api/GET-v1-Client-Matters-MatterId-Sponsors
type MatterSponsor struct {
	ID            int    `json:"MatterSponsorId"`
	GUID          string `json:"MatterSponsorGuid"`
	LastModified  Time   `json:"MatterSponsorLastModifiedUtc"`
	RowVersion    string `json:"MatterSponsorRowVersion"`
	MatterID      int    `json:"MatterSponsorMatterId"`
	MatterVersion string `json:"MatterSponsorMatterVersion"`
	Name          string `json:"MatterSponsorName"`
	NameID        int    `json:"MatterSponsorNameId"`
	Sequence      int    `json:"MatterSponsorSequence"`
	BodyID        int    `json:"MatterSponsorBodyId"`
	LinkFlag      int    `json:"MatterSponsorLinkFlag"`
}

func (s MatterSponsor) Slug() string {
	return slug.MakeLang(s.Name, "en")
}

type MatterSponsors []MatterSponsor

// MatterType
// http://webapi.legistar.com/Help/Api/GET-v1-Client-MatterTypes
type MatterType struct {
	ID           int    `json:"MatterTypeId"`
	GUID         string `json:"MatterTypeGuid"`
	LastModified Time   `json:"MatterTypeLastModifiedUtc"`
	RowVersion   string `json:"MatterTypeRowVersion"`
	Name         string `json:"MatterTypeName"`
	Sort         int    `json:"MatterTypeSort"`
	ActiveFlag   int    `json:"MatterTypeActiveFlag"`
	Description  string `json:"MatterTypeDescription"`
	UsedFlag     int    `json:"MatterTypeUsedFlag"`
}
type MatterTypes []MatterType

// MatterIndex
// http://webapi.legistar.com/Help/Api/GET-v1-Client-MatterIndexes
type MatterIndex struct {
	ID           int    `json:"MatterIndexId"`
	GUID         string `json:"MatterIndexGuid"`
	LastModified Time   `json:"MatterIndexLastModifiedUtc"`
	RowVersion   string `json:"MatterIndexRowVersion"`
	MatterID     int    `json:"MatterIndexMatterId"`
	IndexID      int    `json:"MatterIndexIndexId"`
	Name         string `json:"MatterIndexName"`
}
type MatterIndexes []MatterIndex

// MatterText
// http://webapi.legistar.com/Help/Api/GET-v1-Client-Matters-MatterId-Texts-MatterTextId
type MatterText struct {
	ID           int    `json:"MatterTextId"`
	GUID         string `json:"MatterTextGuid"`
	LastModified Time   `json:"MatterTextLastModifiedUtc"`
	RowVersion   string `json:"MatterTextRowVersion"`
	MatterID     int    `json:"MatterTextMatterId"`
	Version      string `json:"MatterTextVersion"`
	Plain        string `json:"MatterTextPlain"`
	RTF          string `json:"MatterTextRtf"`
}

// SimplifiedText returns text without the "$file\nBy Councilmembers ...\n\n..Title:\n$title" preamble
func (t MatterText) SimplifiedText() string {
	t.Plain = strings.ReplaceAll(t.Plain, "\r\n", "\n")
	s := strings.Split(t.Plain, "\n")
	o := s
	for i, ss := range s {
		if strings.HasPrefix(ss, "By Council Members") {
			o = s[i+1:]
		}
		if strings.TrimSpace(strings.ToLower(ss)) == "..body" {
			o = s[i+1:]
		}
	}
	return strings.TrimSpace(strings.Join(o, "\n"))
}

// SimplifiedRTF returns RTF without the "$file\nBy Councilmembers ...\n\n..Title:\n$title" preamble
func (t MatterText) SimplifiedRTF() string {
	t.RTF = strings.ReplaceAll(t.RTF, "\r\n", "\n")
	s := strings.Split(t.RTF, "\n")
	if len(s) < 3 {
		return t.RTF
	}
	// \viewkindN ; 4 == normal
	// \ucN ; unicode char; 1 == "start of header"
	// \pard  is paragraph definition

	// we want to skip till we get \viewkind and
	// strip till we get the paragraph line for ..Body
	foundStart := -1
	for i, ss := range s {
		// log.Printf("s[%d] %q", i, ss)
		if strings.Contains(ss, `\viewkind4\uc1`) && foundStart == -1 {
			s[i] = `\viewkind4\uc1`
			foundStart = i
			continue
		}
		if strings.Contains(ss, "..Body") || strings.Contains(ss, "..body") {
			if foundStart == -1 {
				break
			}
			s = append(s[:foundStart+1], s[i+1:]...)
			break
		}
	}
	return strings.Join(s, "\n")
}

// MatterTextVersion
// http://webapi.legistar.com/Help/Api/GET-v1-Client-Matters-MatterId-Versions
type MatterTextVersion struct {
	TextID  string `json:"Key"`
	Version string `json:"Value"`
}
type MatterTextVersions []MatterTextVersion

func (m MatterTextVersions) LatestTextID() int {
	if len(m) == 0 {
		return 0
	}
	// TODO: is this the right order
	n, _ := strconv.Atoi(m[len(m)-1].TextID)
	return n
}

type MatterHistory struct {
	ID              int    `json:"MatterHistoryId"`
	GUID            string `json:"MatterHistoryGuid"`
	LastModified    Time   `json:"MatterHistoryLastModifiedUtc"`
	RowVersion      string `json:"MatterHistoryRowVersion"`
	EventID         int    `json:"MatterHistoryEventId"`
	AgendaSequence  int    `json:"MatterHistoryAgendaSequence"`
	MinutesSequence int    `json:"MatterHistoryMinutesSequence"`
	AgendaNumber    string `json:"MatterHistoryAgendaNumber"`
	Video           int    `json:"MatterHistoryVideo"`
	VideoIndex      int    `json:"MatterHistoryVideoIndex"`
	Version         string `json:"MatterHistoryVersion"`
	AgendaNote      string `json:"MatterHistoryAgendaNote"`
	MinutesNote     string `json:"MatterHistoryMinutesNote"`
	ActionDate      Time   `json:"MatterHistoryActionDate"`
	ActionID        int    `json:"MatterHistoryActionId"`
	ActionName      string `json:"MatterHistoryActionName"`
	ActionText      string `json:"MatterHistoryActionText"`
	ActionBodyID    int    `json:"MatterHistoryActionBodyId"`
	ActionBodyName  string `json:"MatterHistoryActionBodyName"`
	PassedFlag      int    `json:"MatterHistoryPassedFlag"`     // make bool
	PassedFlagName  string `json:"MatterHistoryPassedFlagName"` // i.e. "Pass"
	RollCallFlag    int    `json:"MatterHistoryRollCallFlag"`
	FlagExtra       int    `json:"MatterHistoryFlagExtra"`
	Tally           string `json:"MatterHistoryTally"`
	AccelaRecordID  string `json:"MatterHistoryAccelaRecordId"`
	Consent         int    `json:"MatterHistoryConsent"`
	MoverID         int    `json:"MatterHistoryMoverId"`
	MoverName       string `json:"MatterHistoryMoverName"`
	SeconderID      int    `json:"MatterHistorySeconderId"`
	SeconderName    string `json:"MatterHistorySeconderName"`
	MatterStatusID  int    `json:"MatterHistoryMatterStatusId"`
}
type MatterHistories []MatterHistory

type MatterAttachment struct {
	ID                   int    `json:"MatterAttachmentId"`
	GUID                 string `json:"MatterAttachmentGuid"`
	LastModified         Time   `json:"MatterAttachmentLastModifiedUtc"`
	RowVersion           string `json:"MatterAttachmentRowVersion"`
	Name                 string `json:"MatterAttachmentName"`
	Link                 string `json:"MatterAttachmentHyperlink"`
	FileName             string `json:"MatterAttachmentFileName"`
	MatterVersion        string `json:"MatterAttachmentMatterVersion"` // 0
	IsHyperlink          bool   `json:"MatterAttachmentIsHyperlink,omitempty"`
	Binary               string `json:"MatterAttachmentBinary,omitempty"`
	IsSupportingDocument bool   `json:"MatterAttachmentIsSupportingDocument,omitempty"`
	ShowOnInternetPage   bool   `json:"MatterAttachmentShowOnInternetPage"`
	IsMinuteOrder        bool   `json:"MatterAttachmentIsMinuteOrder,omitempty"`
	IsBoardLetter        bool   `json:"MatterAttachmentIsBoardLetter,omitempty"`
	AgiloftID            int    `json:"MatterAttachmentAgiloftId,omitempty"`
	Description          string `json:"MatterAttachmentDescription,omitempty"`
	PrintWithReports     bool   `json:"MatterAttachmentPrintWithReports"`
	Sort                 int    `json:"MatterAttachmentSort"`
}
type MatterAttachments []MatterAttachment
