package main

import (
	"encoding/json"
	"flag"
	"log"
	"os"
)

var (
	ConfPath = flag.String("config", "config.json", "Path to config")
	DBPath   = flag.String("db", "limits.db", "Path to database")
	LogPath  = flag.String("log", "mastodon-group-bot.log", "Path to log")
)

type Config struct {
	Server         string   `json:"Server"`
	ClientID       string   `json:"ClientID"`
	ClientSecret   string   `json:"ClientSecret"`
	AccessToken    string   `json:"AccessToken"`
	WelcomeMessage string   `json:"WelcomeMessage"`
	Max_toots      uint16   `json:"Max_toots"`
	Toots_interval uint16   `json:"Toots_interval"`
	Admins         []string `json:"Admins"`
}

func ReadConf() Config {
	flag.Parse()

	data, err := os.ReadFile(*ConfPath)
	if err != nil {
		log.Fatal("Failed to read config")
	}

	var Conf Config
	json.Unmarshal(data, &Conf)

	return Conf
}
