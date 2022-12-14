package opts

import (
	"encoding/json"
	"log"
)

type Params struct {
	FrigateServer    string `           long:"server"      env:"FRIGATE_SERVER"  description:"Server name or IP of frigate server and port number" default:"http://192.168.4.1:5000"`
	FrigateMqtt      string `           long:"mqtt"      env:"FRIGATE_MQTT_SERVER"  description:"Server name or IP of mqtt server and port number" default:"192.168.4.1:1883"`
	DiscordToken     string `           long:"token"      env:"DISCORD_TOKEN"  description:"Discord Webhook token"`
	DiscordChannelID string `           long:"channel"      env:"DISCORD_CHANNEL_ID"  description:"Discord Channel ID"`
	SnapshotQuality  string `           long:"pixel"      env:"SNAPSHOT_PIXEL"  description:"Snapshot image size" default:"300"`
	TimeZone         string `           long:"timezone"      env:"TIME_ZONE"  description:"Timezone as per timezone db" default:"Asia/Kolkata"`
}

func (o *Params) GetJson() []byte {
	jsonBytes, err := json.Marshal(o)
	if err != nil {
		log.Panic(err)
	}
	return jsonBytes
}
