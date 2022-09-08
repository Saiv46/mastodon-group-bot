package main

import (
	"encoding/json"
	"flag"
	"log"
	"os"
)

var (
	ConfPath = flag.String("config", "config.json", "Path to config")
	DBPath   = flag.String("db", "mastodon-group-bot.db", "Path to database")
	LogPath  = flag.String("log", "mastodon-group-bot.log", "Path to log")

	Conf = ReadConfig()
)

type Config struct {
	Server               string   `json:"Server"`
	ClientID             string   `json:"ClientID"`
	ClientSecret         string   `json:"ClientSecret"`
	AccessToken          string   `json:"AccessToken"`
	WelcomeMessage       string   `json:"WelcomeMessage"`
	NotFollowedMessage   string   `json:"NotFollowedMessage"`
	Max_toots            uint     `json:"Max_toots"`
	Toots_interval       uint     `json:"Toots_interval"`
	Duplicate_buf        uint     `json:"Duplicate_buf"`
	Order_limit          uint     `json:"Order_limit"`
	Del_notices_interval uint     `json:"Del_notices_interval"`
	Admins               []string `json:"Admins"`
}

func ReadConfig() Config {
	flag.Parse()

	data, err := os.ReadFile(*ConfPath)
	if err != nil {
		log.Fatal("Failed to read config")
	}

	var Conf Config
	json.Unmarshal(data, &Conf)

	return Conf
}
