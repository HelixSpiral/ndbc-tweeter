package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"unicode"

	"github.com/helixspiral/ndbc"
	"golang.org/x/exp/slices"
)

func main() {
	// Some initial setup
	buoyIDString := os.Getenv("BUOY_ID")
	buoyID, err := strconv.Atoi(buoyIDString)
	if err != nil {
		log.Fatalf("Error converting BUOY_ID to int")
	}

	n := ndbc.NewAPI()

	buoyPicture, err := n.GetPictureFromBuoy(buoyID)
	if err != nil {
		log.Fatal(err)
	}
	/*
		buoyPictureFile, err := os.Create("buoyPicture.jpg")
		if err != nil {
			log.Fatal(err)
		}
		defer buoyPictureFile.Close()

		pictureIOReader := bytes.NewReader(buoyPicture)

		_, err = io.Copy(buoyPictureFile, pictureIOReader)
		if err != nil {
			log.Fatal(err)
		}
	*/
	buoyInfo, err := n.GetLatestDataFromBuoy(buoyID)
	if err != nil {
		log.Fatal(err)
	}

	// Upload image to Twitter
	mediaResp, err := uploadImage(buoyPicture)
	if err != nil {
		log.Fatal(err)
	}

	// Send tweet
	var tweetMessage string
	if slices.Contains([]string{
		"North", "Northeast",
		"East", "Southeast",
		"South", "Southwest",
		"West", "Northwest",
	}, buoyInfo.WindDirection) {
		tweetMessage += fmt.Sprintf("Winds coming in from the %s, ", buoyInfo.WindDirection)
	}

	if buoyInfo.WindSpeed > 0 {
		tweetMessage += fmt.Sprintf("sustained winds of %f", buoyInfo.WindSpeed)
		if buoyInfo.GustSpeed > 0 {
			tweetMessage += fmt.Sprintf(", and gusting up to %f!", buoyInfo.GustSpeed)
		} else {
			tweetMessage += "."
		}
	}

	// I'm sure there's a better way to do this, but it's simple and it works.
	// We use this to ensure the first letter is capital, there's a change it wont be if there's no
	// wind direction information
	tweetMessageRunes := []rune(tweetMessage)
	tweetMessageRunes[0] = unicode.ToUpper(tweetMessageRunes[0])
	tweetMessage = string(tweetMessageRunes)

	err = sendMessage(tweetMessage, mediaResp.MediaKey)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(buoyPicture)
	log.Println(buoyInfo)
}
