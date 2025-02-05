package main

import (
	"fmt"
	"os"
	"sort"
	"strings"
)

type sortable struct {
	Page   *page
	Phrase string
	Score  int
}

var unsorted []map[string]int = []map[string]int{
	wordmap,
	twowords,
	threewords,
	fourwords,
	fivewords,
	sixwords,
}

var wordmap map[string]int = map[string]int{}
var twowords map[string]int = map[string]int{}
var threewords map[string]int = map[string]int{}
var fourwords map[string]int = map[string]int{}
var fivewords map[string]int = map[string]int{}
var sixwords map[string]int = map[string]int{}

func (p *page) Nimbus() {
	if strings.Contains(strings.ToLower(p.Text), "enewspaper") ||
		strings.Contains(strings.ToLower(p.Text), "jamaica afghan") {
		return
	}
	sentences_ := strings.ReplaceAll(p.Text, "|", ".")
	sentences_ = strings.ReplaceAll(sentences_, "?", ".")
	sentences_ = strings.ReplaceAll(sentences_, "!", ".")
	sentences_ = strings.ReplaceAll(sentences_, "?", ".")
	sentences_ = strings.ReplaceAll(sentences_, ":", ".")
	sentences := strings.Split(sentences_, ".")
	for _, sen := range sentences {
		sen = strings.TrimSpace(sen)
		if len(sen) > 14 {
			if p.Sentences == nil {
				p.Sentences = make(map[string]int)
			}
			if sites.Sentences == nil {
				sites.Sentences = make(map[string]int)
			}
			if !strings.Contains(strings.ToLower(sen), "subscribe") &&
				!strings.Contains(strings.ToLower(sen), "newsletter") &&
				!strings.Contains(strings.ToLower(sen), "since 19") &&
				!strings.Contains(strings.ToLower(sen), "since 18") &&
				!strings.Contains(strings.ToLower(sen), "©") &&
				!strings.Contains(strings.ToLower(sen), "cookie") &&
				!strings.Contains(strings.ToLower(sen), "your list") &&
				!strings.Contains(strings.ToLower(sen), "your list") &&
				!strings.Contains(strings.ToLower(sen), "inbox") &&
				!strings.Contains(strings.ToLower(sen), "local event") &&
				!strings.Contains(strings.ToLower(sen), "get latest news") &&
				!strings.Contains(strings.ToLower(sen), "weather,") &&
				!strings.Contains(strings.ToLower(sen), "latest headl") &&
				!strings.Contains(strings.ToLower(sen), "deltamin") &&
				!strings.Contains(strings.ToLower(sen), "links submit") &&
				!strings.Contains(strings.ToLower(sen), "editorial") &&
				!strings.Contains(strings.ToLower(sen), "updates") &&
				!strings.Contains(strings.ToLower(sen), "skip to content") &&
				!strings.Contains(strings.ToLower(sen), "new window") &&
				!strings.Contains(strings.ToLower(sen), "new tab") &&
				!(strings.Contains(strings.ToLower(sen), "facebook") &&
					strings.Contains(strings.ToLower(sen), "twitter")) &&

				!(strings.Contains(strings.ToLower(sen), "season") &&
					strings.Contains(strings.ToLower(sen), "episode")) &&
				!(strings.Contains(strings.ToLower(sen), "sports") &&
					strings.Contains(strings.ToLower(sen), "journal")) &&
				!(strings.Contains(strings.ToLower(sen), "and has") &&
					strings.Contains(strings.ToLower(sen), "maybe") &&
					strings.Contains(strings.ToLower(sen), "said")) &&
				!strings.Contains(strings.ToLower(sen), "copyright") &&
				!strings.Contains(strings.ToLower(sen), "recent article") &&
				!strings.Contains(strings.ToLower(sen), "rights reserved") &&
				!strings.Contains(strings.ToLower(sen), "sally kestin") &&
				!strings.Contains(strings.ToLower(sen), "email") &&
				!strings.Contains(strings.ToLower(sen), "podcast") &&
				!strings.Contains(strings.ToLower(sen), "miss a beat") &&
				!strings.Contains(strings.ToLower(sen), "latest stor") &&
				!strings.Contains(strings.ToLower(sen), "roll call") &&
				!strings.Contains(strings.ToLower(sen), "breaking news") &&
				!strings.Contains(strings.ToLower(sen), "hot topic") &&
				!strings.Contains(strings.ToLower(sen), "sports betting") &&
				!strings.Contains(strings.ToLower(sen), "expand expand") &&
				!strings.Contains(strings.ToLower(sen), "this website") &&
				!strings.Contains(strings.ToLower(sen), "weightloss") &&
				!strings.Contains(strings.ToLower(sen), "weight loss") &&
				!strings.Contains(strings.ToLower(sen), "weight-loss") &&
				!strings.Contains(strings.ToLower(sen), "}}") &&
				!strings.Contains(strings.ToLower(sen), "{{") &&
				!strings.Contains(strings.ToLower(sen), "sign in") &&
				!strings.Contains(strings.ToLower(sen), "sign-in") &&
				!strings.Contains(strings.ToLower(sen), "sign-up") &&
				!strings.Contains(strings.ToLower(sen), "sign up") &&
				!strings.Contains(strings.ToLower(sen), "and more") &&
				!strings.Contains(strings.ToLower(sen), "up to $") &&
				!strings.Contains(strings.ToLower(sen), "year in review") &&
				!strings.Contains(strings.ToLower(sen), "time deposit") &&
				!strings.Contains(strings.ToLower(sen), "time payment") &&
				!strings.Contains(strings.ToLower(sen), "dedicated to") &&
				!strings.Contains(strings.ToLower(sen), "forecast") &&
				!strings.Contains(strings.ToLower(sen), "bet $") &&
				!strings.Contains(strings.ToLower(sen), "promo") &&
				!strings.Contains(strings.ToLower(sen), "reporting") &&
				!strings.Contains(strings.ToLower(sen), "privacy policy") &&
				!strings.Contains(strings.ToLower(sen), "obituaries") &&
				!strings.Contains(strings.ToLower(sen), "rights reserve") &&
				!strings.Contains(strings.ToLower(sen), "24/7 live") {
				p.Sentences[sen] = p.Sentences[sen] + 1
				sites.Sentences[sen] = sites.Sentences[sen] + 1
			}
		}
	}
	words := strings.Split(p.Text, " ")
	depth := 3
	for i := 0; i < len(words)-depth; i++ {
		formed := []string{words[i]}
		for j := 0; j < depth; j++ {
			formed = append(formed, strings.TrimSpace(strings.ToLower(formed[j]+" "+words[i+1+j])+" "))
		}
		if filter_nimbus_1(formed) {
			for _, ph := range formed {
				indx := strings.Count(strings.TrimSpace(ph), " ")
				unsorted[indx][ph] = unsorted[indx][ph] + 1
			}
		}
	}
}
func multiSort() {
	for _, un := range unsorted {
		var sortables []*sortable = []*sortable{}
		for w, c := range un {
			sortables = append(sortables, &sortable{Phrase: w, Score: c})
		}
		sort.Slice(sortables, func(i, j int) bool {
			return sortables[i].Score < sortables[j].Score
		})
		for _, s := range sortables {
			if s.Score > 8 && s.Score < sortables[len(sortables)-1].Score {
				for sen := range sites.Sentences {
					if len(sen) < 200 {
						if strings.Contains(strings.ToLower(sen), strings.ToLower(s.Phrase)) {
							// fmt.Println(fmt.Sprint(s.Score), "\t", s.Phrase)
							sites.Sentences[sen] = sites.Sentences[sen] + 1
							// fmt.Println(sen)
							// fmt.Println()
						}
					}
				}
			}
		}

	}
	var sortables []*sortable = []*sortable{}
	for w, c := range sites.Sentences {
		sortables = append(sortables, &sortable{Phrase: w, Score: c})
	}
	sort.Slice(sortables, func(i, j int) bool {
		return sortables[i].Score < sortables[j].Score
	})
	for _, s := range sortables {
		if s.Score > 4 {
			fmt.Println(s.Score, s.Phrase)
			fmt.Println()
		}
	}
	os.Exit(0)
}
