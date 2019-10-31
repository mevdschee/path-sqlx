package pathsqlx

import (
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
)

// DB is a wrapper around sqlx.DB which keeps track of the driverName upon Open,
// used mostly to automatically bind named queries using the right bindvars.
type DB struct {
	*sqlx.DB
}

// Q is the query that returns nested paths
func (db *DB) Q(query string, arg interface{}) (interface{}, error) {
	rows, err := db.NamedQuery(query, arg)
	if err != nil {
		log.Fatalln(err)
	}
	for rows.Next() {
		row, err := rows.SliceScan()
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Printf("%#v\n", row)
	}
	return nil, nil
}

// Create a sqlx connection
func Create(user, password, dbname, driver, host, port string) (*DB, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err := sqlx.Connect(driver, dsn)
	return &DB{db}, err
}
