package main

/**
 *
 *	Global types
 *
 */

type Config struct {
	Token         string		`json:"BOT_TOKEN"`
	Extension     string		`json:"STORAGE_EXT"`
}

type Storage struct {
	ID            string		`json:"id"`
	Title         string		`json:"title"`
	Login         string		`json:"login"`
	Pass          string		`json:"pass"`
	Url           string		`json:"url"`
}

type Current struct {
	Command       string
	Step          int
	Entry         Storage
	Element       int
	Continuous    bool
}

type Output struct {
	Response      string
	Buttons       []InlineButton
}

type InlineButton struct {
	Text          string
	URL           string
}

/**
 *
 *	Global variables
 *
 */

var (
	config        Config
	current       map[int]Current
	masters       map[int]string
)
