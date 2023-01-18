package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"
)

var (
	JiraBaseURL        = "https://%s.atlassian.net"
	JiraRestURL        = JiraBaseURL + "/rest/api/3/"
	JiraIssueBrowseURL = JiraBaseURL + "/browse/%s"
)

func main() {
	env := GetEnvironment()
	flag.Parse()

	jiraClient := NewJira(fmt.Sprintf(JiraRestURL, env.JiraDomain), env.JiraToken, env.JiraEmail)

	session, err := discordgo.New("Bot " + env.DiscordToken)
	if err != nil {
		log.Fatalf("Invalid bot parameters: %v", err)
	}
	session.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) { log.Println("Bot is up!") })
	session.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := getCommandHandlers()[i.ApplicationCommandData().Name]; ok {
			h(s, i, jiraClient, env.JiraDomain)
		}
	})

	err = session.Open()
	if err != nil {
		log.Fatalf("Cannot open the session: %v", err)
	}
	defer session.Close()

	_, err = session.ApplicationCommandBulkOverwrite(session.State.User.ID, env.GuildID, commands)
	if err != nil {
		log.Fatalf("Cannot register commands: %v", err)
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop
	log.Println("Gracefully shutting down")
}
