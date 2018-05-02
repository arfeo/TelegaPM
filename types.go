package main

type Config struct {
	Token           string        `json:"BOT_TOKEN"`
	Extension       string        `json:"STORAGE_EXT"`
	Lock            []string      `json:"LOCK_FOR_USER_IDS"`
}

type Storage struct {
	ID              string        `json:"id"`
	Title           string        `json:"title"`
	Login           string        `json:"login"`
	Pass            string        `json:"pass"`
	Url             string        `json:"url"`
}

type Current struct {
	Command         string
	Step            int
	Entry           Storage
	Element         int
	Continuous      bool
}

type Output struct {
	Response        string
	Buttons         []InlineButton
}

type InlineButton struct {
	Text            string
	URL             string
}
