package legistar

import (
	"net/url"
	"time"
)

type OfficeRecord struct {
	ID              int     `json:"OfficeRecordId"`
	GUID            string  `json:"OfficeRecordGuid"`
	LastModified    Time    `json:"OfficeRecordLastModifiedUtc"`
	RowVersion      string  `json:"OfficeRecordRowVersion"`
	FirstName       string  `json:"OfficeRecordFirstName"`
	LastName        string  `json:"OfficeRecordLastName"`
	Email           string  `json:"OfficeRecordEmail"`
	FullName        string  `json:"OfficeRecordFullName"`
	StartDate       Time    `json:"OfficeRecordStartDate"`
	EndDate         Time    `json:"OfficeRecordEndDate"`
	Sort            int     `json:"OfficeRecordSort"`
	PersonID        int     `json:"OfficeRecordPersonId"`
	BodyID          int     `json:"OfficeRecordBodyId"`
	BodyName        string  `json:"OfficeRecordBodyName"`
	Title           string  `json:"OfficeRecordTitle"`
	VoteDivider     float64 `json:"OfficeRecordVoteDivider"`
	ExtendFlag      int     `json:"OfficeRecordExtendFlag"`
	MemberTypeID    int     `json:"OfficeRecordMemberTypeId"`
	MemberType      string  `json:"OfficeRecordMemberType"`
	SupportNameID   int     `json:"OfficeRecordSupportNameId"`
	SupportFullName string  `json:"OfficeRecordSupportFullName"`
}

type OfficeRecords []OfficeRecord

type OfficeRecordLastModifiedFilter time.Time

func (f OfficeRecordLastModifiedFilter) Paramters() url.Values {
	return DateTimeFilter("OfficeRecordLastModifiedUtc", time.Time(f))
}
