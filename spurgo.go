package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/cenkalti/backoff/v5"
	"github.com/kari/fmi"
	urldescribe "github.com/kari/urldescribe"
	"github.com/lrstanley/girc"
)

var Version = "development"

var serv = flag.String("server", "irc.quakenet.org", "hostname of the irc server to connect to")
var port = flag.Int("port", 6667, "server port")
var nick = flag.String("nick", "spurgo", "nickname for the bot")
var chans = flag.String("chans", "#spurgo", "channels to join")
var versionFlag = flag.Bool("version", false, "print version information and quit")

func main() {
	flag.Parse()

	if *versionFlag {
		fmt.Printf("Version: %s\n", Version)
		return
	}

	config := girc.Config{
		Server: *serv,
		Port:   *port,
		Nick:   *nick,
		User:   *nick,
		Out:    os.Stdout,
	}

	client := girc.New(config)

	client.Handlers.Add(girc.CONNECTED, func(c *girc.Client, _ girc.Event) {
		c.Cmd.Join(strings.Split(*chans, ",")...)
	})

	client.Handlers.Add(girc.PRIVMSG, func(c *girc.Client, e girc.Event) {
		if e.Last() != "!info" {
			return
		}
		c.Cmd.Replyf(e, "Hei, olen %s. Kysy minulta vaikka säätä: !sää helsinki", c.Config.Nick)
	})

	client.Handlers.Add(girc.PRIVMSG, func(c *girc.Client, e girc.Event) {
		if e.IsFromChannel() || e.Source.Name != "zyx" || e.Last() != "!quit" {
			return
		}
		log.Printf("Quit trigger activated by %s", e.Source.Name)
		c.Quit("Time to die.")
	})

	client.Handlers.Add(girc.PRIVMSG, func(c *girc.Client, e girc.Event) {
		if e.IsFromChannel() || e.Source.Name != "zyx" || !strings.HasPrefix(e.Last(), "!op ") {
			return
		}

		channel := strings.TrimPrefix(e.Last(), "!op ")

		if !girc.IsValidChannel(channel) {
			return
		}

		// FIXME: Check that bot has op on channel, and that caller is on the channel as well
		c.Cmd.Mode(channel, "+o", e.Source.Name)
	})

	client.Handlers.Add(girc.PRIVMSG, func(c *girc.Client, e girc.Event) {
		msg := strings.TrimSpace(e.Last())

		if !strings.HasPrefix(msg, "!sää ") && !strings.HasPrefix(msg, "!fmi ") {
			return
		}

		// Remove whichever prefix was used
		location := strings.TrimSpace(msg)
		location = strings.TrimPrefix(location, "!sää ")
		location = strings.TrimPrefix(location, "!fmi ")

		if location == "" {
			c.Cmd.ReplyTo(e, "Käyttö: !sää <paikkakunta>")
			return
		}

		weather, _ := fmi.Weather(location)

		c.Cmd.Reply(e, weather)
	})

	client.Handlers.Add(girc.PRIVMSG, func(c *girc.Client, e girc.Event) {
		msg := strings.TrimSpace(e.Last())
		if !strings.HasPrefix(msg, "!vertaus") {
			return
		}
		arg := strings.TrimSpace(strings.TrimPrefix(msg, "!vertaus"))
		vertaus, _ := Sample("data/vertauskuvat.txt", arg)

		c.Cmd.Reply(e, vertaus)
	})

	client.Handlers.Add(girc.PRIVMSG, func(c *girc.Client, e girc.Event) {
		re := regexp.MustCompile(`https?://[^\s]+`)
		urlMatch := re.FindString(e.Last())
		if urlMatch == "" {
			return
		}
		desc, err := urldescribe.DescribeURL(context.Background(), urlMatch)
		if err != nil {
			// Original ignores; optional reply
			return
		}
		c.Cmd.Reply(e, desc)
	})

	log.Printf("Connecting to %s:%d ...", config.Server, config.Port)

	operation := func() (string, error) {
		return "", client.Connect()
	}

	_, err := backoff.Retry(context.TODO(), operation, backoff.WithBackOff(&backoff.ExponentialBackOff{InitialInterval: 2 * time.Second,
		Multiplier:          2.0,
		RandomizationFactor: 0.5,
		MaxInterval:         5 * time.Minute}), backoff.WithMaxElapsedTime(30*time.Minute), backoff.WithNotify(func(err error, delay time.Duration) {
		log.Printf("Connect failed (%v); retrying in %v", err, delay)
	}))

	if err != nil {
		log.Printf("Permanent failure after retries: %v – exiting", err)
		return
	}

	log.Print("Bot shutting down.")
}
