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
	LastModified Time   `json:"PersonLastModifiedUtc"`
	FullName     string `json:"PersonFullName"`
	FirstName    string `json:"PersonFirstName"`
	LastName     string `json:"PersonLastName"`
	WWW          string `json:"PersonWWW"`

	// Address, City, State, Zip, Phone, Fax
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
