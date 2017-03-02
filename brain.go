package main

import (
	irc "github.com/fluffle/goirc/client"
	cleverbot "github.com/ugjka/cleverbot-go"
	"log"
	"strings"
)

func initializeBrain(key string) (brain *cleverbot.Session) {
	brain = cleverbot.New(key)

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
			case strings.Contains(line.Args[1], "teddy"):
				// Send Cleverbot a message.
				response, err := brain.Ask(strings.Trim(line.Args[1], "teddy"))
				if err != nil {
					log.Fatal(err)
				}
				c <- IRCMessage{line.Args[0], response}
			}
		})

}
