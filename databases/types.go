package databases

import "database/sql"

type LoginDBClient struct {
	DBFile string
	db     *sql.DB
}

type MMDB struct {
	MMFile string
}
type SupermanDatabases struct {
	LoginDBClient *LoginDBClient
	MMDB          MMDB
}
