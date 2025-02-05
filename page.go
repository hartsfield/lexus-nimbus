package main

import (
	"fmt"
	"html"
	"io"
	"net/http"
	"strings"
	"time"
)

func (p *page) Init() {
	p.Phrases = make(map[string]int)
	p.ImgsWithAlts = map[string][]string{}
	p.ContentMap = make(map[string]int)
	p.Links = make(map[string]int)
	sites.Pages = append(sites.Pages, p)
	p.Status = stat.Init
}
func (p *page) IsTime() bool {
	return time.Since(p.LastChecked) > 12*time.Hour || p.LastChecked.IsZero()
}

func (p *page) StatErr(err error) {
	p.Status = stat.NetErr
	p.LastErr = fmt.Sprint(err)
	p.LastChecked = time.Now()
	statchange <- p
}
func (p *page) StatInit() {
	p.Status = stat.Init
	statchange <- p

}
func (p *page) GetHTML(client *http.Client) {
	r, err := http.NewRequest("GET", p.Link, nil)
	if err != nil {
		p.StatErr(err)
	}
	r.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36")
	r.Header.Set("Content-Type", "text/plain")

	res, err := client.Do(r)
	if err != nil {
		p.StatErr(err)
		return
	}

	b, err := io.ReadAll(res.Body)
	if err != nil {
		p.StatErr(err)
		return
	}
	p.HTML = strings.Join(strings.Fields(string(b)), " ")
}
func (p *page) Analyze(client *http.Client) {
	if p.IsTime() {
		p.GetHTML(client)
		switch p.Status {
		case stat.Init:
			p.GetTitleAndContent()
			p.FindLinks()
			p.FindWords()
			p.GetImageLinkAndAltText()
			p.LastChecked = time.Now()
		case stat.NetErr:
			p.StatInit()
		}
	}
	wg.Done()
}
func (p *page) FindLinks() {
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
		if filter_1(ss_) && strings.Contains(ss_, "http://") && strings.Contains(ss_, ".") {
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
func (p *page) FindWords() {
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
						p.Text = p.Text + " " + final
						p.Phrases[final] = p.Phrases[final] + 1
					}
				}
			}
		}
	}
}

func (p *page) GetImageLinkAndAltText() {
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

func (p *page) GetTitleAndContent() {
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
}
