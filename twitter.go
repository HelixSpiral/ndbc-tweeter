package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"os"

	"github.com/dghubble/oauth1"
)

func uploadImage(image []byte) (twitterMediaResponse, error) {
	consumerKey := os.Getenv("CONSUMER_KEY")
	consumerSecret := os.Getenv("CONSUMER_SECRET")
	accessToken := os.Getenv("ACCESS_TOKEN")
	accessSecret := os.Getenv("ACCESS_SECRET")

	config := oauth1.NewConfig(consumerKey, consumerSecret)
	token := oauth1.NewToken(accessToken, accessSecret)

	httpClient := config.Client(oauth1.NoContext, token)

	imageReader := bytes.NewReader(image)
	imageBuf := &bytes.Buffer{}
	form := multipart.NewWriter(imageBuf)

	fw, err := form.CreateFormFile("media", "buoyPicture.jpg")
	if err != nil {
		return twitterMediaResponse{}, fmt.Errorf("error creating file: %w", err)
	}

	_, err = io.Copy(fw, imageReader)
	if err != nil {
		return twitterMediaResponse{}, fmt.Errorf("error in io copy: %w", err)
	}
	form.Close()

	resp, err := httpClient.Post("https://upload.twitter.com/1.1/media/upload.json?media_category=tweet_image", form.FormDataContentType(), bytes.NewReader(imageBuf.Bytes()))
	if err != nil {
		return twitterMediaResponse{}, fmt.Errorf("error in http POST: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return twitterMediaResponse{}, fmt.Errorf("error reading http body: %w", err)
	}

	var mediaResp twitterMediaResponse

	err = json.Unmarshal(body, &mediaResp)
	if err != nil {
		return twitterMediaResponse{}, fmt.Errorf("error unmarshaling http body to Json: %w", err)
	}

	return mediaResp, nil
}

func sendMessage(message, mediaID string) error {
	consumerKey := os.Getenv("CONSUMER_KEY")
	consumerSecret := os.Getenv("CONSUMER_SECRET")
	accessToken := os.Getenv("ACCESS_TOKEN")
	accessSecret := os.Getenv("ACCESS_SECRET")

	config := oauth1.NewConfig(consumerKey, consumerSecret)
	token := oauth1.NewToken(accessToken, accessSecret)

	httpClient := config.Client(oauth1.NoContext, token)

	log.Printf("Tweeting message: %s\r\n", message)

	resp, err := httpClient.Post("https://api.twitter.com/2/tweets", "application/json",
		bytes.NewBuffer([]byte(fmt.Sprintf(`{"text": "%s", "media": {"media_ids": ["%s"]}}`, message, mediaID))))
	if err != nil {
		return fmt.Errorf("error in http POST: %w", err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading http body: %w", err)
	}

	log.Println("Tweet:", string(body))

	return nil
}
