package legistar

import "github.com/gosimple/slug"

// http://webapi.legistar.com/Help/ResourceModel?modelName=GranicusVote
type Vote struct {
	ID           int    `json:"VoteId"`
	GUID         string `json:"VoteGuid"`
	PersonID     int    `json:"VotePersonId"`
	PersonName   string `json:"VotePersonName"`
	EventItemID  int    `json:"VoteEventItemId"`
	LastModified Time   `json:"VoteLastModifiedUtc"`
	RowVersion   []byte `json:"VoteRowVersion"`

	Result    int    `json:"VoteResult"` // 1=Afirmative, 2=Negative
	ValueID   int    `json:"VoteValueId"`
	ValueName string `json:"VoteValueName"` // i.e. Affirmative, Negative, ABsent
	Sort      int    `json:"VoteSort"`
}

func (s Vote) Slug() string {
	return slug.MakeLang(s.PersonName, "en")
}

type Votes []Vote

type VoteType struct {
	ID           int    `json:"VoteTypeId"`
	GUID         string `json:"VoteTypeGuid"`
	LastModified Time   `json:"VoteTypeLastModifiedUtc"`
	RowVersion   string `json:"VoteTypeRowVersion"`
	Name         string `json:"VoteTypeName"`
	PluralName   string `json:"VoteTypePluralName"`
	UsedFor      int    `json:"VoteTypeUsedFor"`
	Result       int    `json:"VoteTypeResult"`
	Sort         int    `json:"VoteTypeSort"`
}
type VoteTypes []VoteType

func (v VoteTypes) Find(ID int) VoteType {
	for _, vv := range v {
		if vv.ID == ID {
			return vv
		}
	}
	return VoteType{}
}
