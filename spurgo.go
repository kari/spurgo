// This is an example program showing the usage of hellabot
// kts myös https://github.com/go-chat-bot/bot

package main

import (
	"flag"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/kari/fmi"
	urldescribe "github.com/kari/urldescribe"
	hbot "github.com/whyrusleeping/hellabot"
	log "gopkg.in/inconshreveable/log15.v2"
)

// var serv = flag.String("server", "irc.quakenet.org:6667", "hostname and port for irc server to connect to")

var serv = flag.String("server", "chat.freenode.net:6667", "hostname and port for irc server to connect to")

var nick = flag.String("nick", "spurgo", "nickname for the bot")

func main() {
	flag.Parse()

	hijackSession := func(bot *hbot.Bot) {
		bot.HijackSession = true
	}
	channels := func(bot *hbot.Bot) {
		// bot.Channels = []string{"#spurgo"}
		bot.Channels = []string{"#reddit-suomi"}

	}
	irc, err := hbot.NewBot(*serv, *nick, hijackSession, channels)
	if err != nil {
		panic(err)
	}

	irc.AddTrigger(SayInfoMessage)
	// irc.AddTrigger(LongTrigger)
	irc.AddTrigger(QuitTrigger)
	irc.AddTrigger(WeatherTrigger)
	irc.AddTrigger(WeatherTrigger2)
	irc.AddTrigger(URLTrigger)
	irc.Logger.SetHandler(log.StdoutHandler)
	// logHandler := log.LvlFilterHandler(log.LvlInfo, log.StdoutHandler)
	// or
	// irc.Logger.SetHandler(logHandler)
	// or
	// irc.Logger.SetHandler(log.StreamHandler(os.Stdout, log.JsonFormat()))

	// Start up bot (this blocks until we disconnect)
	irc.Run()
	fmt.Println("Bot shutting down.")
}

// SayInfoMessage replies Hello when you say !info
var SayInfoMessage = hbot.Trigger{
	func(bot *hbot.Bot, m *hbot.Message) bool {
		return m.Command == "PRIVMSG" && m.Content == "!info"
	},
	func(irc *hbot.Bot, m *hbot.Message) bool {
		irc.Reply(m, fmt.Sprintf("Hello, I am %s", irc.Nick))
		return false
	},
}

// LongTrigger sends two messages with 5 second delay
var LongTrigger = hbot.Trigger{
	func(bot *hbot.Bot, m *hbot.Message) bool {
		return m.Command == "PRIVMSG" && m.Content == "!long"
	},
	func(irc *hbot.Bot, m *hbot.Message) bool {
		irc.Reply(m, "This is the first message")
		time.Sleep(5 * time.Second)
		irc.Reply(m, "This is the second message")

		return false
	},
}

// QuitTrigger makes the bot shut down
var QuitTrigger = hbot.Trigger{
	func(bot *hbot.Bot, m *hbot.Message) bool {
		return m.Command == "PRIVMSG" && m.To == bot.Nick && m.Content == "!quit" && m.From == "zyx"
	},
	func(irc *hbot.Bot, m *hbot.Message) bool {
		irc.Info("Quit trigger activated")
		irc.Send("QUIT :Time to die.")

		return true
	},
}

// WeatherTrigger check weather
var WeatherTrigger = hbot.Trigger{
	func(bot *hbot.Bot, m *hbot.Message) bool {
		return m.Command == "PRIVMSG" && strings.HasPrefix(m.Content, "!sää ")
	},
	func(irc *hbot.Bot, m *hbot.Message) bool {
		irc.Reply(m, fmi.Weather(strings.TrimPrefix(m.Content, "!sää ")))
		return false
	},
}

// WeatherTrigger2 check weather
var WeatherTrigger2 = hbot.Trigger{
	func(bot *hbot.Bot, m *hbot.Message) bool {
		return m.Command == "PRIVMSG" && strings.HasPrefix(m.Content, "!fmi ")
	},
	func(irc *hbot.Bot, m *hbot.Message) bool {
		irc.Reply(m, fmi.Weather(strings.TrimPrefix(m.Content, "!fmi ")))
		return false
	},
}

// URLTrigger attempts to describe the link
var URLTrigger = hbot.Trigger{
	func(bot *hbot.Bot, m *hbot.Message) bool {
		re := regexp.MustCompile("https?:\\/\\/[^\\ ]+")
		return m.Command == "PRIVMSG" && re.MatchString(m.Content)
	},
	func(irc *hbot.Bot, m *hbot.Message) bool {
		irc.Reply(m, urldescribe.DescribeURL(m.Content))
		return false
	},
}
