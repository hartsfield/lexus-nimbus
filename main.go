package main

import (
	"html"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/fatih/color"
)

type database struct {
	Pages []*page
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
}

func getGetList() {
	b, err := os.ReadFile("sites.csv")
	if err != nil {
		log.Println(err)
	}
	getlist = append(getlist, strings.Split(string(b), "\n")[1:]...)
}

func getHTML(path string, client *http.Client) string {
	r, err := http.NewRequest("GET", path, nil)
	if err != nil {
		log.Println(err)
	}
	r.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36")
	r.Header.Set("Content-Type", "text/plain")

	res, err := client.Do(r)
	if err != nil {
		log.Println(err)
		return ""
	}

	b, err := io.ReadAll(res.Body)
	if err != nil {
		log.Println(err)
		return ""
	}
	return string(b)
}

var getlist []string = []string{}
var complete *database = &database{Pages: []*page{}}
var errored *database = &database{Pages: []*page{}}
var started *database = &database{Pages: []*page{}}
var props []string = []string{"src=\"", "alt=\"", "srcset=\""}
var red, yel, blu, grn *color.Color = &color.Color{}, &color.Color{}, &color.Color{}, &color.Color{}
var colorslice []*color.Color = []*color.Color{red, yel, blu, grn}
var black *color.Color = &color.Color{}

func main() {
	red.AddBgRGB(168, 88, 50).AddRGB(255, 255, 255)
	yel.AddBgRGB(235, 235, 52).AddRGB(0, 0, 0)
	blu.AddBgRGB(6, 158, 128).AddRGB(255, 255, 255)
	grn.AddBgRGB(125, 158, 6).AddRGB(255, 255, 255)
	black.AddBgRGB(0, 0, 0).AddRGB(255, 255, 255)
	getGetList()
	client := &http.Client{}
	for _, path := range getlist[:50] {
		var p *page = &page{
			Link: strings.Split(path, ",")[1],
			HTML: strings.Join(strings.Fields(getHTML(strings.Split(path, ",")[1], client)), " "),
		}
		p.Phrases = make(map[string]int)
		p.ImgsWithAlts = map[string][]string{}
		p.ContentMap = make(map[string]int)
		p.Links = make(map[string]int)
		started.Pages = append(started.Pages, p)

		if p.HTML != "" {
			p.Status = "started"
			getTitleAndContent(p)
			p.Status = "linking"
			findLinks(p)
			p.Status = "wording"
			findWords(p)
			p.Status = "imaging"
			getImageLinkAndAltText(p)
			p.Status = "complete"
			complete.Pages = append(complete.Pages, p)
			printData(p)
			writeJSON(complete, "complete.json")
		}
	}
	// printData(complete)
}

var types []string = []string{".png", ".jpg", ".jpeg", ".webp", ".webm", ".mp4"}

func findLinks(p *page) {
	s := p.HTML
	for _, t := range types {
		s = strings.ReplaceAll(s, t, t+"\n")
	}
	s = strings.ReplaceAll(s, "url(", "http(")
	s = strings.ReplaceAll(s, "data- http", "http(")
	s = strings.ReplaceAll(s, "data-http", "http(")
	s = strings.ReplaceAll(s, "href=\"", "http(")
	s = strings.ReplaceAll(s, "href='", "http(")
	s = strings.ReplaceAll(s, "src=\"", "http(")
	s = strings.ReplaceAll(s, "src='", "http(")
	s = strings.ReplaceAll(s, p.Link, "http("+p.Link)
	s = strings.ReplaceAll(s, "http(https://", "https://")
	s = strings.ReplaceAll(s, "http(s", "https://")
	s = strings.ReplaceAll(s, "http(http://", "http://")
	s = strings.ReplaceAll(s, "http(http", "http://")
	s = strings.ReplaceAll(s, "http(://", "http://")
	s = strings.ReplaceAll(s, "\\/", "/")
	s = strings.ReplaceAll(s, "http(", "http://")
	s = strings.ReplaceAll(s, "),http", "http://")
	s = strings.ReplaceAll(s, "https://", "http://")
	s = strings.ReplaceAll(s, "http://", "\nhttp://")
	s = strings.ReplaceAll(s, "///", "//")
	s_ := strings.Split(s, "\"")
	for _, ss_ := range s_ {
		if filter_1(ss_) && strings.Contains(s, "http://") && strings.Contains(s, ".") {
			if len(ss_) > 10 {
				ss__ := strings.Split(strings.Join(strings.Fields(ss_), " "), " ")
				for _, _s := range ss__ {
					if strings.Contains(_s, "http") && len(_s) > 12 {
						_s = fixback(_s)
						p.Links[_s] = p.Links[_s] + 1
					}
				}
			}
		}
	}
}
func findWords(p *page) {
	istextthere := strings.Split(p.HTML, "</")
	if len(istextthere) > 0 {
		for _, s := range istextthere {
			s = html.UnescapeString(s)
			if filter_1(s) {
				test := strings.Split(s, ">")
				if len(test[len(test)-1]) > 2 {
					final := strings.Join(strings.Fields(test[len(test)-1]), " ")
					if len(final) > 1 {
						final = fixfront(final)
					}
					if len(final) > 1 && isMostlyLetters(final) {
						p.Phrases[final] = p.Phrases[final] + 1
					}
				}
			}
		}
	}
}

func getImageLinkAndAltText(p *page) {
	imgs := strings.Split(strings.ReplaceAll(p.HTML, "<img", "\n<img"), "\n")
	if len(imgs) > 0 {
		for _, im := range imgs {
			var alt string
			var set []string
			for _, pr := range props {
				if strings.Contains(im, props[0]) && strings.Contains(im, props[1]) {
					s := strings.Split(strings.Split(im, ">")[0], pr)
					if len(s) > 1 {
						switch pr {
						case props[2]:
							set = append(set, strings.Split(strings.Split(s[1], "\"")[0], ",")...)
						case "src=\"":
							if p.Image == "" {
								iimg_ := strings.Split(strings.Split(s[1], "\"")[0], "?")[0]
								if !strings.Contains(iimg_, "logo") {
									p.Image = iimg_
								}
							}
							set = append(set, strings.Split(strings.Split(s[1], "\"")[0], "?")[0])
						case "alt=\"":
							alt = strings.Split(strings.Split(s[1], "\"")[0], "?")[0]
						}
					}
				}
			}
			for _, sst := range set {
				p.ImgsWithAlts[alt] = append(p.ImgsWithAlts[alt], strings.Split(sst, "?")[0])
			}
		}
	}
}

func getTitleAndContent(p *page) *page {
	for _, l := range strings.Split(p.HTML, "content=\"") {
		c := strings.Split(l, "\"")[0]
		if !strings.ContainsAny(c, "/><=") && strings.Count(c, " ") > 1 {
			if p.Content == "" {
				p.Content = c
			}
			p.ContentMap[c] = 1
		}
	}

	title_ := strings.Split(strings.Split(p.HTML, "</title>")[0], ">")
	p.Title = html.UnescapeString(title_[len(title_)-1])
	p.Status = "content"
	return p
}

var twowords map[string]int = map[string]int{}
var threewords map[string]int = map[string]int{}
var fourwords map[string]int = map[string]int{}

func wordup() {
	for _, p := range complete.Pages {
		for ph := range p.Phrases {
			phs := strings.Split(ph, " ")
			for i, phss := range phs {
				nxt := phss + " "
				switch len(phs[i:]) - 1 {
				case 1:
					nxt = nxt + phs[i+1]
					twowords[nxt] = twowords[nxt] + 1
				case 2:
					nxt = nxt + phs[i+1] + " " + phs[i+2]
					threewords[nxt] = threewords[nxt] + 1
				case 3:
					nxt = nxt + phs[i+1] + " " + phs[i+2] + " " + phs[i+3]
					fourwords[nxt] = fourwords[nxt] + 1
				}
			}
		}
	}
}

func printData(p *page) {
	for ph := range p.Phrases {
		colorslice[rand.Intn(len(colorslice))].Print(ph + " ")
	}
}

// func printData(db *database) {
// 	// wordup()
// 	// for w, scr := range fourwords {
// 	// 	if scr > 20 {
// 	// 		fmt.Println(w, scr)
// 	// 	}
// 	// }
// 	// for w, scr := range threewords {
// 	// 	fmt.Println(w, scr)
// 	// }
// 	// for w, scr := range twowords {
// 	// 	fmt.Println(w, scr)
// 	// }

// 	for _, p := range db.Pages {
// 		fmt.Println()
// 		fmt.Println(p.Title)
// 		fmt.Println()
// 		for ph := range p.Phrases {
// 			colorslice[rand.Intn(len(colorslice))].Print(ph + " ")
// 		}
// 		// 	// for link := range p.Links {
// 		// 	// 	fmt.Println(link)
// 		// 	// }
// 		// 	// for alt := range p.ImgsWithAlts {
// 		// 	// 	if strings.Count(alt, " ") > 5 {
// 		// 	// 		fmt.Println(alt, i)
// 		// 	// 	}
// 		// 	// }

// 	}
// }
