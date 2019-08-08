package main

type Webhook struct {
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

type RegisterWebHook struct {
	EnvName     string `json:"env_name"`
	AppURL      string `json:"app_url"`
	WebhookPath string `json:"webhook_path"`
}
