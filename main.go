package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	red.AddBgRGB(168, 88, 50).AddRGB(255, 255, 255)
	yel.AddBgRGB(235, 235, 52).AddRGB(0, 0, 0)
	blu.AddBgRGB(6, 158, 128).AddRGB(255, 255, 255)
	grn.AddBgRGB(125, 158, 6).AddRGB(255, 255, 255)
	black.AddBgRGB(0, 0, 0).AddRGB(255, 255, 255)
	csvToJSON()
	b, err := os.ReadFile("stops.txt")
	if err != nil {
		fmt.Println(err)
	}
	for _, w := range strings.Split(string(b), "\n") {
		stopsmap[w] = true
	}
}
func cycle() {
	for {
		p := <-analyzing
		wg.Add(1)
		go p.Analyze(client)
	}
}
func main() {
	for _, p := range sites.Pages {
		if len(p.Text) > 0 {
			p.Nimbus()
		}
	}
	multiSort()
	defer close(analyzing)
	go cycle()
	go writeOnStatusChange()
	for {
		for _, p := range sites.Pages {
			analyzing <- p
		}
		wg.Wait()
		time.Sleep(1 * time.Second)
	}
}
