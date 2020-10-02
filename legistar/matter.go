package legistar

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
	EnactmentNumber    string `json:"MatterEnactmentNumber"`
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

// MatterSponsor
// http://webapi.legistar.com/Help/Api/GET-v1-Client-Matters-MatterId-Sponsors
type MatterSponsor struct {
	ID            int    `json:"MatterSponsorId"`
	GUID          string `json:"MatterSponsorGuid"`
	LastModified  Time   `json:"MatterSponsorLastModifiedUtc"`
	RowVersion    string `json:"MatterSponsorRowVersion"`
	MatterID      int    `json:"MatterSponsorMatterId"`
	MatterVersion string `json:"MatterSponsorMatterVersion"`
	NameID        int    `json:"MatterSponsorNameId"`
	BodyID        int    `json:"MatterSponsorBodyId"`
	Name          string `json:"MatterSponsorName"`
	Sequence      int    `json:"MatterSponsorSequence"`
	LinkFlag      int    `json:"MatterSponsorLinkFlag"`
}

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
