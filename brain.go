package main

import (
	"github.com/CleverbotIO/go-cleverbot.io"
	irc "github.com/fluffle/goirc/client"
	"log"
	"strings"
)

func initializeBrain(creds BrainCreds) (brain *cleverbot.Session) {
	brain, err := cleverbot.New(creds.Apiuser, creds.Apikey, creds.Apinick)

	if err != nil {
		log.Fatal(err)
	}

	return
}

func brainActions(bot *irc.Conn, brain *cleverbot.Session, channels map[string]IRCChannels) {

	c := make(chan IRCMessage)
	go sendMsg(bot, c)

	// line.Args[0] contains the channel/sender
	// line.Args[1] contains the message
	bot.HandleFunc(irc.PRIVMSG,
		func(conn *irc.Conn, line *irc.Line) {
			switch {
			case strings.HasPrefix(line.Args[1], "teddy: "):
				// Send Cleverbot a message.
				response, err := brain.Ask(strings.TrimPrefix(line.Args[1], "teddy:"))
				if err != nil {
					log.Fatal(err)
				}
				c <- IRCMessage{line.Args[0], response}
			}
		})

}
