package pathsqlx

import (
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
)

// DB is a wrapper around sqlx.DB
type DB struct {
	*sqlx.DB
}

// Q is the query that returns nested paths
func (db *DB) Q(query string, arg interface{}) (interface{}, error) {
	rows, err := db.NamedQuery(query, arg)
	if err != nil {
		log.Fatalln(err)
	}
	results := []interface{}{}
	for rows.Next() {
		row, err := rows.SliceScan()
		if err != nil {
			log.Fatalln(err)
		}
		results = append(results, row)
	}
	return results, nil
}

// Create a pathsqlx connection
func Create(user, password, dbname, driver, host, port string) (*DB, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err := sqlx.Connect(driver, dsn)
	return &DB{db}, err
}
