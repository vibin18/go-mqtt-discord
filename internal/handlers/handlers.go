package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/bwmarrin/discordgo"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/vibin18/go-mqtt-discord/internal/models"
	"github.com/vibin18/go-mqtt-discord/internal/repos"
	"io"
	"log"
	"net/http"
)

type MessageContent discordgo.MessageSend

var config *repos.Repository

func NewConfig(c *repos.Repository) {
	config = c
}

func (c *MessageContent) NewMessageContent() *discordgo.MessageSend {
	return &discordgo.MessageSend{}

}

func (c *MessageContent) SetName(content string) {
	c.Content = content
}

func (c *MessageContent) SetEmbed(eventID string) {
	var embedImageList []discordgo.MessageEmbedImage
	var embedList []*discordgo.MessageEmbed

	embedImage := discordgo.MessageEmbedImage{

		URL: config.Params.FrigateServer + "/events/" + eventID + "/snapshot.jpg",
	}
	embedImageList = append(embedImageList, embedImage)

	embedNew := discordgo.MessageEmbed{
		Image: &embedImage,
		Type:  "image",
	}

	embedList = append(embedList, &embedNew)

	c.Embeds = embedList

}

func getImage(eventID string) (io.Reader, error) {

	res, error := http.Get(config.Params.FrigateServer + "/events/" + eventID + "/snapshot.jpg")
	if error != nil {
		return nil, error
	}
	if res.ContentLength >= 0 {
		return res.Body, nil
	}
	defer res.Body.Close()
	return nil, nil

}

var MessagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {

	bot, err := discordgo.New("Bot " + config.Params.DiscordToken)
	if err != nil {
		log.Panicf(err.Error())
		return
	}
	var events models.Events

	err = json.Unmarshal(msg.Payload(), &events)
	if err != nil {
		log.Println("Error unmarshalling")
	}

	var content MessageContent

	content.NewMessageContent()

	if events.Type == "new" {
		fmt.Printf("New event alerts\n")
		eventID := events.Before.ID
		label := events.Before.Label
		camera := events.Before.Camera

		imageStatus, err := getImage(eventID)
		if err != nil {
			log.Println("Image status is: " + fmt.Sprint(imageStatus))
		}
		if imageStatus != nil {
			content.SetName(fmt.Sprintf("New %v detetced on %v", label, camera))
			content.SetEmbed(eventID)

		} else {
			content.SetName(fmt.Sprintf("A %v is detetced on %v, But image not found for event with ID: %v", label, camera, eventID))
		}

		messageContent := discordgo.MessageSend(content)

		//messageSend, err := bot.ChannelMessageSend(discordChannelID, "New movement detetced")
		messageSend, error := bot.ChannelMessageSendComplex(config.Params.DiscordChannelID, &messageContent)
		if err != nil {
			log.Println(error)
			return
		}
		log.Printf(": ==> %v", messageSend.Content)
		//	fmt.Printf("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())
	}

}

var ConnectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	fmt.Println("Connected")
}

var ConnectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	fmt.Printf("Connect lost: %v\n", err)
}

func Sub(client mqtt.Client, topic string) {
	token := client.Subscribe(topic, 1, MessagePubHandler)
	token.Wait()
	fmt.Printf("Subscribed to topic: %s\n", topic)
}
