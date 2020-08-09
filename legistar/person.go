package legistar

import (
	"github.com/gosimple/slug"
)

// http://webapi.legistar.com/Help/ResourceModel?modelName=GranicusPerson
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
