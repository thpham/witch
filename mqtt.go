package main

import (
	"encoding/json"
	"time"
	"os"

	paho "github.com/eclipse/paho.mqtt.golang"

	"github.com/Eagle-X/witch/system"
)

// Client is the MQTT client.
type Client struct {
	broker string
	client paho.Client
	cfg *Config
	control *system.Controller
}

//define a function for the default message handler
var defaultHandler paho.MessageHandler = func(client paho.Client, msg paho.Message) {
	log.Printf("TOPIC: %s MSG: %s\n", msg.Topic(), msg.Payload())
}

func handleActionMessage(client paho.Client, msg paho.Message, control *system.Controller) {
	log.Printf("Message action: %s", msg.Payload())
	action := &system.Action{}
	if err := json.Unmarshal(msg.Payload(), action); err != nil {
		log.Printf("Invalid action format: %s", err)
		return
	}
	control.Handle(action)
}


// MqttClient inits a MQTT client.
func MqttClient(broker string, control *system.Controller, cfg *Config) *Client {
	//create a ClientOptions struct setting the broker address, clientid, turn
	//off trace output and set the default message handler
	opts := paho.NewClientOptions().AddBroker(broker)
	opts.SetClientID(cfg.Mqtt.ClientID)
	opts.SetProtocolVersion(4)
	opts.SetUsername(cfg.Mqtt.Username)
	opts.SetPassword(cfg.Mqtt.Password)
	opts.SetKeepAlive(time.Duration(cfg.Mqtt.Keepalive))
	opts.SetDefaultPublishHandler(defaultHandler)

	// opts.SetWill("my/will/topic", "Goodbye", 1, true)

	//create and start a client using the above ClientOptions
	cl := &Client{
		broker: broker,
		client: paho.NewClient(opts),
		cfg: cfg,
		control: control,
	}

	return cl
}

// Start starts the client.
func (cl *Client) Start() error {
	log.Printf("MQTT client connects to %s", cl.broker)
	if token := cl.client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	token := cl.client.Subscribe(cl.cfg.Mqtt.Actions_message.Topic, cl.cfg.Mqtt.Actions_message.Qos, func(client paho.Client, msg paho.Message) {
		handleActionMessage(client, msg, cl.control)
	})
	//subscribe to the topic and request messages to be delivered,
	// wait for the receipt to confirm the subscription
	if token.Wait() && token.Error() != nil {
		log.Println(token.Error())
		os.Exit(1)
	}
	return nil
}

// Stop stops the client.
func (cl *Client) Stop() {
	cl.client.Disconnect(250)
}
