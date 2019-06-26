package databases

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"github.com/oschwald/geoip2-golang"
	"log"
	"net"
	"os"
)

func CheckIfExists(dbFile string) error {

	if _, err := os.Stat(dbFile); os.IsNotExist(err) {
		return fmt.Errorf("db file does not exists %v", err)
	}

	return nil
}

func NewDBClient(client *SupermanDatabases) (*SupermanDatabases, error) {


	db, err := sql.Open("sqlite3", client.LoginDBClient.DBFile)
	if err != nil {
		return nil, fmt.Errorf("could not open new DB due to %v", err)
	}

	client.LoginDBClient.db = db

	return client, nil
}

func (super *SupermanDatabases) CreateDB() error {

	var createTableStmt = `create table logins(
    							id        integer     not null	constraint table_name_pk primary key autoincrement,
                                username  string  default 'testuser' not null,
                                ipaddress varchar(13) not null,
                                timestamp integer default 0 not null
                            );

							create unique index logins_timestamp_uindex
								on logins (timestamp);
							
							create unique index table_name_id_uindex
								on logins (id);
							
							`
	if _, err := super.LoginDBClient.db.Exec(createTableStmt); err != nil {
		return fmt.Errorf("ould not create new table for logins database %v", err)

	}

	return nil

}

func (super *SupermanDatabases) LoadDataset(ipaddress, username, timestamp string) error {

	insertRow := `INSERT into logins (username, ipaddress, timestamp)
                   VALUES ($1, $2, $3)`

	timestampInt := string(timestamp)

	if _, err := super.LoginDBClient.db.Exec(insertRow, username, ipaddress, timestampInt); err != nil {
		return fmt.Errorf("row could not be inserted due to %v", err)
	}

	return nil

}

func (super *SupermanDatabases) RetrieveEventsByUsername(user string) ([]map[string]interface{}, error) {

	var resultMapArr []map[string]interface{}
	var resultMap map[string]interface{}

	queryString := `select * from logins where username=$1 order by timestamp limit 3;`

	rows, err := super.LoginDBClient.db.Query(queryString, user)
	if err != nil {
		return nil, fmt.Errorf("error occured when preparing statement %v", err)
	}

	defer rows.Close()

	for rows.Next() {
		var id string
		var username string
		var ipaddress string
		var timestamp int32

		if err := rows.Scan(&id, &username, &ipaddress, &timestamp); err != nil {
			return nil, fmt.Errorf("error occured when scanning resultset %v", err)
		}

		resultMap = map[string]interface{}{
			"ipaddress": ipaddress,
			"username":  username,
			"timestamp": timestamp,
		}

		resultMapArr = append(resultMapArr, resultMap)

	}

	return resultMapArr, nil

}

func (super *SupermanDatabases) LookupIp(ipAddress string) (*geoip2.City, error) {
	mm, err := geoip2.Open(super.MMDB.MMFile)
	if err != nil {
		log.Fatalf("could not open database %v", err)
	}

	defer mm.Close()

	parsedIp := net.ParseIP(ipAddress)

	cityData, err := mm.City(parsedIp)
	if err != nil {
		return nil, fmt.Errorf("could not retrieve ip address due to %v", err)
	}

	return cityData, nil

}

func (super *SupermanDatabases) Close() error {

	if err := super.LoginDBClient.db.Close(); err != nil {
		return fmt.Errorf("error when attempting to close DB %v", err)
	}

	return nil

}
