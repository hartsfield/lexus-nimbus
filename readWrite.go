package main

import (
	"encoding/json"
	"log"
	"os"
	"strings"
)

func writeOnStatusChange() {
	for {
		for range statchange {
			writeJSON(sites, jsonfile)
		}
	}
}
func writeJSON(d *database, fn string) {
	b, err := json.MarshalIndent(d, "", "    ")
	if err != nil {
		log.Println(err)
	}
	appendFile(string(b), fn)
}
func appendFile(l, fn string) {
	f, err := os.OpenFile(fn, os.O_TRUNC|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		log.Println(err)
	}

	defer f.Close()

	if _, err = f.WriteString(l); err != nil {
		log.Println(err)
	}

}

func csvToJSON() {
	b, err := os.ReadFile(jsonfile)
	if err != nil {
		log.Println(err)
	}

	err = json.Unmarshal(b, sites)
	if err != nil {
		log.Println(err)
		b, err = os.ReadFile("sites.csv")
		if err != nil {
			log.Println(err)
		}
		csv := strings.Split(string(b), "\n")[1:]
		for _, line := range csv {
			vals := strings.Split(line, ",")
			if len(vals) >= 10 {
				var p *page = &page{}
				p.Twitter = vals[0]
				p.Link = vals[1]
				p.Name = vals[2]
				p.Location = vals[3]
				p.TimeZone = vals[4]
				p.Country = vals[5]
				p.Language = vals[6]
				p.Bundle = vals[7]
				p.Wait = vals[8]
				p.HasAdGuard = vals[9]
				p.Init()
			}
		}
		writeJSON(sites, jsonfile)
	}
}

// var statmap map[string]int = make(map[string]int)

// func dbStats() {
// 	for _, p := range sites.Pages {
// 		statmap[p.Status] = 0
// 	}
// 	for _, p := range sites.Pages {
// 		statmap[p.Status] = statmap[p.Status] + 1
// 	}
// 	// clearTerm(20)
// 	fmt.Println(stat.NetErr, "\t\t", statmap[stat.NetErr])
// 	fmt.Println(stat.Init, "\t", statmap[stat.Init])
// 	fmt.Println(stat.Started, "\t", statmap[stat.Started])
// 	fmt.Println(stat.Lexed, "\t\t", statmap[stat.Lexed])
// 	fmt.Println(totalsent, len(analyzing), fmt.Sprint(checked)+"/"+fmt.Sprint(len(sites.Pages)))
// 	fmt.Println()

// }
//
//	func clearTerm(lines int) {
//		for i := 0; i <= lines; i++ {
//			fmt.Println()
//		}
//	}
