package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	irc "github.com/fluffle/goirc/client"
	"github.com/yhat/scrape"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
	"net/http"
	"os"
	"runtime"
	"strings"
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
	Password string
	Networks map[string]IRCNetworks
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

func getTitle(url string) string {

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("error:", err)
		return "error"
	}

	root, err := html.Parse(resp.Body)
	if err != nil {
		fmt.Println("error:", err)
		return "error"
	}

	title, ok := scrape.Find(root, scrape.ByTag(atom.Title))

	if ok {
		return scrape.Text(title)
	}

	return "unknown"
}

func dance() [3]string {
	var a [3]string
	a[0] = ":D\\<"
	a[1] = ":D|<"
	a[2] = ":D/<"
	return a
}

func angrydance() [3]string {
	var a [3]string
	a[0] = ">/D:"
	a[1] = ">|D:"
	a[2] = ">\\D:"
	return a
}

func teddyBot(nick string, password string, config IRCNetworks) {

	cfg := irc.NewConfig(nick)
	cfg.SSL = config.Ssl
	cfg.SSLConfig = &tls.Config{ServerName: config.Host, InsecureSkipVerify: true}
	cfg.Server = config.Host
	cfg.NewNick = func(n string) string { return n + "^" }

	bot := irc.Client(cfg)
	bot.EnableStateTracking()

	bot.HandleFunc(irc.CONNECTED,
		func(conn *irc.Conn, line *irc.Line) {
			conn.Mode(conn.Me().Nick, "+B")
			bot.Privmsg("NickServ", fmt.Sprintf("identify %s", password))
			for key, channel := range config.Channels {
				fmt.Printf("Connecting to channel #%s\n", key)
				conn.Join(channel.Name + " " + channel.Key)
			}
		})

	bot.HandleFunc(irc.PRIVMSG,
		func(conn *irc.Conn, line *irc.Line) {
			if strings.HasPrefix(line.Text(), "http") {
				bot.Privmsg(line.Args[0], getTitle(line.Text()))
			} else if strings.HasPrefix(line.Text(), "!dance") {
				for _, moves := range dance() {
					bot.Privmsg(line.Args[0], moves)
				}
			} else if strings.HasPrefix(line.Text(), "!angrydance") {
				for _, moves := range angrydance() {
					bot.Privmsg(line.Args[0], moves)
				}
			}

		})

	quit := make(chan bool)

	bot.HandleFunc(irc.DISCONNECTED,
		func(conn *irc.Conn, line *irc.Line) { quit <- true })

	if err := bot.Connect(); err != nil {
		fmt.Printf("Connection error: %s\n", err.Error())
	}

	// go func(line *irc.Line) {
	// 	fmt.Println(line)
	// }

	<-quit

}

func main() {
	runtime.GOMAXPROCS(2)
	config := readConfig("config.json")

	for _, network := range config.Networks {
		teddyBot(config.Nick, config.Password, network)
	}
}
