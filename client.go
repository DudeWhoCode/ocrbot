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

// registerWebhook registers the webhook path to twitter environment and returns respective webhookID
func registerWebhook(envName, webhookPath string) (string, error) {
	httpClient := createClient()
	path := fmt.Sprintf(webhookEndpoint, envName)
	values := url.Values{}
	values.Set("url", webhookPath)

	// TODO: Check whether you can achieve the same using anaconda
	resp, err := httpClient.PostForm(path, values)
	if err != nil {
		return "", errors.Wrap(err, "register webhook POST failed")
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", errors.Wrap(err, "read failed")
	}

	var data map[string]interface{}
	if err := json.Unmarshal([]byte(body), &data); err != nil {
		return "", errors.Wrap(err, "json unmarshall failed")
	}
	log.Infof("Registered Webhook: %s", data)
	webhookID, err := subscribeWebhook(envName)
	if err != nil {
		return "", errors.Wrap(err, "subscribe webhook failed")
	}

	return webhookID, nil
}

// subscribeWebhook subscribes this endpoint to events
func subscribeWebhook(envName string) (string, error) {
	client := createClient()
	path := fmt.Sprintf(subscriptionEndpoint, envName)

	// TODO: Check whether you can achieve the same using anaconda
	resp, err := client.PostForm(path, nil)
	if err != nil {
		return "", errors.Wrap(err, "webhook post call failed")
	}
	if resp.StatusCode != 204 {
		return "", errors.New("204 status code not received")
	}

	webhookID, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", errors.Wrap(err, "webhook read failed")
	}

	return string(webhookID), nil
}
