package main

type twitterMediaResponse struct {
	MediaID       int              `json:"media_id"`
	MediaIDString string           `json:"media_id_string"`
	MediaKey      string           `json:"media_key"`
	Size          int              `json:"size"`
	ExpiresAfter  int              `json:"expires_after_secs"`
	Image         twitterImageInfo `json:"image"`
}

type twitterImageInfo struct {
	Type   string `json:"image_type"`
	Width  int    `json:"w"`
	Height int    `json:"h"`
}
