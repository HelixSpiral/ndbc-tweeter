package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
	"unicode"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/helixspiral/ndbc"
	"golang.org/x/exp/slices"
)

func main() {
	// Some initial Twitter setup
	consumerKey := os.Getenv("CONSUMER_KEY")
	consumerSecret := os.Getenv("CONSUMER_SECRET")
	accessToken := os.Getenv("ACCESS_TOKEN")
	accessSecret := os.Getenv("ACCESS_SECRET")

	// Some initial Mastodon setup
	mastodonServer := os.Getenv("MASTODON_SERVER")
	mastodonClientID := os.Getenv("MASTODON_CLIENT_ID")
	mastodonClientSecret := os.Getenv("MASTODON_CLIENT_SECRET")
	mastodonUser := os.Getenv("MASTODON_USERNAME")
	mastodonPass := os.Getenv("MASTODON_PASSWORD")

	// Some initial MQTT setup
	mqttBroker := os.Getenv("MQTT_BROKER")
	mqttClientId := os.Getenv("MQTT_CLIENT_ID")
	mqttTopic := os.Getenv("MQTT_TOPIC")

	options := mqtt.NewClientOptions().AddBroker(mqttBroker).SetClientID(mqttClientId)
	options.WriteTimeout = 20 * time.Second
	mqttClient := mqtt.NewClient(options)

	// Some initial setup
	buoyIDString := os.Getenv("BUOY_ID")
	buoyLocation := os.Getenv("BUOY_LOCATION") // There's no clean and easy way to get this from the API, since these buoys don't move we'll just throw it in this way.

	buoyID, err := strconv.Atoi(buoyIDString)
	if err != nil {
		log.Fatalf("Error converting BUOY_ID to int")
	}

	n := ndbc.NewAPI()

	buoyPicture, err := n.GetPictureFromBuoy(buoyID)
	if err != nil {
		log.Fatal(err)
	}

	buoyInfo, err := n.GetLatestDataFromBuoy(buoyID)
	if err != nil {
		log.Fatal(err)
	}

	// Setup the MQTT message
	var message string
	if buoyLocation != "" {
		message += fmt.Sprintf("%s: ", buoyLocation)
	}

	// Figure out wind direction
	if slices.Contains([]string{
		"North", "Northeast",
		"East", "Southeast",
		"South", "Southwest",
		"West", "Northwest",
	}, buoyInfo.WindDirection) {
		message += fmt.Sprintf("Winds coming in from the %s, ", buoyInfo.WindDirection)
	}

	// Add wind speed
	if buoyInfo.WindSpeed > 0 {
		message += fmt.Sprintf("sustained winds of %.1f m/s", buoyInfo.WindSpeed)
		if buoyInfo.GustSpeed > 0 && buoyInfo.GustSpeed != buoyInfo.WindSpeed {
			message += fmt.Sprintf(", and gusting up to %.1f m/s!", buoyInfo.GustSpeed)
		} else {
			message += "!"
		}
	}

	// Some newlines and then hashtags
	if message != "" {
		message += "\\r\\n\\r\\n"
	}

	message += "#noaa #ndbc #buoy #Maine #coast #weather #ocean"

	// I'm sure there's a better way to do this, but it's simple and it works.
	// We use this to ensure the first letter is capital, there's a change it wont be if there's no
	// wind direction information
	messageRunes := []rune(message)
	messageRunes[0] = unicode.ToUpper(messageRunes[0])
	message = string(messageRunes)

	jsonMsg, err := json.Marshal(&MqttMessage{
		MastodonServer:       mastodonServer,
		MastodonClientID:     mastodonClientID,
		MastodonClientSecret: mastodonClientSecret,
		MastodonUser:         mastodonUser,
		MastodonPass:         mastodonPass,

		TwitterConsumerKey:    consumerKey,
		TwitterConsumerSecret: consumerSecret,
		TwitterAccessToken:    accessToken,
		TwitterAccessSecret:   accessSecret,

		Message: message,
		Image:   buoyPicture,
	})
	if err != nil {
		log.Fatal(err)
	}

	// Connect to the MQTT broker
	if token := mqttClient.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	log.Println("Sending message:", message)
	token := mqttClient.Publish(mqttTopic, 2, false, jsonMsg)
	_ = token.Wait()
	if token.Error() != nil {
		panic(err)
	}
}
