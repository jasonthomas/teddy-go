package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	irc "github.com/fluffle/goirc/client"
	"log"
	"os"
	"strings"
)

type JenkinsCreds struct {
	ApiToken string
	BaseUrl  string
	Username string
}

type IRCMessage struct {
	Channel string
	Msg     string
}

type IRCChannels struct {
	Key string
}

type IRCConfig struct {
	Nick     string
	Password string
	Host     string
	Ssl      bool
	Port     int
	Channels map[string]IRCChannels
}

func readConfig(configFile string) IRCConfig {
	config := IRCConfig{}

	file, _ := os.Open(configFile)
	decoder := json.NewDecoder(file)
	err := decoder.Decode(&config)

	if err != nil {
		log.Println("Cannot parse config:", err)
		os.Exit(1)
	}

	return config
}

func Bot(config IRCConfig) {

	c := irc.NewConfig(config.Nick, config.Nick, config.Nick)
	c.SSL = config.Ssl
	c.Server = config.Host
	c.SSLConfig = &tls.Config{ServerName: c.Server}

	bot := irc.Client(c)

	quit := make(chan bool)

	// on CONNECTED identify with NickServ, set mode to bot, connect to channels in config
	bot.HandleFunc(irc.CONNECTED,
		func(conn *irc.Conn, line *irc.Line) {
			conn.Mode(conn.Me().Nick, "+B")
			bot.Privmsgf("NickServ", "identify %s", config.Password)
			fmt.Println(line.Raw)
			for channel, cinfo := range config.Channels {
				if cinfo.Key != "" {
					conn.Join(channel + " " + cinfo.Key)
				} else {
					conn.Join(channel)
				}
			}
		})

	bot.HandleFunc(irc.PRIVMSG,
		func(conn *irc.Conn, line *irc.Line) {
			log.Println(line.Raw)
			switch {
			case strings.HasPrefix(line.Args[1], "!quit"):
				quit <- true

			}
		})

	bot.HandleFunc(irc.DISCONNECTED,
		func(conn *irc.Conn, line *irc.Line) { quit <- true })

	if err := bot.Connect(); err != nil {
		log.Printf("Connection error: %s\n", err.Error())
	}

	go danceActions(bot, config.Channels)

	<-quit

}

// sends IRC messages within routines
func sendMsg(bot *irc.Conn, c chan IRCMessage) {
	for item := range c {
		bot.Privmsg(item.Channel, item.Msg)
		log.Println(item.Channel, item.Msg)
	}

}

func main() {
	config := readConfig("config.json")
	Bot(config)
}
