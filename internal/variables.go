package internal

import "database/sql"

var (
	// TableName is the name of the DynamoDB Table
	TableName     string
	database      *sql.DB
	initStatement string
)
