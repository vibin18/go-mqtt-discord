package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/bwmarrin/discordgo"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/vibin18/go-mqtt-discord/internal/models"
	"github.com/vibin18/go-mqtt-discord/internal/repos"
	"log"
	"net/http"
	"strings"
)

var config *repos.Repository

func NewConfig(c *repos.Repository) {
	config = c
}

var MessagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {

	bot, err := discordgo.New("Bot " + config.Params.DiscordToken)
	if err != nil {
		log.Panicf("failed to create discord client %v", err)
		return
	}
	var events models.Events

	err = json.Unmarshal(msg.Payload(), &events)
	if err != nil {
		log.Println("error unmarshalling")
	}

	if events.Type == "new" {
		var eventStartime float64
		var snapShotURL strings.Builder

		eventID := events.Before.ID
		label := events.Before.Label
		camera := events.Before.Camera
		eventStartime = events.Before.StartTime
		startTime := int64(eventStartime)

		snapShotURL.WriteString(config.Params.FrigateServer)
		snapShotURL.WriteString("/api/")
		snapShotURL.WriteString(camera)
		snapShotURL.WriteString("/latest.jpg?h=")
		snapShotURL.WriteString(config.Params.SnapshotQuality)

		response, err := http.Get(snapShotURL.String())
		if err != nil {
			log.Println(err)
			return
		}
		defer response.Body.Close()

		var files []*discordgo.File

		NewFile := discordgo.File{
			Name:        fmt.Sprintf("%v.jpeg", eventID),
			Reader:      response.Body,
			ContentType: "image/jpeg",
		}

		files = append(files, &NewFile)

		mc := discordgo.MessageSend{
			Content: fmt.Sprintf("New %v detetced on %v at %v", label, camera, startTime),
			Files:   files,
		}

		messageContent := discordgo.MessageSend(mc)

		messageSend, error := bot.ChannelMessageSendComplex(config.Params.DiscordChannelID, &messageContent)
		if err != nil {
			log.Println(error)
			return
		}
		log.Printf("%v", messageSend.Content)

	}

}

var ConnectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	log.Println("connected to mqtt")
}

var ConnectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	log.Printf("connection lost to mqtt: %v\n", err)
}

func Sub(client mqtt.Client, topic string) {
	token := client.Subscribe(topic, 1, MessagePubHandler)
	token.Wait()
	log.Printf("subscribed to topic: %s\n", topic)
}
