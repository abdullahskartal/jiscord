package main

import (
	"encoding/json"
	"io"
	"os"

	"github.com/bwmarrin/discordgo"
)

type UserInfo struct {
	DiscordID string `json:"discordID"`
	JiraID    string `json:"jiraID"`
}

type Projects struct {
	ID         string       `json:"id"`
	IssueTypes []IssueTypes `json:"issue_types"`
}

type IssueTypes struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func buildAssigneeChoices() []*discordgo.ApplicationCommandOptionChoice {
	jsonFile, err := os.Open("users.json")
	if err != nil {
		panic(err)
	}
	defer jsonFile.Close()
	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		panic(err)
	}

	var users map[string]UserInfo
	json.Unmarshal(byteValue, &users)

	index := 0
	choices := make([]*discordgo.ApplicationCommandOptionChoice, len(users))
	for name, info := range users {
		choices[index] = &discordgo.ApplicationCommandOptionChoice{
			Name:  name,
			Value: info.JiraID,
		}
		index++
	}
	return choices
}

func buildIssueTypesAndProjectChoices() ([]*discordgo.ApplicationCommandOptionChoice, []*discordgo.ApplicationCommandOptionChoice) {
	jsonFile, err := os.Open("projects.json")
	if err != nil {
		panic(err)
	}
	defer jsonFile.Close()
	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		panic(err)
	}

	var projects map[string]Projects
	json.Unmarshal(byteValue, &projects)

	var issueTypeChoices []*discordgo.ApplicationCommandOptionChoice
	projectNameChoices := make([]*discordgo.ApplicationCommandOptionChoice, len(projects))

	index := 0
	for projectName, project := range projects {
		projectNameChoices[index] = &discordgo.ApplicationCommandOptionChoice{
			Name:  projectName,
			Value: project.ID,
		}

		for _, issueType := range project.IssueTypes {
			issueTypeChoices = append(issueTypeChoices, &discordgo.ApplicationCommandOptionChoice{
				Name:  issueType.Name,
				Value: issueType.ID,
			})
		}
		index++
	}
	return issueTypeChoices, projectNameChoices
}
