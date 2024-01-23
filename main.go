// Main file for the blizh.go web server.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/aymerick/raymond"
	"goji.io"
	"goji.io/pat"
)

type secrets struct {
	ClientID string `json:"clientid"`
	Secret   string `json:"secret"`
}

var (
	clientID, secret, jsonFile string
	ServerstartTime            time.Time
	// params is a map of the query parameters that will be sent to the api.
	params = map[string]any{
		"sort":     "ID:asc",
		"manaCost": 7,
		"rarity":   legendary,
		"class":    []any{warlock, druid},
	}
)

// initializaztion function, called before main.
func init() {
	ServerstartTime = time.Now()
	flag.StringVar(&jsonFile, "json", "secrets.json", "json file containing the clientID and secret")
	flag.StringVar(&clientID, "clientid", "", "clientID for the blizzard api")
	flag.StringVar(&secret, "secret", "", "secret for the blizzard api")
	flag.Parse()
}

func main() {
	secs, err := getSecrets()
	if err != nil {
		panic(fmt.Errorf("failed to obtain secretes: %v", err))
	}
	// Using the goji muxer to handle requests.
	// I chose goji because it's a simple and fast muxer.
	mux := goji.NewMux()
	log.Default().Println("Starting server on localhost:8080")
	mux.HandleFunc(pat.Get("/"), func(w http.ResponseWriter, r *http.Request) {
		log.Default().Printf("Request recieved from %s", r.RemoteAddr)
		w.Write(constructSite(secs, params))
	})
	http.ListenAndServe("localhost:8080", mux)
	log.Default().Println("Server stopped")
	log.Default().Printf("Server ran for %s", time.Since(ServerstartTime))
}

func PrettyStruct(data interface{}) string {
	val, _ := json.MarshalIndent(data, "", "    ")
	return string(val)
}

// construcSite is a helper function that constructs the html page.
// It uses raymond to render the html template.
// I chose raymond because it's a simple and fast template engine. Same reason I chose goji.
// This is the entry point for the client server that speaks to the blizzard api.
func constructSite(s secrets, c criteria) []byte {
	// Site is the html page that will be returned to the client.
	site := []byte(`<!DOCTYPE html>`)

	// Source defines a card representation, and is used as a template for the raymond template engine.
	source := `<div class="card">
	<img src="{{Image}}" alt="{{Name}}">
  </div>`

	// Get the cards from the API.
	cards, err := cardPicker(s, c)
	if err != nil {
		panic(err)
	}

	// Parse the template.
	tpl, err := raymond.Parse(source)
	if err != nil {
		panic(err)
	}

	// Render the template, and append the result to the site.
	for _, card := range cards {
		result, err := tpl.Exec(card)
		if err != nil {
			panic(err)
		}
		site = append(site, []byte(result)...)
	}
	// return the site.
	log.Default().Printf("Rendering Page")
	return site
}

// getSecrets obtains the clientID and secret from the secrets.json file.
func getSecrets() (secrets, error) {
	// If the clientID and secret are passed as flags, use those.
	if clientID != "" && secret != "" {
		return secrets{ClientID: clientID, Secret: secret}, nil
	}
	// Otherwise, read the secrets.json file.
	var s secrets
	file, err := os.ReadFile("secrets.json")
	if err != nil {
		return s, err
	}
	err = json.Unmarshal(file, &s)
	if err != nil {
		return s, err
	}
	return s, nil
}
