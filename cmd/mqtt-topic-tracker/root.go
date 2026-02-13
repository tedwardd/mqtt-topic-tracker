package main

import (
	"database/sql"
	"os"
	"strconv"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/cobra"
	"github.com/tedwardd/mqtt-topic-tracker/internal/messages"
	"github.com/tedwardd/mqtt-topic-tracker/internal/mqtt"
)

const create string = `
	CREATE TABLE IF NOT EXISTS topics (
	topic TEXT NOT NULL PRIMARY KEY,
	count INTEGER NOT NULL
	);`

func newRootCmd(basics *commandBasics) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "mqtt-topic-tracker",
		Short: "Log and increment a count of mqtt topics seen to a DB",
		Long: `mqtt-topic-tracker is a tool that will connect to an MQTT Broker and log all
topics seen and increment the count of how many times a given topic was seen`,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			basics.setupLogCache(cmd.OutOrStderr(), cmd.ErrOrStderr())

			config := appConfig{}
			config.co = initColorOutput(cmd.OutOrStdout(), cmd.ErrOrStderr())

			if basics.verbose {
				basics.LogCache.EnableVerbose()
			}
			if basics.debug {
				basics.LogCache.EnableDebug()
			}

			calledCmd := cmd.CalledAs()

			basics.LogCache.SetCommandName(calledCmd)

			return nil
		},
		Run: func(c *cobra.Command, args []string) {
			run(c, basics)
		},
	}
	rootCmd.PersistentFlags().String("server", "tcp://127.0.0.1:1883", "The full uri of the MQTT server to connect to ex: tcp://127.0.0.1:1883")
	rootCmd.PersistentFlags().String("topic", "#", "Topic to subscribe to")
	rootCmd.PersistentFlags().Int("qos", 0, "The QoS to subscribe to messages at")
	rootCmd.PersistentFlags().String("clientid", hostname+strconv.Itoa(time.Now().Second()), "A clientid for the connection")
	rootCmd.PersistentFlags().String("username", "", "A username to authenticate to the MQTT server")
	rootCmd.PersistentFlags().String("password", "", "Password to match username")
	rootCmd.PersistentFlags().String("db", "topics.db", "DB File Path")
	rootCmd.PersistentFlags().String("logfile", "mqtt-topic-tracker.log", "Log file path")
	rootCmd.PersistentFlags().BoolVar(&basics.debug, "debug", false, "Enable debug logging")
	rootCmd.PersistentFlags().BoolVar(&basics.verbose, "verbose", false, "Enable verbose logging")

	return rootCmd
}

func run(c *cobra.Command, basics *commandBasics) {
	var err error
	var ret bool

	l := basics.LogCache

	channel := make(chan os.Signal, 1)

	options := mqtt.Options{}

	options.Username, ret = os.LookupEnv("MQTT_USERNAME")
	if !ret {
		options.Username = c.Flag("username").Value.String()
		l.Debugf("Using flag value for username")
	}

	options.Password, ret = os.LookupEnv("MQTT_PASSWORD")
	if !ret {
		options.Password = c.Flag("password").Value.String()
		l.Debugf("Using flag value for password")
	} else {
		l.Debugf("Using value from MQTT_PASSWORD as password")
	}

	options.Server, ret = os.LookupEnv("MQTT_SERVER")
	if !ret {
		options.Server = c.Flag("server").Value.String()
		l.Debugf("Using flag value for server")
	}

	options.Topic, ret = os.LookupEnv("MQTT_TOPIC")
	if !ret {
		options.Topic = c.Flag("topic").Value.String()
		l.Debugf("Using flag value for topic")
	}

	options.Id, ret = os.LookupEnv("MQTT_CLIENTID")
	if !ret {
		options.Id = c.Flag("clientid").Value.String()
		l.Debugf("Using flag value for client ID")
	}

	options.Qos, err = strconv.Atoi(c.Flag("qos").Value.String())
	if err != nil {
		l.Errorf("Unable to parse QoS value")
	}

	messages.DbConn, ret = os.LookupEnv("MQTT_DB")
	messages.Logger = l
	if !ret {
		messages.DbConn = c.Flag("db").Value.String()
		l.Debugf("Using flag value for DB Connection")
	}

	dbconn, err := sql.Open("sqlite3", messages.DbConn)
	if err != nil {
		panic(err)
	}
	defer dbconn.Close()

	if _, table_check := dbconn.Query("select * from topics;"); table_check == nil {
		l.Debugf("Table 'topics' already exists")
	} else {
		if _, err := dbconn.Exec(create); err != nil {
			l.Errorf(err.Error())
		} else {
			l.Debugf("Creating new DB Table")
		}
	}

	mqtt.Connect(&options, channel, l)
}
