# NYC
https://council.nyc.gov/legislation/api/
http://webapi.legistar.com/Home/Examples
http://webapi.legistar.com/Help

# environment variable NYC_LEGISLATOR_TOKEN


curl "https://webapi.legistar.com/v1/nyc/matters?token=${NYC_LEGISLATOR_TOKEN}&\$top=2&\$filter=MatterIntroDate+ge+datetime'2020-06-01'" > output.json
http://127.0.0.1:7001/Matters?$filter=MatterIntroDate+ge+datetime%272020-01-01%27



GET v1/{Client}/Persons	
GET v1/{Client}/Matters	
	http://127.0.0.1:7001/Matters?$filter=MatterIntroDate+ge+datetime%272020-01-01%27
GET v1/{Client}/Matters/{MatterId}/Sponsors	
	http://127.0.0.1:7001/Matters/52409/Sponsors
	http://127.0.0.1:7001/Matters/65667/Sponsors
	
GET v1/{Client}/Matters/{MatterId}/Relations
GET v1/{Client}/Matters/{MatterId}/Histories?AgendaNote={AgendaNote}&MinutesNote={MinutesNote}	
	http://127.0.0.1:7001/Matters/52409/Histories
	http://127.0.0.1:7001/Matters/66050/Histories


GET v1/{Client}/Indexes	
GET v1/{Client}/Indexes/{IndexId}	
GET v1/{Client}/MatterIndexes	
GET v1/{Client}/MatterIndexes/{MatterIndexId}	


GET v1/{Client}/EventItems/{EventItemId}/RollCalls	
GET v1/{Client}/Persons/{PersonId}/RollCalls	


GET v1/{Client}/EventItems/{EventItemId}/Votes	
GET v1/{Client}/Persons/{PersonId}/Votes	

GET v1/{Client}/VoteTypes	


GET v1/{Client}/OfficeRecords	
GET v1/{Client}/Persons/{PersonId}/OfficeRecords	

GET v1/{Client}/Events	
GET v1/{Client}/Events/{EventId}

GET v1/{Client}/Actions	
GET v1/{Client}/Actions/{ActionId}	
