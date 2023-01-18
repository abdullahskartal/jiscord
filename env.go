package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type EnvironmentVariable struct {
	JiraDomain   string
	JiraEmail    string
	JiraToken    string
	DiscordToken string
	GuildID      string
}

// ParseEnv Get environment value from os
// If an environment required and not set raises a panic
func ParseEnv(key string, required bool, dft string) string {
	value := os.Getenv(key)
	if value == "" && required {
		panic(fmt.Sprintf("Environment variable not found: %v", key))
	} else if value == "" {
		return dft
	}
	return value
}

func GetEnvironment() *EnvironmentVariable {
	_ = godotenv.Load()
	return &EnvironmentVariable{
		JiraDomain:   ParseEnv("JIRA_DOMAIN", true, ""),
		JiraEmail:    ParseEnv("JIRA_EMAIL", true, ""),
		JiraToken:    ParseEnv("JIRA_TOKEN", true, ""),
		DiscordToken: ParseEnv("DISCORD_TOKEN", true, ""),
		GuildID:      ParseEnv("GUILD_ID", true, ""),
	}
}
