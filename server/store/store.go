package store

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

func InitDB() (*sql.DB, error) {
	user := os.Getenv("MYSQL_USER")
	pass := os.Getenv("MYSQL_PASSWORD")
	host := os.Getenv("MYSQL_HOST") // "localhost:13306"
	port := os.Getenv("MYSQL_PORT")
	db := os.Getenv("MYSQL_DATABASE") // "mmplugin_parser"
	if len(port) == 0 {
		host = host + ":" + "3306"
	} else {
		host = host + ":" + port
	}
	return sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s", user, pass, host, db))
}

func InsertRepository(db *sql.DB, repo, id, at string) error {
	log.Println(at)
	_, err := db.Exec("INSERT IGNORE INTO repositories SET url = ?, commit_id = ?, created_at = ?", repo, id, at)
	return err
}

func InsertUsage(db *sql.DB, commitId, api, file string, line int, t string) error {
	_, err := db.Exec("INSERT IGNORE INTO usages SET commit_id = ?, api = ?, path = ?, line = ?, type = ?", commitId, api, file, line, t)
	return err
}
