package store

import (
	"database/sql"
	"fmt"
	"os"
	"reflect"

	_ "github.com/go-sql-driver/mysql"

	"github.com/mattermost/mattermost-server/model"
)

func InitDB() (*sql.DB, error) {
	user := os.Getenv("MYSQL_USER")
	pass := os.Getenv("MYSQL_PASSWORD")
	host := os.Getenv("MYSQL_HOST")
	port := os.Getenv("MYSQL_PORT")
	db := os.Getenv("MYSQL_DATABASE")
	if len(port) == 0 {
		host = host + ":" + "3306"
	} else {
		host = host + ":" + port
	}
	return sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s", user, pass, host, db))
}

func InsertRepository(db *sql.DB, repo, id, at, refs string) error {
	_, err := db.Exec("INSERT IGNORE INTO repositories SET url = ?, commit_id = ?, created_at = ?, refs = ?", repo, id, at, refs)
	return err
}

type DummyScanner struct{}

func (DummyScanner) Scan(interface{}) error {
	return nil
}

func InsertManifest(db *sql.DB, commitId string, manifest *model.Manifest) error {
	id := manifest.Id
	name := manifest.Name
	version := manifest.Version
	minServerVersion := manifest.MinServerVersion
	// TODO: If duplicated, insert data after deleting existing data
	_, err := db.Exec("INSERT IGNORE INTO manifest SET commit_id = ?, id = ?, name = ?, version = ?, min_server_version = ?", commitId, id, name, version, minServerVersion)
	if err != nil {
		return err
	}

	settingsSchema := manifest.SettingsSchema
	if settingsSchema != nil {
		var header bool
		var footer bool
		if settingsSchema.Header != "" {
			header = true
		}
		if settingsSchema.Footer != "" {
			footer = true
		}
		rows, err := db.Query("INSERT IGNORE INTO settings_schema VALUE (?, ?, ?)", commitId, header, footer)
		if err != nil {
			return err
		}
		defer rows.Close()

		for _, v := range settingsSchema.Settings {
			if _, err := db.Exec("INSERT IGNORE INTO plugin_settings VALUE (?, ?, ?)", commitId, v.Key, v.Type); err != nil {
				return err
			}
		}
		for k, v := range manifest.Props {
			if _, err := db.Exec("INSERT IGNORE INTO props VALUE (?, ?, ?)", commitId, k, reflect.TypeOf(v).String()); err != nil {
				return err
			}
		}
	}
	return nil
}

func InsertUsage(db *sql.DB, commitId, api, file string, line int, t string) error {
	_, err := db.Exec("INSERT IGNORE INTO usages SET commit_id = ?, api = ?, path = ?, line = ?, type = ?", commitId, api, file, line, t)
	return err
}
