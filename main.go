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

	"goji.io"
	"goji.io/pat"
)

var (
	clientID, secret, jsonFile string
	ServerstartTime            time.Time
	// params is a map of the query parameters that will be sent to the api.
	params = map[string]any{
		"sort":     "ID:asc",
		"manaCost": 7,
		"rarity":   LEGENDARY,
		"class":    []any{WARLOCK, DRUID},
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

// Construct the site using the secrets and criteria.
func constructSite(s secrets, c criteria) []byte {
	// Site is the html page that will be returned to the client.
	site := []byte(header) // header is defined in pagetemps.go
	css := []byte(css)     // css is defined in pagetemps.go
	site = append(site, css...)
	site = append(site, []byte(`<div class="cards-container">`)...)

	// Get the cards from the API.
	cards, err := cardPicker(s, c)
	if err != nil {
		panic(err)
	}

	// Render the template, and append the result to the site.
	for _, card := range cards {
		result := fmt.Sprintf(source,
			card.Image, card.Name, card.Name,
			card.ID, card.CardTypeID, card.ClassID,
			card.CardSetID, card.RarityID)
		site = append(site, []byte(result)...)
	}
	site = append(site, []byte(`</div>`)...)
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
