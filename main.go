package main

import (
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/jessevdk/go-flags"
	"github.com/vibin18/go-mqtt-discord/internal/handlers"
	"github.com/vibin18/go-mqtt-discord/internal/ops"
	"github.com/vibin18/go-mqtt-discord/internal/repos"
	"os"
)

var (
	argparser *flags.Parser
	arg       ops.Params
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

	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%v", arg.FrigateMqtt))
	opts.SetClientID("go_mqtt_client")
	//opts.SetUsername("emqx")
	//opts.SetPassword("public")
	//opts.SetDefaultPublishHandler(messagePubHandler)
	opts.OnConnect = handlers.ConnectHandler
	opts.OnConnectionLost = handlers.ConnectLostHandler
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	go func() {
		handlers.Sub(client, "frigate/events")

	}()
	select {}
	client.Disconnect(250)
}
