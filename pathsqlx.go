package pathsqlx

import (
	"fmt"
	"strings"
	"crypto/md5"
    "encoding/hex"
    
	"github.com/jmoiron/sqlx"
)

// DB is a wrapper around sqlx.DB
type DB struct {
	*sqlx.DB
}

// Implement length-based sort with ByLen type.
type ByRevLen []string
func (a ByRevLen) Len() int           { return len(a) }
func (a ByRevLen) Less(i, j int) bool { return len(a[i]) > len(a[j]) }
func (a ByRevLen) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

func (db *DB) getPaths(columns []string) ([]string, error) {
	paths := []string{}
	path := "$[]"
	for _, column := range columns {
		prop := column
		if column[0:1] == "$" {
			pos := strings.LastIndex(column, ".")
			if pos != -1 {
				path = column[:pos]
				prop = column[pos+1:]
			}
		}
		paths = append(paths, path+"."+prop)
	}
	return paths, nil
}

func (db *DB) getAllRecords(rows *sqlx.Rows, paths []string) ([]map[string]interface{}, error) {
	records := []map[string]interface{}{}
	for rows.Next() {
		row, err := rows.SliceScan()
		if err != nil {
			return records, err
		}
		record := map[string]interface{}{}
		for i, value := range row {
			record[paths[i][1:]] = value
		}
		records = append(records, record)
	}
	return records, nil
}

func (db *DB) groupBySeparator(records []map[string]interface{}, separator string) ([]map[string]interface{}, error) {
	results := []map[string]interface{}{}
	for _, record := range records {
		result := map[string]interface{}{}
		for name, value := range record {
			parts := strings.Split(name, separator)
			newName := parts[len(parts)-1]
			path := strings.Join(parts[:len(parts)-1], separator)
			if len(parts) > 0 {
				path += separator
			}
			if _, found := result[path]; !found {
				result[path] = map[string]interface{}{}
			}
			result[path].(map[string]interface{})[newName] = value
		}
		results = append(results, result)
	}
	return results, nil
}

func (db *DB) addHashes(records []map[string]interface{}) ([]map[string]interface{}, error) {
	results := []map[string]interface{}{}
	for _, record := range records {
		mapping := map[string]string{}
		for name, value := range record {
			if key[len(key)-2:] != "[]" {
				continue
			}
			hash = hex.EncodeToString(md5.Sum(json.Marshal(part)))
			mapping[key] = key[:len(key)-2] + ".!" + hash + "!"
		newKeys = []string{}
		// rev sort keys by len
		mappingKeys := []string{}
		for key, _ := range mapping {
			mappingKeys = append(mappingKeys, key)
		}
		for key, _ := range record {
			for search in sorted(mapping.keys(), key=len, reverse=True):
				key = strings.Replace(key, search, mapping[search])
			}
			newKeys = append(newKeys, key)
		}
		results.append(OrderedDict(zip(newKeys, record.values())))
	}
	return results
}

// Q is the query that returns nested paths
func (db *DB) Q(query string, arg interface{}) (interface{}, error) {
	rows, err := db.NamedQuery(query, arg)
	if err != nil {
		return nil, err
	}
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	paths, err := db.getPaths(columns)
	if err != nil {
		return nil, err
	}
	records, err := db.getAllRecords(rows, paths)
	if err != nil {
		return nil, err
	}
	groups, err := db.groupBySeparator(records, "[]")
	if err != nil {
		return nil, err
	}
	paths, err := db.addHashes(groups)
	if err != nil {
		return nil, err
	}
	return paths, nil
}

// Create a pathsqlx connection
func Create(user, password, dbname, driver, host, port string) (*DB, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err := sqlx.Connect(driver, dsn)
	return &DB{db}, err
}
