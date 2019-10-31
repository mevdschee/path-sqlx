package pathsqlx

import (
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
)

func main() {
	// this Pings the database trying to connect, panics on error
	// use sqlx.Open() for sql.Open() semantics
	host := "127.0.0.1"
	port := "5432"
	user := "php-crud-api"
	password := "php-crud-api"
	dbname := "php-crud-api"
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		log.Fatalln(err)
	}

	// Selects Mr. Smith from the database
	rows, err := db.NamedQuery(`SELECT * FROM posts WHERE id=:id`, map[string]interface{}{"id": 1})
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
}
