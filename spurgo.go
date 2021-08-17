package main

import (
	"flag"
	"fmt"
	"regexp"
	"strings"

	"github.com/kari/fmi"
	urldescribe "github.com/kari/urldescribe"
	hbot "github.com/whyrusleeping/hellabot"
	log "gopkg.in/inconshreveable/log15.v2"
)

var serv = flag.String("server", "irc.quakenet.org:6667", "hostname and port for irc server to connect to")
var nick = flag.String("nick", "spurgo", "nickname for the bot")
var chans = flag.String("chans", "#spurgo", "channels to join")

func main() {
	flag.Parse()

	hijackSession := func(bot *hbot.Bot) {
		bot.HijackSession = true
	}
	channels := func(bot *hbot.Bot) {
		bot.Channels = strings.Split(*chans, ",")

	}
	irc, err := hbot.NewBot(*serv, *nick, hijackSession, channels)
	if err != nil {
		panic(err)
	}

	irc.AddTrigger(SayInfoMessage)
	irc.AddTrigger(QuitTrigger)
	irc.AddTrigger(OpTrigger)
	irc.AddTrigger(WeatherTrigger)
	irc.AddTrigger(WeatherTrigger2)
	irc.AddTrigger(URLTrigger)
	irc.AddTrigger(SimileTrigger)
	// WrongBotTrigger needs to be last
	// irc.AddTrigger(WrongBotTrigger)
	irc.Logger.SetHandler(log.StdoutHandler)

	// Start up bot (this blocks until we disconnect)
	irc.Run()
	fmt.Println("Bot shutting down.")
}

// SayInfoMessage replies Hello when you say !info
var SayInfoMessage = hbot.Trigger{
	Condition: func(bot *hbot.Bot, m *hbot.Message) bool {
		return m.Command == "PRIVMSG" && m.Content == "!info"
	},
	Action: func(irc *hbot.Bot, m *hbot.Message) bool {
		irc.Reply(m, fmt.Sprintf("Hei, olen %s. Kysy minulta vaikka säätä: !sää helsinki", irc.Nick))
		return true
	},
}

// QuitTrigger makes the bot shut down
var QuitTrigger = hbot.Trigger{
	Condition: func(bot *hbot.Bot, m *hbot.Message) bool {
		return m.Command == "PRIVMSG" && m.To == bot.Nick && m.Content == "!quit" && m.From == "zyx"
	},
	Action: func(irc *hbot.Bot, m *hbot.Message) bool {
		irc.Info("Quit trigger activated")
		irc.Send("QUIT :Time to die.")

		return true
	},
}

// OpTrigger makes the bot op the owner
var OpTrigger = hbot.Trigger{
	Condition: func(bot *hbot.Bot, m *hbot.Message) bool {
		return m.Command == "PRIVMSG" && m.To == bot.Nick && strings.HasPrefix(m.Content, "!op ") && m.From == "zyx"
	},
	Action: func(irc *hbot.Bot, m *hbot.Message) bool {
		irc.ChMode("zyx", strings.TrimPrefix(m.Content, "!op "), "+o")

		return true
	},
}

// WeatherTrigger check weather
var WeatherTrigger = hbot.Trigger{
	Condition: func(bot *hbot.Bot, m *hbot.Message) bool {
		return m.Command == "PRIVMSG" && strings.HasPrefix(m.Content, "!sää ")
	},
	Action: func(irc *hbot.Bot, m *hbot.Message) bool {
		irc.Reply(m, fmi.Weather(strings.TrimPrefix(m.Content, "!sää ")))
		return true
	},
}

// WeatherTrigger2 check weather
var WeatherTrigger2 = hbot.Trigger{
	Condition: func(bot *hbot.Bot, m *hbot.Message) bool {
		return m.Command == "PRIVMSG" && strings.HasPrefix(m.Content, "!fmi ")
	},
	Action: func(irc *hbot.Bot, m *hbot.Message) bool {
		irc.Reply(m, fmi.Weather(strings.TrimPrefix(m.Content, "!fmi ")))
		return true
	},
}

// SimileTrigger fetches a simile
var SimileTrigger = hbot.Trigger{
	Condition: func(bot *hbot.Bot, m *hbot.Message) bool {
		return m.Command == "PRIVMSG" && strings.HasPrefix(m.Content, "!vertaus")
	},
	Action: func(irc *hbot.Bot, m *hbot.Message) bool {
		irc.Reply(m, sample("", strings.TrimSpace(strings.TrimPrefix(m.Content, "!vertaus"))))
		return true
	},
}

// WrongBotTrigger redirects to correct bot
var WrongBotTrigger = hbot.Trigger{
	Condition: func(bot *hbot.Bot, m *hbot.Message) bool {
		re := regexp.MustCompile(`^![^!\?]+$`)
		return m.Command == "PRIVMSG" && re.MatchString(m.Content)
	},
	Action: func(irc *hbot.Bot, m *hbot.Message) bool {
		irc.Reply(m, "Tarkoititko ."+strings.TrimPrefix(m.Content, "!")+"?")
		return true
	},
}

// URLTrigger attempts to describe the link
var URLTrigger = hbot.Trigger{
	Condition: func(bot *hbot.Bot, m *hbot.Message) bool {
		re := regexp.MustCompile(`https?:\/\/[^\ ]+`)
		return m.Command == "PRIVMSG" && re.MatchString(m.Content)
	},
	Action: func(irc *hbot.Bot, m *hbot.Message) bool {
		re := regexp.MustCompile(`https?:\/\/[^\ ]+`)
		irc.Reply(m, urldescribe.DescribeURL(re.FindString(m.Content)))
		return true
	},
}
