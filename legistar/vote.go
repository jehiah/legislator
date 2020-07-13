package legistar

// http://webapi.legistar.com/Help/ResourceModel?modelName=GranicusVote
type Vote struct {
	ID           string `json:"VoteId"`
	GUID         string `json:"VoteGuid"`
	PersonID     string `json:"VotePersonId"`
	PersonName   string `json:"VotePersonName"`
	EventItemID  string `json:"EventItemId"`
	LastModified string `json:"VoteLastModifiedUtc"`
	RowVersion   []byte `json:"VoteRowVersion"`

	Result    int    `json:"VoteResult"`
	ValueID   int    `json:"ValueId"`
	ValueName string `json:"ValueName"` // i.e. Affirmative
	Sort      int    `json:"VoteSort"`
}
