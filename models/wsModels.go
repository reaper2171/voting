package models

type Vote struct {
	User    string `json:"user"`
	Session string `json:"session"`
	Opt     string `json:"opt"`
}

type Session struct {
	Id      string
	Creator string
	Votes   map[string]int
}
