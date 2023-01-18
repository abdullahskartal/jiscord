package main

type CreateIssueRequest struct {
	Fields Fields `json:"fields"`
}

type Fields struct {
	Summary     string    `json:"summary"`
	Issuetype   Issuetype `json:"issuetype"`
	Project     Project   `json:"project"`
	Assignee    Assignee  `json:"assignee"`
	CustomField uint64    `json:"customfield_10032"`
}

type Issuetype struct {
	ID string `json:"id"`
}

type Project struct {
	ID string `json:"id"`
}

type Assignee struct {
	ID string `json:"id"`
}

type CreateIssueResponse struct {
	ID         string     `json:"id"`
	Key        string     `json:"key"`
	Self       string     `json:"self"`
	Transition Transition `json:"transition"`
}

type Errors struct {
}

type Errorcollection struct {
	Errormessages []interface{} `json:"errorMessages"`
	Errors        Errors        `json:"errors"`
}

type Transition struct {
	Status          int             `json:"status"`
	Errorcollection Errorcollection `json:"errorCollection"`
}
