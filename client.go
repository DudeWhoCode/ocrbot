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
	"github.com/pkg/errors"
)

var (
	apiKey               = os.Getenv("APIKEY")
	apiSecret            = os.Getenv("APISECRET")
	accessToken          = os.Getenv("ACCESSTOKEN")
	accessSecret         = os.Getenv("ACCESSSECRET")
	webhookEndpoint      = "https://api.twitter.com/1.1/account_activity/all/%s/webhooks.json"
	subscriptionEndpoint = "https://api.twitter.com/1.1/account_activity/all/%s/subscriptions.json"
)

var api *anaconda.TwitterApi

func init() {
	anaconda.SetConsumerKey(apiKey)
	anaconda.SetConsumerSecret(apiSecret)
	api = anaconda.NewTwitterApi(accessToken, accessSecret)
}

func replyTweet(tweet string, replyID string) error {
	v := url.Values{}
	v.Set("in_reply_to_status_id", replyID)
	_, err := api.PostTweet(tweet, v)
	if err != nil {
		return errors.Wrap(err, "reply tweet failed")
	}

	return nil
}

func createClient() *http.Client {
	config := oauth1.NewConfig(apiKey, apiSecret)
	token := oauth1.NewToken(accessToken, accessSecret)
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
