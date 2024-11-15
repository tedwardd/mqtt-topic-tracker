package mqtt

import (
	"crypto/tls"
	"log"
	"os"
	"os/signal"
	"syscall"

	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/tedwardd/mqtt-topic-tracker/internal/logger"
	"github.com/tedwardd/mqtt-topic-tracker/internal/messages"
)

type Options struct {
	Server   string
	Username string
	Password string
	Id       string
	Topic    string
	Qos      int
	Db       string
}

func Connect(options *Options, channel chan os.Signal, logger *logger.LogCache) {
	signal.Notify(channel, os.Interrupt, syscall.SIGTERM)

	connOpts := MQTT.NewClientOptions().AddBroker(options.Server).SetClientID(options.Id).SetCleanSession(true)
	connOpts.SetUsername(options.Username)
	connOpts.SetPassword(options.Password)

	tlsConfig := &tls.Config{InsecureSkipVerify: true, ClientAuth: tls.NoClientCert}

	connOpts.SetTLSConfig(tlsConfig)

	connOpts.OnConnect = func(channel MQTT.Client) {
		if token := channel.Subscribe(options.Topic, byte(options.Qos), messages.OnMessageReceived); token.Wait() && token.Error() != nil {
			log.Fatal(token.Error())
		}
	}

	client := MQTT.NewClient(connOpts)
	logger.Printf("Connecting to Broker...")
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatal(token.Error())
	}
	logger.Printlnf("Connected")

	<-channel
}
