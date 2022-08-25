package main

import (
	"encoding/json"
	"flag"
	"log"
	"os"
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

func read_conf() Config {
	ConfPath := flag.String("config", "config.json", "Path to config")
	flag.Parse()

	data, err := os.ReadFile(*ConfPath)
	if err != nil {
		log.Fatal(err)
	}

	var Conf Config
	json.Unmarshal(data, &Conf)

	return Conf
}
