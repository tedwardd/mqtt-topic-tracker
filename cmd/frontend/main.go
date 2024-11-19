package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
	"github.com/tedwardd/mqtt-topic-tracker/internal/topics"
)

var DbPath *string

func handler(w http.ResponseWriter, r *http.Request) {
	var loggedTopics []topics.Topic

	dbconn, err := sql.Open("sqlite3", "/mqtt_data/topics.db")
	if err != nil {
		panic(err)
	}
	defer dbconn.Close()

	query := `SELECT * from topics ORDER BY count DESC`
	rows, err := dbconn.Query(query)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		var topic topics.Topic
		if err := rows.Scan(&topic.Name, &topic.Count); err != nil {
			panic(err)
		}
		fmt.Fprintf(w, "%s - %d\n", topic.Name, topic.Count)
		loggedTopics = append(loggedTopics, topic)
		if err = rows.Err(); err != nil {
			panic(err)
		}
	}

	// json, err := json.Marshal(loggedTopics)
	// fmt.Fprintf(w, "%v", json)
	return
}

func main() {
	DbPath = flag.String("db", "topics.db", "Path of db file")
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
