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

var DbConn *string

func handler(w http.ResponseWriter, r *http.Request) {
	var loggedTopics []topics.Topic

	conn := *DbConn
	dbconn, err := sql.Open("sqlite3", conn)
	if err != nil {
		panic(err)
	}

	query := `SELECT * from topics ORDER BY count DESC`
	rows, err := dbconn.Query(query)
	if err != nil {
		defer rows.Close()
		panic(err)
	}

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
}

func main() {
	DbConn = flag.String("db", "topics.db", "Path of db file")
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
