package main

import (
	"database/sql"
	"flag"
	"log"
	"net/http"
	"text/template"

	_ "github.com/mattn/go-sqlite3"
	"github.com/tedwardd/mqtt-topic-tracker/internal/topics"
)

var DbPath *string

type TopicPageData struct {
	PageTitle string
	Topics    []topics.Topic
}

func handler(w http.ResponseWriter, r *http.Request) {
	data := TopicPageData{
		PageTitle: "Topic Traffic",
		Topics:    []topics.Topic{},
	}

	dbconn, err := sql.Open("sqlite3", "topics.db")
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
	var topic topics.Topic

	tmpl := template.New("topics")

	for rows.Next() {

		if err := rows.Scan(&topic.Name, &topic.Count); err != nil {
			panic(err)
		}

		// fmt.Fprintf(w, "%s - %d\n", topic.Name, topic.Count)
		tmpl, err = template.ParseFiles("web/templates/mqtt_stats.gohtml")
		if err != nil {
			panic(err)
		}

		data.Topics = append(data.Topics, topic)
		if err = rows.Err(); err != nil {
			panic(err)
		}

	}
	tmpl.Execute(w, data)
	return
}

func main() {
	DbPath = flag.String("db", "topics.db", "Path of db file")
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
