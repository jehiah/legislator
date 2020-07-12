# NYC
https://council.nyc.gov/legislation/api/
http://webapi.legistar.com/Home/Examples
http://webapi.legistar.com/Help

# environment variable NYC_LEGISLATOR_TOKEN


curl "https://webapi.legistar.com/v1/nyc/matters?token=${NYC_LEGISLATOR_TOKEN}&\$top=2&\$filter=MatterIntroDate+ge+datetime'2020-06-01'" > output.json



GET v1/{Client}/Persons	
GET v1/{Client}/Matters	
GET v1/{Client}/Matters/{MatterId}/Sponsors	
GET v1/{Client}/Matters/{MatterId}/Relations	

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
