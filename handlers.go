package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/gorilla/mux"
)

func startServer() {
	m := mux.NewRouter()
	m.HandleFunc("/webhook/twitter", crcCheck).Methods("GET")
	m.HandleFunc("/webhookdev/twitter", crcCheck).Methods("GET") // Dev environment
	m.HandleFunc("/webhook/twitter", webhookHandler).Methods("POST")
	m.HandleFunc("/webhookdev/twitter", webhookHandler).Methods("POST") // Dev environment
	m.HandleFunc("/register-webhook", registerNewWebhook).Methods("POST")
	m.HandleFunc("/ping", ping).Methods("GET")

	server := &http.Server{
		Handler: m,
	}
	server.Addr = "0.0.0.0:8000"
	log.Info("Starting server...")
	server.ListenAndServe()
}

func webhookHandler(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)
	var load Webhook
	err := json.Unmarshal(body, &load)
	if err != nil {
		log.Error(err)
		return
	}

	//Check if it was a tweet_create_event and tweet was in the payload and it was not tweeted by the bot
	if len(load.TweetCreateEvent) < 1 || load.UserId == load.TweetCreateEvent[0].User.IdStr {
		log.Warn("Tweet create event is not a mention")
		return
	}
	replyHandle := load.TweetCreateEvent[0].User.Handle
	log.Infof("Got mentioned by %s", replyHandle)

	parentID := load.TweetCreateEvent[0].ParentID
	if parentID == 0 {
		log.Error("Unable to find parent tweet ID")
		return
	}
	v := url.Values{}
	parentTweet, err := api.GetTweet(parentID, v)
	if err != nil {
		log.Error("Unable to get parent tweet")
		return
	}
	if len(parentTweet.Entities.Media) == 0 {
		log.Error("Parent tweet has no media")
		return
	}

	// TODO: handle multiple media files
	media := parentTweet.Entities.Media[0]
	if media.Type != "photo" {
		return
	}

	// TODO: use dynamic file names for saving images
	err = downloadImage(media.Media_url_https, "pic.jpg")
	if err != nil {
		log.Errorf("Unable to download image. \n %s", err)
		return
	}

	text, err := read("pic.jpg")
	if err != nil {
		log.Errorf("Unable to read image. \n %s", err)
		return
	}

	// TODO: If pastebin API limit is reached, retry here. Or add queue for creating pastes
	pasteURL, err := createPaste(text)
	if err != nil {
		log.Errorf("Unable to create paste. \n %s", err)
		return
	}

	err = replyTweet("@"+replyHandle+" Here is the text: "+pasteURL, load.TweetCreateEvent[0].IdStr)
	if err != nil {
		log.Errorf("Error while replying to %s \n %s", replyHandle, err) // Log tweet URL instead of just handle
	} else {
		log.Infof("Reply sent successfully to %s", replyHandle)
	}

}

func crcCheck(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Into crcCheck")
	w.Header().Set("Content-Type", "application/json")

	token := r.URL.Query()["crc_token"]
	if len(token) < 1 {
		fmt.Fprintf(w, "No crc_token given")
	}

	h := hmac.New(sha256.New, []byte(apiSecret))
	h.Write([]byte(token[0]))
	encoded := base64.StdEncoding.EncodeToString(h.Sum(nil))

	response := make(map[string]string)
	response["response_token"] = "sha256=" + encoded

	responseJson, _ := json.Marshal(response)
	fmt.Fprintf(w, string(responseJson))
}

func ping(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "pong")
}

func registerNewWebhook(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Into registerNewWebhook")
	var payload RegisterWebHook
	decoder := json.NewDecoder(r.Body)
	decoder.Decode(&payload)
	webhookID, err := registerWebhook(payload.EnvName, payload.AppURL+payload.WebhookPath)
	if err != nil {
		log.Errorf("Unable to register webhook.\n %s", err)
	}
	log.Infof("subscribed to webhook: %s", webhookID)
}
