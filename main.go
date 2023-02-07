package main

import (
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/jessevdk/go-flags"
	"github.com/vibin18/go-mqtt-discord/internal/handlers"
	"github.com/vibin18/go-mqtt-discord/internal/opts"
	"github.com/vibin18/go-mqtt-discord/internal/repos"
	"log"
	"os"
	"time"
)

var (
	argparser *flags.Parser
	arg       opts.Params
)

func initArgparser() {
	argparser = flags.NewParser(&arg, flags.Default)
	_, err := argparser.Parse()

	// check if there is an parse error
	if err != nil {
		if flagsErr, ok := err.(*flags.Error); ok && flagsErr.Type == flags.ErrHelp {
			os.Exit(0)
		} else {
			fmt.Println()
			argparser.WriteHelp(os.Stdout)
			os.Exit(1)
		}
	}
}

func main() {

	initArgparser()

	var params repos.Repository
	params.Params = &arg
	handlers.NewConfig(&params)

	mqttClient := mqtt.NewClientOptions()
	mqttClient.ConnectTimeout = 30 * time.Second
	mqttClient.ConnectRetry = true
	mqttClient.AutoReconnect = true
	mqttClient.KeepAlive = 25
	mqttClient.CleanSession = true
	mqttClient.ConnectRetryInterval = 20 * time.Second
	mqttClient.PingTimeout = 30 * time.Second
	mqttClient.MaxReconnectInterval = 30 * time.Second
	mqttClient.ResumeSubs = true
	mqttClient.AddBroker(fmt.Sprintf("tcp://%v", arg.FrigateMqtt))
	mqttClient.SetClientID("go_mqtt_client")
	mqttClient.OnConnect = handlers.ConnectHandler
	mqttClient.OnConnectionLost = handlers.ConnectLostHandler
	client := mqtt.NewClient(mqttClient)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	go func() {
		timeout := 10

		ok := false
		for !ok {
			if timeout <= 0 {
				log.Println("Connection not ready after timeout, exiting..")
				return
			}
			ok = client.IsConnectionOpen()
			if !ok {
				log.Println("Connection not ready")
				time.Sleep(30 * time.Second)
				timeout--
			}
			log.Println("Connection retrying..")

		}

		handlers.Sub(client, "frigate/events")

	}()

	select {}
	client.Disconnect(250)
}
