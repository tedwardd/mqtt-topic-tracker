// Provides the frontend for mqtt-topic-tracker
package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"text/template"

	_ "github.com/mattn/go-sqlite3"
	"github.com/tedwardd/mqtt-topic-tracker/internal/topics"
)

type options interface {
	Validate() error
}

type serverOpts struct {
	dbFile, tmplPath string
}

func (options serverOpts) Validate() error {
	var ret error

	dbFileInfo, err := os.Stat(options.dbFile)
	if errors.Is(err, os.ErrNotExist) {
		ret = err
	} else {
		if dbFileInfo.IsDir() {
			ret = fmt.Errorf("db path must be file")
		}
	}

	tmplFileInfo, err := os.Stat(options.tmplPath)
	if errors.Is(err, os.ErrNotExist) {
		ret = err
	} else {
		if !tmplFileInfo.IsDir() {
			ret = fmt.Errorf("tmplate path must be directory")
		}
	}
	return ret
}

func check(o options) {
	o.Validate()
}

type topicPageData struct {
	PageTitle string
	Topics    []topics.Topic
}

func (options *serverOpts) handler(w http.ResponseWriter, _ *http.Request) {
	data := topicPageData{
		PageTitle: "Topic Traffic",
		Topics:    []topics.Topic{},
	}

	dbconn, err := sql.Open("sqlite3", options.dbFile)
	if err != nil {
		log.Fatal(err.Error())
	}

	defer dbconn.Close()

	query := `SELECT * from topics ORDER BY count DESC`

	rows, err := dbconn.Query(query)
	if err != nil {
		log.Fatal(err.Error())
	}

	defer rows.Close()
	var topic topics.Topic

	tmpl := template.New("topics")

	for rows.Next() {
		if err := rows.Scan(&topic.Name, &topic.Count); err != nil {
			log.Fatal(err.Error())
		}

		// fmt.Fprintf(w, "%s - %d\n", topic.Name, topic.Count)
		tmpl, err = template.ParseFiles(options.tmplPath + "/mqtt_stats.gohtml")
		if err != nil {
			log.Fatal(err.Error())
		}

		data.Topics = append(data.Topics, topic)
		if err = rows.Err(); err != nil {
			log.Fatal(err.Error())
		}
	}

	if tmplErr := tmpl.Execute(w, data); tmplErr != nil {
		log.Fatal(tmplErr.Error())
	}
}
