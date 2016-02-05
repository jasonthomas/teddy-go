package main

import (
	"encoding/json"
	"fmt"
	// irc "github.com/fluffle/goirc/client"
	"os"
)

type IRCChannels struct {
	Name string
	Key  string
}

type IRCNetworks struct {
	Host     string
	Port     int
	Ssl      bool
	Channels map[string]IRCChannels
}

type IRCConfig struct {
	Nick     string
	Networks map[string]IRCNetworks
}

func main() {
	config := readConfig("config.json")
	fmt.Println(config.Networks["test"].Channels)
}

func readConfig(configFile string) IRCConfig {
	config := IRCConfig{}

	file, _ := os.Open(configFile)
	decoder := json.NewDecoder(file)
	err := decoder.Decode(&config)

	if err != nil {
		fmt.Println("error:", err)
	}

	return config
}
