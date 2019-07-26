package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"

	"github.com/ChimeraCoder/anaconda"
	"github.com/dghubble/oauth1"
)

var (
	APIKEY       = os.Getenv("APIKEY")
	APISECRET    = os.Getenv("APISECRET")
	ACCESSTOKEN  = os.Getenv("ACCESSTOKEN")
	ACCESSSECRET = os.Getenv("ACCESSSECRET")
	WEBHOOKENV   = "dev3"
	APPURL       = "https://dfde80a7.ngrok.io"
)

var api *anaconda.TwitterApi

type WebhookLoad struct {
	UserId           string  `json:"for_user_id"`
	TweetCreateEvent []Tweet `json:"tweet_create_events"`
}

type Tweet struct {
	Id       int64
	IdStr    string `json:"id_str"`
	User     User
	Text     string
	ParentID int64  `json:"in_reply_to_status_id"`
	Entities Entity `json:"entities"`
}

type Entity struct {
	Media []Media `json:"media"`
}

type Media struct {
	ID   int64  `json:"id"`
	URL  string `json:"media_url_https"`
	Type string `json:"type"`
}

type User struct {
	Id     int64
	IdStr  string `json:"id_str"`
	Name   string
	Handle string `json:"screen_name"`
}

func init() {
	anaconda.SetConsumerKey(APIKEY)
	anaconda.SetConsumerSecret(APISECRET)
	api = anaconda.NewTwitterApi(ACCESSTOKEN, ACCESSSECRET)
}

func postTweet(tweet string) {
	v := url.Values{}
	api.PostTweet(tweet, v)
}

func createClient() *http.Client {
	config := oauth1.NewConfig(APIKEY, APISECRET)
	token := oauth1.NewToken(ACCESSTOKEN, ACCESSSECRET)
	return config.Client(oauth1.NoContext, token)
}

func registerWebhook() {
	fmt.Println("register webhook")
	httpClient := createClient()

	path := "https://api.twitter.com/1.1/account_activity/all/" + WEBHOOKENV + "/webhooks.json"
	values := url.Values{}
	values.Set("url", APPURL+"/webhook/twitter")

	resp, err := httpClient.PostForm(path, values)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	var data map[string]interface{}
	if err := json.Unmarshal([]byte(body), &data); err != nil {
		panic(err)
	}
	fmt.Println("registerWebHook data: ", data)
	fmt.Println("Webhook id of " + data["id"].(string) + " has been registered")
	fmt.Println("Before calling subscribeWebhook")
	subscribeWebhook()
}

func subscribeWebhook() {
	fmt.Println("Into subscribe webhook")
	client := createClient()
	path := "https://api.twitter.com/1.1/account_activity/all/" + WEBHOOKENV + "/subscriptions.json"
	resp, err := client.PostForm(path, nil)
	if err != nil {
		panic(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	if resp.StatusCode == 204 {
		fmt.Println("subscribed to WH: ", string(body))
	} else {
		fmt.Println("Unable to subscribe to WH: ", string(body))
	}
}

func SendTweet(tweet string, reply_id string) (*Tweet, error) {
	fmt.Println("Sending tweet as reply to " + reply_id)
	//Initialize tweet object to store response in
	var responseTweet Tweet
	//Add params
	params := url.Values{}
	params.Set("status", tweet)
	params.Set("in_reply_to_status_id", reply_id)
	//Grab client and post
	client := createClient()
	resp, err := client.PostForm("https://api.twitter.com/1.1/statuses/update.json", params)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	//Decode response and send out
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
	err = json.Unmarshal(body, &responseTweet)
	if err != nil {
		return nil, err
	}
	return &responseTweet, nil
}
