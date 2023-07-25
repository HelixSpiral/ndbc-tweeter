package main

import (
	"log"
	"os"
	"strconv"

	"github.com/helixspiral/ndbc"
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

	buoyInfo, err := n.GetLatestDataFromBuoy(buoyID)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(buoyPicture)
	log.Println(buoyInfo)
}
