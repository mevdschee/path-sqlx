package pathsqlx

import (
	"encoding/json"
	"log"
	"reflect"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
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
			"single record no path", db,
			args{"select id, content from posts where id=:id", `{"id": 1}`},
			`[{"id":1,"content":"blog started"}]`, false,
		}, {
			"two records no path", db,
			args{"select id from posts where id<=2 order by id", `{}`},
			`[{"id":1},{"id":2}]`, false,
		}, {
			"two records named no path", db,
			args{"select id from posts where id<=:two and id>=:one order by id", `{"one": 1, "two": 2}`},
			`[{"id":1},{"id":2}]`, false,
		}, {
			"two tables with path", db,
			args{`select posts.id as "$[].posts.id", comments.id as "$[].comments.id" from posts left join comments on post_id = posts.id where posts.id=1`, `{}`},
			`[{"posts":{"id":1},"comments":{"id":1}},{"posts":{"id":1},"comments":{"id":2}}]`, false,
		}, {
			"posts with comments properly nested", db,
			args{`select posts.id as "$.posts[].id", comments.id as "$.posts[].comments[].id" from posts left join comments on post_id = posts.id where posts.id<=2 order by posts.id, comments.id`, `{}`},
			`{"posts":[{"id":1,"comments":[{"id":1},{"id":2}]},{"id":2,"comments":[{"id":3},{"id":4},{"id":5},{"id":6}]}]}`, false,
		}, {
			"comments with post properly nested", db,
			args{`select posts.id as "$.comments[].post.id", comments.id as "$.comments[].id" from posts left join comments on post_id = posts.id where posts.id<=2 order by comments.id, posts.id`, `{}`},
			`{"comments":[{"id":1,"post":{"id":1}},{"id":2,"post":{"id":1}},{"id":3,"post":{"id":2}},{"id":4,"post":{"id":2}},{"id":5,"post":{"id":2}},{"id":6,"post":{"id":2}}]}`, false,
		}, {
			"count posts with simple alias", db,
			args{`select count(*) as "posts" from posts`, `{}`},
			`[{"posts":12}]`, false,
		}, {
			"count posts with path", db,
			args{`select count(*) as "$[].posts" from posts`, `{}`},
			`[{"posts":12}]`, false,
		}, {
			"count posts as object with path", db,
			args{`select count(*) as "$.posts" from posts`, `{}`},
			`{"posts":12}`, false,
		}, {
			"count posts grouped no path", db,
			args{`select categories.name, count(posts.id) as "post_count" from posts, categories where posts.category_id = categories.id group by categories.name order by categories.name`, `{}`},
			`[{"name":"announcement","post_count":11},{"name":"article","post_count":1}]`, false,
		}, {
			"count posts with added root set in path", db,
			args{`select count(*) as "$.statistics.posts" from posts`, `{}`},
			`{"statistics":{"posts":12}}`, false,
		}, {
			"count posts and comments as object with path", db,
			args{`select (select count(*) from posts) as "$.stats.posts", (select count(*) from comments) as "comments"`, `{}`},
			`{"stats":{"posts":12,"comments":6}}`, false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var args map[string]interface{}
			err := json.Unmarshal([]byte(tt.args.arg), &args)
			if err != nil {
				log.Fatal("Cannot decode to JSON ", err)
			}
			got, err := tt.db.PathQuery(tt.args.query, args)
			if (err != nil) != tt.wantErr {
				t.Errorf("DB.Q() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			json, err := json.Marshal(got)
			if err != nil {
				log.Fatal("Cannot encode to JSON ", err)
			}
			if !reflect.DeepEqual(string(json), tt.want) {
				t.Errorf("DB.Q() = %v, want %v", string(json), tt.want)
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
