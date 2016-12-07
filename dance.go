package main

import (
	irc "github.com/fluffle/goirc/client"
	"strings"
)

func danceActions(bot *irc.Conn, channels map[string]IRCChannels) {

	c := make(chan IRCMessage)

	go sendMsg(bot, c)

	// line.Args[0] contains the channel/sender
	// line.Args[1] contains the message
	bot.HandleFunc(irc.PRIVMSG,
		func(conn *irc.Conn, line *irc.Line) {
			switch {
			case strings.HasPrefix(line.Args[1], "!dance"):
				c <- IRCMessage{line.Args[0], ":D\\<"}
				c <- IRCMessage{line.Args[0], ":D|<"}
				c <- IRCMessage{line.Args[0], ":D/<"}
				c <- IRCMessage{line.Args[0], ":D|<"}
				c <- IRCMessage{line.Args[0], ":D\\<"}
			case strings.HasPrefix(line.Args[1], "!angrydance"):
				c <- IRCMessage{line.Args[0], ">\\D:"}
				c <- IRCMessage{line.Args[0], ">|D:"}
				c <- IRCMessage{line.Args[0], ">/D:"}
				c <- IRCMessage{line.Args[0], ">|D:"}
				c <- IRCMessage{line.Args[0], ">\\D:"}
			}
		})

}
