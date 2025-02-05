package main

import (
	"net/http"
	"sync"
	"time"

	"github.com/fatih/color"
)

type database struct {
	Pages     []*page
	Sentences map[string]int
}
type page struct {
	Link         string    `json:"Link"`
	ID           string    `json:"ID"`
	Submitted    time.Time `json:"submitted"`
	LastChecked  time.Time
	Status       string
	Title        string `json:"title"`
	Content      string `json:"content"`
	HTML         string
	Text         string `json:"text"`
	Image        string
	Links        map[string]int `json:"links"`
	Images       []string       `json:"images"`
	LastErr      string
	Phrases      map[string]int
	ContentMap   map[string]int
	ImgsWithAlts map[string][]string
	Sentences    map[string]int

	Twitter    string
	Name       string
	Location   string
	TimeZone   string
	Country    string
	Language   string
	Bundle     string
	Wait       string
	HasAdGuard string
}

var stat *struct{ NetErr, Init, Started, Lexed string } = &struct{ NetErr, Init, Started, Lexed string }{
	NetErr:  "error",
	Init:    "initialized",
	Started: "started",
	Lexed:   "lexed",
}

var checked int
var totalsent int
var stopsmap map[string]bool = make(map[string]bool)
var statchange chan *page = make(chan *page)
var analyzing chan *page = make(chan *page, 100)
var jsonfile string = "skink.json"
var types []string = []string{".png", ".jpg", ".jpeg", ".webp", ".webm", ".mp4"}
var sites *database = &database{Pages: []*page{}}
var props []string = []string{"src=\"", "alt=\"", "srcset=\""}
var red, yel, blu, grn *color.Color = &color.Color{}, &color.Color{}, &color.Color{}, &color.Color{}
var colorslice []*color.Color = []*color.Color{red, yel, blu, grn}
var black *color.Color = &color.Color{}
var wg sync.WaitGroup = sync.WaitGroup{}
var client *http.Client = &http.Client{
	Timeout: 4 * time.Second,
}
