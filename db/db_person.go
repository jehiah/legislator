package db

import (
	"time"

	"github.com/jehiah/legislator/legistar"
)

type PersonReference struct {
	ID       int
	Slug     string
	FullName string
}

type Person struct {
	ID                  int
	GUID                string
	Slug                string
	IsActive            bool
	Email               string
	FullName            string
	FirstName, LastName string
	WWW                 string
	LastModified        time.Time
	OfficeRecords       []OfficeRecord
	Start, End          time.Time
	DistrictOffice      Address `json:",omitempty"`
	LegislativeOffice   Address `json:",omitempty"`
}

type Address struct {
	Address string
	City    string
	State   string
	Zip     string
}

func (p Person) Reference() PersonReference {
	return PersonReference{
		ID:       p.ID,
		Slug:     p.Slug,
		FullName: p.FullName,
	}
}

type OfficeRecord struct {
	ID   int
	GUID string

	// Body
	BodyID            int
	BodyName          string
	MemberType, Title string
	MemberTypeID      int

	// Person
	FullName     string
	PersonID     int
	Start, End   time.Time
	LastModified time.Time
}

func NewOfficeRecord(o legistar.OfficeRecord) OfficeRecord {
	return OfficeRecord{
		ID:           o.ID,
		GUID:         o.GUID,
		BodyID:       o.BodyID,
		BodyName:     o.BodyName,
		MemberType:   o.MemberType,
		MemberTypeID: o.MemberTypeID,
		Title:        o.Title,
		FullName:     o.FullName,
		PersonID:     o.PersonID,
		Start:        o.StartDate.Time,
		End:          o.EndDate.Time,
		LastModified: o.LastModified.Time,
	}
}

func NewPerson(p legistar.Person, o legistar.OfficeRecords) Person {
	v := Person{
		ID:           p.ID,
		GUID:         p.GUID,
		Slug:         p.Slug(),
		IsActive:     p.ActiveFlag == 1,
		Email:        p.Email,
		FullName:     p.FullName,
		FirstName:    p.FirstName,
		LastName:     p.LastName,
		WWW:          p.WWW,
		LastModified: p.LastModified.Time,
		DistrictOffice: Address{
			Address: p.Address1,
			City:    p.City1,
			State:   p.State1,
			Zip:     p.Zip1,
		},
		LegislativeOffice: Address{
			Address: p.Address2,
			City:    p.City2,
			State:   p.State2,
			Zip:     p.Zip2,
		},
	}
	for _, oo := range o {
		v.OfficeRecords = append(v.OfficeRecords, NewOfficeRecord(oo))
	}
	v.Start = Min(v.OfficeRecords, func(i int) time.Time { return v.OfficeRecords[i].Start })
	v.End = Max(v.OfficeRecords, func(i int) time.Time { return v.OfficeRecords[i].End })
	return v
}
