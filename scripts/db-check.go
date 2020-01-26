package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	user := os.Getenv("MYSQL_USER")
	pass := os.Getenv("MYSQL_PASSWORD")
	host := os.Getenv("MYSQL_HOST")
	port := os.Getenv("MYSQL_PORT")
	database := os.Getenv("MYSQL_DATABASE")
	if len(port) == 0 {
		host = host + ":" + "3306"
	} else {
		host = host + ":" + port
	}

	// connect
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s", user, pass, host, database))
	if err != nil {
		log.Fatal(err.Error())
	}

	// cheeck if db has data
	row := db.QueryRow("SELECT COUNT(*) from repositories;")
	var val int
	if err := row.Scan(&val); err != nil {
		log.Fatal(err.Error())
	}
	if val > 0 {
		log.Printf("DB has already had %d repositories.", val)
		os.Exit(0)
	}

	// import data from text if any data are inserted
	if len(os.Args) == 0 {
		log.Printf("Could not find data file. You should specify json file for setting up database.")
		os.Exit(1)
	}
	file := os.Args[0]
	data, err := parseJson(file)
	if err != nil {
		log.Fatal(err.Error())
	}
	for _, d := range data {
		if _, err := db.Exec("INSERT IGNORE INTO repositories SET url = ?, commit_id = ?, created_at = ?, refs = ?", d.Url, d.CommitId, d.CreatedAt, d.Refs); err != nil {
			log.Println("Failed to insert repository data.")
			log.Printf("  Data: %#v", d)
			log.Printf("  Details: %v", err)
		}
		if _, err := db.Exec("INSERT IGNORE INTO usages SET commit_id = ?, api = ?, path = ?, line = ?, type = ?", d.CommitId, d.Api, d.Path, d.Line, d.Type); err != nil {
			log.Println("Failed to insert usage data.")
			log.Printf("  Data: %#v", d)
			log.Printf("  Details: %v", err)
		}
	}
	log.Println("Complete to set up database successfully.")
}

func parseJson(file string) ([]data, error) {
	b, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	var ret []data
	if err := json.Unmarshal(b, &ret); err != nil {
		return nil, err
	}
	return ret, nil
}

type data struct {
	Url       string `json:"url"`
	CommitId  string `json:"commit_id"`
	CreatedAt string `json:"created_at"`
	Refs      string `json:"refs"`
	Api       string `json:"api"`
	Path      string `json:"path"`
	Line      int    `json:"line"`
	Type      string `json:"type"`
}
