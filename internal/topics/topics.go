package topics

import (
	"database/sql"
	"log"
)

type Topic struct {
	Name  string `field:"topic"`
	Count int64  `field:"count"`
}

func InsertNewTopic(dbconn *sql.DB, topic string) {
	_, err := dbconn.Exec("INSERT INTO topics VALUES(?, 1);", topic)
	if err != nil {
		log.Fatal(err)
	}

	defer dbconn.Close()
}

func IncrementTopicCount(dbconn *sql.DB, topic string) {
	_, err := dbconn.Exec("UPDATE topics SET count = count + 1 WHERE topic = ?;", topic)
	if err != nil {
		log.Fatal(err)
	}

	defer dbconn.Close()
}
