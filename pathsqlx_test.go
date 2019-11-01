package pathsqlx

import (
	"encoding/json"
	"log"
	"reflect"
	"testing"

	_ "github.com/lib/pq"
)

var db *DB

func init() {
	var err error
	db, err = Create("php-crud-api", "php-crud-api", "php-crud-api", "postgres", "127.0.0.1", "5432")
	if err != nil {
		log.Fatalln(err)
	}
}

func TestDB_Q(t *testing.T) {
	type args struct {
		query string
		arg   string
	}
	tests := []struct {
		name    string
		db      *DB
		args    args
		want    string
		wantErr bool
	}{
		{
			"first",
			db,
			args{
				"SELECT * from posts where id=:id",
				`{"id": 1}`,
			},
			`[[1,1,1,"blog started"]]`,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var args map[string]interface{}
			err := json.Unmarshal([]byte(tt.args.arg), &args)
			if err != nil {
				log.Fatal("Cannot decode to JSON ", err)
			}
			got, err := tt.db.Q(tt.args.query, args)
			if (err != nil) != tt.wantErr {
				t.Errorf("DB.Q() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			json, err := json.Marshal(got)
			if err != nil {
				log.Fatal("Cannot encode to JSON ", err)
			}
			if !reflect.DeepEqual(string(json), tt.want) {
				t.Errorf("DB.Q() = %s, want %s", json, tt.want)
			}
		})
	}
}

func TestCreate(t *testing.T) {
	type args struct {
		user     string
		password string
		dbname   string
		driver   string
		host     string
		port     string
	}
	tests := []struct {
		name    string
		args    args
		want    *DB
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Create(tt.args.user, tt.args.password, tt.args.dbname, tt.args.driver, tt.args.host, tt.args.port)
			if (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Create() = %v, want %v", got, tt.want)
			}
		})
	}
}
