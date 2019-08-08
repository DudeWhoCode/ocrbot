package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/ChimeraCoder/anaconda"
	"github.com/dghubble/oauth1"
)

const (
	WEBHOOKENV = "dev4"
	APPURL     = "https://ocrbot.dudewho.codes"
)

var (
	APIKEY               = os.Getenv("APIKEY")
	APISECRET            = os.Getenv("APISECRET")
	ACCESSTOKEN          = os.Getenv("ACCESSTOKEN")
	ACCESSSECRET         = os.Getenv("ACCESSSECRET")
	webhookEndpoint      = "https://api.twitter.com/1.1/account_activity/all/%s/webhooks.json"
	subscriptionEndpoint = "https://api.twitter.com/1.1/account_activity/all/%s/subscriptions.json"
)

var api *anaconda.TwitterApi

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

func registerWebhook(envName, webhookPath string) {
	fmt.Println("Into register webhook")
	httpClient := createClient()

	path := fmt.Sprintf(webhookEndpoint, envName)
	values := url.Values{}
	values.Set("url", webhookPath)

	resp, err := httpClient.PostForm(path, values)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var data map[string]interface{}
	if err := json.Unmarshal([]byte(body), &data); err != nil {
		log.Fatal(err)
	}
	log.Println("registerWebHook data: ", data)
	log.Println("Webhook id of " + data["id"].(string) + " has been registered")
	subscribeWebhook(envName)
}

func subscribeWebhook(envName string) {
	fmt.Println("Into subscribe webhook")
	client := createClient()
	path := fmt.Sprintf(subscriptionEndpoint, envName)
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
