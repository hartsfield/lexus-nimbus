package main

import (
	"encoding/json"
	"log"
	"os"
	"time"
)

func cycleWrite() {
	for {
		time.Sleep(5 * time.Second)
		writeJSON(started, "started.json")
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
