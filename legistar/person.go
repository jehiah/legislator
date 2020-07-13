package legistar

// http://webapi.legistar.com/Help/ResourceModel?modelName=GranicusPerson
type Person struct {
	ID         int    `json:"PersonId"`
	GUID       string `json:"PersonGuid"`
	ActiveFlag int    `json:"PersonActiveFlag"`
	Email      string `json:"PersonEmail"`
	// LastModified string `json:"PersonLastModifiedUtc"`
	FullName  string `json:"PersonFullName"`
	FirstName string `json:"PersonFirstName"`
	LastName  string `json:"PersonLastName"`
	WWW       string `json:"PersonWWW"`

	// Address, City, State, Zip, Phone, Fax
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
