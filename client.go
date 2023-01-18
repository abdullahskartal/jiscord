package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

type JiraInterface interface {
	CreateIssue(body *CreateIssueRequest) (*CreateIssueResponse, error)
}

type Jira struct {
	baseURL string
	token   string
	email   string
}

func NewJira(baseURL, token, email string) JiraInterface {
	return &Jira{baseURL: baseURL, token: token, email: email}
}

func (j *Jira) CreateIssue(body *CreateIssueRequest) (*CreateIssueResponse, error) {
	endpoint := j.baseURL + "issue"

	bodyByte, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	request, _ := http.NewRequest(http.MethodPost, endpoint, bytes.NewBuffer(bodyByte))
	request.Header = getHeader()
	request.SetBasicAuth(j.email, j.token)
	response, _ := client.Do(request)

	if response == nil {
		return nil, errors.New("nil response")
	}

	if response.StatusCode == http.StatusCreated {
		createIssueResponse := &CreateIssueResponse{}
		bodyBytes, _ := io.ReadAll(response.Body)
		err = json.Unmarshal(bodyBytes, &createIssueResponse)
		if err != nil {
			return nil, errors.New("response unmarshal error")
		}
		return createIssueResponse, nil
	}
	return nil, errors.New("status code different from 201")
}

func getHeader() http.Header {
	return http.Header{
		"Content-Type": {"application/json"},
		"Accept":       {"application/json"},
	}
}
