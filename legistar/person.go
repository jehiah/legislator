package legistar

import (
	"net/url"
	"time"

	"github.com/gosimple/slug"
)

// http://webapi.legistar.com/Help/ResourceModel?modelName=GranicusPerson
// Note: LastModified not updated for Email, WWW or OfficeRecords updates
type Person struct {
	ID           int    `json:"PersonId"`
	GUID         string `json:"PersonGuid"`
	ActiveFlag   int    `json:"PersonActiveFlag"`
	Email        string `json:"PersonEmail"`
	Email2       string `json:"PersonEmail2"`
	LastModified Time   `json:"PersonLastModifiedUtc"`
	FullName     string `json:"PersonFullName"`
	FirstName    string `json:"PersonFirstName"`
	LastName     string `json:"PersonLastName"`
	WWW          string `json:"PersonWWW"`
	WWW2         string `json:"PersonWWW2"`

	// Address, City, State, Zip, Phone, Fax
	Address1 string `json:"PersonAddress1"`
	City1    string `json:"PersonCity1"`
	State1   string `json:"PersonState1"`
	Zip1     string `json:"PersonZip1"`
	Phone    string `json:"PersonPhone"`
	Fax      string `json:"PersonFax"`

	Address2 string `json:"PersonAddress2"`
	City2    string `json:"PersonCity2"`
	State2   string `json:"PersonState2"`
	Zip2     string `json:"PersonZip2"`
	Phone2   string `json:"PersonPhone2"`
	Fax2     string `json:"PersonFax2"`
}

type PersonLastModifiedFilter time.Time

func (p PersonLastModifiedFilter) Paramters() url.Values {
	return DateTimeFilter("PersonLastModifiedUtc", "gt", time.Time(p))
}

func (p Person) Slug() string {
	return slug.MakeLang(p.FullName, "en")
}

type Persons []Person

// Filter to active persons
func (pp Persons) Active() Persons {
	var o Persons
	for _, p := range pp {
		if p.ActiveFlag > 0 {
			o = append(o, p)
		}
	}
	return o
}
func (pp Persons) Lookup() map[string]Person {
	m := make(map[string]Person, len(pp))
	for _, p := range pp {
		m[p.Slug()] = p
	}
	return m
}
