package main

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

var (
	commands = []*discordgo.ApplicationCommand{
		{
			Name:        "create-issue",
			Description: "Create jira issue",
			Type:        discordgo.ChatApplicationCommand,
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:         "type",
					Description:  "Type of issue",
					Type:         discordgo.ApplicationCommandOptionString,
					Required:     true,
					Autocomplete: true,
				},
				{
					Name:         "project",
					Description:  "Project id",
					Type:         discordgo.ApplicationCommandOptionString,
					Required:     true,
					Autocomplete: true,
				},
				{
					Name:         "assignee",
					Description:  "Assignee of issue",
					Type:         discordgo.ApplicationCommandOptionString,
					Required:     true,
					Autocomplete: true,
				},
				{
					Name:         "title",
					Description:  "Title of issue",
					Type:         discordgo.ApplicationCommandOptionString,
					Required:     true,
					Autocomplete: false,
				},
			},
		},
	}
)

func getCommandHandlers() map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate, jiraClient JiraInterface, jiraDomain string) {
	return map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate, jiraInterface JiraInterface, jiraDomain string){
		"create-issue": createIssueHandler,
	}
}

func createIssueHandler(
	s *discordgo.Session,
	i *discordgo.InteractionCreate,
	jiraClient JiraInterface,
	jiraDomain string) {
	data := i.ApplicationCommandData()
	switch i.Type {
	case discordgo.InteractionApplicationCommand:
		handleSubmitForCreateIssue(data, jiraClient, i, jiraDomain, s)
	case discordgo.InteractionApplicationCommandAutocomplete:
		handleAutoCompleteForCreateIssue(data, s, i)
	}
}

func handleSubmitForCreateIssue(
	data discordgo.ApplicationCommandInteractionData,
	jiraClient JiraInterface,
	i *discordgo.InteractionCreate,
	jiraDomain string,
	s *discordgo.Session) {

	typeOfIssue := data.Options[0].StringValue()
	projectOfIssue := data.Options[1].StringValue()
	assigneeOfIssue := data.Options[2].StringValue()
	titleOfIssue := data.Options[3].StringValue()
	createIssueRequest := &CreateIssueRequest{
		Fields: Fields{
			Summary:     titleOfIssue,
			Issuetype:   Issuetype{ID: typeOfIssue},
			Project:     Project{ID: projectOfIssue},
			Assignee:    Assignee{ID: assigneeOfIssue},
			CustomField: 1,
		},
	}
	resp, err := jiraClient.CreateIssue(createIssueRequest)
	if err != nil {
		panic(err)
	}
	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf(JiraIssueBrowseURL, jiraDomain, resp.Key),
		},
	})
	if err != nil {
		panic(err)
	}
}

func handleAutoCompleteForCreateIssue(
	data discordgo.ApplicationCommandInteractionData,
	s *discordgo.Session,
	i *discordgo.InteractionCreate) {

	var choices []*discordgo.ApplicationCommandOptionChoice
	issueTypeChoice, projectChoice := buildIssueTypesAndProjectChoices()
	switch {
	case data.Options[0].Focused:
		choices = issueTypeChoice
	case data.Options[1].Focused:
		choices = projectChoice
	case data.Options[2].Focused:
		choices = buildAssigneeChoices()
	}

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionApplicationCommandAutocompleteResult,
		Data: &discordgo.InteractionResponseData{
			Choices: choices,
		},
	})
	if err != nil {
		panic(err)
	}
}
