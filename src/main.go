package main

import (
	"crypto/tls"
	"database/sql"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	MQTT "github.com/eclipse/paho.mqtt.golang"
	_ "github.com/mattn/go-sqlite3"
)

type Topic struct {
	Name  string `field:"topic"`
	Count int64  `field:"count"`
}

const file string = "topics.db"

const create string = `
	CREATE TABLE IF NOT EXISTS topics (
	topic TEXT NOT NULL PRIMARY KEY,
	count INTEGER NOT NULL
	);`

func insertNewTopic(dbconn *sql.DB, topic string) {
	_, err := dbconn.Exec("INSERT INTO topics VALUES(?, 1);", topic)
	if err != nil {
		panic(err)
	}

	defer dbconn.Close()
}

func incrementTopicCount(dbconn *sql.DB, topic string) {
	_, err := dbconn.Exec("UPDATE topics SET count = count + 1 WHERE topic = ?;", topic)
	if err != nil {
		panic(err)
	}

	defer dbconn.Close()
}

func onMessageReceived(_ MQTT.Client, message MQTT.Message) {
	dbconn, err := sql.Open("sqlite3", file)
	if err != nil {
		panic(err)
	}

	messageTopic := Topic{
		Name:  message.Topic(),
		Count: 1,
	}

	query := `SELECT topic FROM topics WHERE topic = ?`
	topic := messageTopic.Name

	rows, err := dbconn.Query(query, topic)
	if err != nil {
		defer rows.Close()
		panic(err)
	}

	if rows.Next() {
		defer rows.Close()
		incrementTopicCount(dbconn, topic)
	} else {
		defer rows.Close()
		insertNewTopic(dbconn, topic)
	}
}

func main() {
	dbconn, err := sql.Open("sqlite3", file)
	if err != nil {
		panic(err)
	}

	if _, err := dbconn.Exec(create); err != nil {
		panic(err)
	}

	channel := make(chan os.Signal, 1)
	signal.Notify(channel, os.Interrupt, syscall.SIGTERM)

	hostname, _ := os.Hostname()

	server := flag.String("server", "tcp://127.0.0.1:1883", "The full url of the MQTT server to connect to ex: tcp://127.0.0.1:1883")
	topic := flag.String("topic", "#", "Topic to subscribe to")
	qos := flag.Int("qos", 0, "The QoS to subscribe to messages at")
	clientid := flag.String("clientid", hostname+strconv.Itoa(time.Now().Second()), "A clientid for the connection")
	username := flag.String("username", "", "A username to authenticate to the MQTT server")
	password := flag.String("password", "", "Password to match username")
	flag.Parse()

	connOpts := MQTT.NewClientOptions().AddBroker(*server).SetClientID(*clientid).SetCleanSession(true)
	// Start by checking to see if we have user/pass set in environment
	user, ret := os.LookupEnv("MQTT_USERNAME")
	if ret {
		connOpts.SetUsername(user)
	}

	pass, ret := os.LookupEnv("MQTT_PASSWORD")
	if ret {
		connOpts.SetPassword(pass)
	}

	// Override env values with flags if passed
	if *username != "" {
		connOpts.SetUsername(*username)
	}

	if *password != "" {
		connOpts.SetPassword(*password)
	}

	tlsConfig := &tls.Config{InsecureSkipVerify: true, ClientAuth: tls.NoClientCert}
	connOpts.SetTLSConfig(tlsConfig)

	connOpts.OnConnect = func(channel MQTT.Client) {
		if token := channel.Subscribe(*topic, byte(*qos), onMessageReceived); token.Wait() && token.Error() != nil {
			panic(token.Error())
		}
	}

	client := MQTT.NewClient(connOpts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	fmt.Printf("Connected to %s\n", *server)

	<-channel
}
