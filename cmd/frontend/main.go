// Provides the frontend for mqtt-topic-tracker
package main

import (
	"flag"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	options := &serverOpts{}

	flag.StringVar(&options.dbFile, "db", "topics.db", "Path of db file")
	flag.StringVar(&options.tmplPath, "t", "/app/web/templates", "Path to go templates")
	flag.Parse()

	check(options)

	http.HandleFunc("/", options.handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
