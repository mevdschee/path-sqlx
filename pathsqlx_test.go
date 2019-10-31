package pathsqlx

import (
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
		arg   interface{}
	}
	tests := []struct {
		name    string
		db      *DB
		args    args
		want    interface{}
		wantErr bool
	}{
		{
			"first",
			db,
			args{
				"SELECT * from posts where id=:id",
				map[string]interface{}{"id": 1},
			},
			nil,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.db.Q(tt.args.query, tt.args.arg)
			if (err != nil) != tt.wantErr {
				t.Errorf("DB.Q() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DB.Q() = %v, want %v", got, tt.want)
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