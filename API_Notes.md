# NYC
https://council.nyc.gov/legislation/api/
http://webapi.legistar.com/Home/Examples
http://webapi.legistar.com/Help

# environment variable NYC_LEGISLATOR_TOKEN


curl "https://webapi.legistar.com/v1/nyc/matters?token=${NYC_LEGISLATOR_TOKEN}&\$top=2&\$filter=MatterIntroDate+ge+datetime'2020-06-01'" > output.json
http://127.0.0.1:7001/Matters?$filter=MatterIntroDate+ge+datetime%272020-01-01%27


 
GET v1/{Client}/Persons	
	// PersonActiveFlag
	// filter on PersonLastModifiedUtc
	http://127.0.0.1:7001/Persons?$filter=PersonLastModifiedUtc+ge+datetime%272021-01-01%27
	http://127.0.0.1:7001/Persons/7631

// i.e introductions
GET v1/{Client}/Matters	
	http://127.0.0.1:7001/Matters?$filter=MatterIntroDate+ge+datetime%272020-01-01%27+and+MatterTypeName+eq+%27Introduction%27
	// query by MatterLastModifiedUtc to keep updated
	// MatterEXText9 -> Intro'd by "last name"  ( use sponsors index 0)
	// MatterEXText10; MatterStatusName -> Status
	// MatterEXText5 -> Summary
	http://127.0.0.1:7001/Matters/66477

GET v1/{Client}/Matters/{MatterId}/Sponsors	
	http://127.0.0.1:7001/Matters/52409/Sponsors
	http://127.0.0.1:7001/Matters/65667/Sponsors
	http://127.0.0.1:7001/Matters/66477/Sponsors

	MatterSponsorSequence=null is the first sponsor; MatterSponsorSequence starts at 1 after that

	20649EFF-1236-40E3-8CE0-150057BA2A2F

	
GET v1/{Client}/Matters/{MatterId}/Relations

// "Introduced by Council", "Referred to Comm by Council", "Hearing Held by Committee"
GET v1/{Client}/Matters/{MatterId}/Histories?AgendaNote={AgendaNote}&MinutesNote={MinutesNote}	
	http://127.0.0.1:7001/Matters/52409/Histories
	http://127.0.0.1:7001/Matters/66050/Histories
	http://127.0.0.1:7001/Matters/66477/Histories

// i.e. committee report, Hearing Transcript, etc
GET v1/{Client}/Matters/{MatterId}/Attachments	
	http://127.0.0.1:7001/Matters/66050/Attachments

GET v1/{Client}/Matters/{MatterId}/Versions	
	http://127.0.0.1:7001/Matters/66477/Versions
		> Version ID "key"
	http://127.0.0.1:7001/Matters/66477/Texts/69354  // Body of bill for a version

// i.e. "report required", "Oversight"
GET v1/{Client}/Indexes	
GET v1/{Client}/Indexes/{IndexId}	

// i.e "report required", "oversight", "Agency Rule-making Required"
GET v1/{Client}/MatterIndexes	
GET v1/{Client}/MatterIndexes/{MatterIndexId}	
http://127.0.0.1:7001/MatterIndexes/1237
	http://127.0.0.1:7001/MatterIndexes?$filter=MatterIndexMatterId+eq+53046 // search by MatterId
	http://127.0.0.1:7001/MatterIndexes/Matter/66477 ==

GET v1/{Client}/EventItems/{EventItemId}/RollCalls	
GET v1/{Client}/Persons/{PersonId}/RollCalls
	http://127.0.0.1:7001/Persons/7631/RollCalls	


GET v1/{Client}/EventItems/{EventItemId}/Votes	
GET v1/{Client}/Persons/{PersonId}/Votes
	http://127.0.0.1:7001/Persons/7631/Votes	

GET v1/{Client}/VoteTypes	


GET v1/{Client}/OfficeRecords
	http://127.0.0.1:7001/OfficeRecords?$filter=OfficeRecordLastModifiedUtc+ge+datetime%272022-01-01%27

GET v1/{Client}/Persons/{PersonId}/OfficeRecords
	http://127.0.0.1:7001/Persons/7631/OfficeRecords
		OfficeRecordStartDate -> OfficeRecordEndDate	

GET v1/{Client}/Events	
GET v1/{Client}/Events/{EventId}
	http://127.0.0.1:7001/Events?$filter=EventDate+ge+datetime%272022-01-01%27
	http://127.0.0.1:7001/Events/19019?EventItems=1&AgendaNode=1&MinutesNote=1&EventItemAttachments=1
	
	// GET v1/{Client}/EventDates/{BodyId}?FutureDatesOnly={FutureDatesOnly}	
	//	GET v1/{Client}/Events/{EventId}?EventItems={EventItems}&AgendaNote={AgendaNote}&MinutesNote={MinutesNote}&EventItemAttachments={EventItemAttachments}	

	http://127.0.0.1:7001/Events/379233/EventItems?AgendaNote=1&MinutesNote=1&Attachments=1


GET v1/{Client}/Actions	
	http://127.0.0.1:7001/Actions?$filter=ActionLastModifiedUtc+ge+datetime%272020-01-01%27
GET v1/{Client}/Actions/{ActionId}	

GET v1/{Client}/Bodies // aka Committes
	http://127.0.0.1:7001/Bodies
	http://127.0.0.1:7001/Bodies/1 // City Council

	http://127.0.0.1:7001/EventDates/1 // Future City Council Events
	http://127.0.0.1:7001/EventDates/1?FutureDatesOnly=false

	http://127.0.0.1:7001/Bodies/29 // Transportation Committee
	http://127.0.0.1:7001/EventDates/29?FutureDatesOnly=false
