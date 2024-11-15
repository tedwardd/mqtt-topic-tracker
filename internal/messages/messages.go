package messages

import (
	"database/sql"

	MQTT "github.com/eclipse/paho.mqtt.golang"
	_ "github.com/mattn/go-sqlite3"
	"github.com/tedwardd/mqtt-topic-tracker/internal/logger"
	"github.com/tedwardd/mqtt-topic-tracker/internal/topics"
)

var (
	DbConn string
	Logger *logger.LogCache
)

func OnMessageReceived(_ MQTT.Client, message MQTT.Message) {
	dbconn, err := sql.Open("sqlite3", DbConn)
	if err != nil {
		panic(err)
	}

	messageTopic := topics.Topic{
		Name:  message.Topic(),
		Count: 1,
	}

	query := `SELECT topic FROM topics WHERE topic = ?`
	topic := messageTopic.Name

	rows, err := dbconn.Query(query, topic)
	if err != nil {
		defer rows.Close()
		Logger.Errorf(err.Error())
	}
	defer dbconn.Close()

	if rows.Next() {
		rows.Close()
		Logger.Verbosef("Incrementing count for %s", topic)
		topics.IncrementTopicCount(dbconn, topic)
	} else {
		Logger.Verbosef("First time logging %s", topic)
		topics.InsertNewTopic(dbconn, topic)
		defer rows.Close()
	}
}
