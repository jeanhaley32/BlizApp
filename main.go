// Main file for the blizh.go web server.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
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
	interruptlog()
	defer func() {
		log.Default().Println("Server stopped")
		log.Default().Printf("Server ran for %s", time.Since(ServerstartTime))
	}()
	secs, err := getSecrets()
	if err != nil {
		exitSeq(fmt.Errorf("failed to obtain secrets: %v", err))
	}
	client := &client{
		criteria: params,
		secrets:  &secs,
	}
	// mux is the multiplexer that will handle the requests.
	mux := goji.NewMux()
	log.Default().Println("Starting server on localhost:8080")
	mux.HandleFunc(pat.Get("/"), func(w http.ResponseWriter, r *http.Request) {
		log.Default().Printf("Request recieved from %s", r.RemoteAddr)
		w.Write(constructSite(client))
	})
	http.ListenAndServe("localhost:8080", mux)
}

// used for debugging. Kept around for future use.
func PrettyStruct(data interface{}) string {
	val, _ := json.MarshalIndent(data, "", "    ")
	return string(val)
}

// Construct site triggers the api call, and constructs the html page that will be returned to the client.
func constructSite(c *client) []byte {
	renderStart := time.Now()
	// Site is the html page that will be returned to the client.
	site := []byte(header) // header is defined in pagetemps.go
	css := []byte(css)     // css is defined in pagetemps.go
	site = append(site, css...)
	site = append(site, []byte(`<div class="cards-container">`)...)

	// Get the cards from the API.
	cards, err := c.CardPicker()
	if err != nil {
		exitSeq(err)
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
	log.Default().Printf("Page rendered in %s", time.Since(renderStart))
	return site
}

// getSecrets returns a secrets struct containing the clientID and secret,
// based on either the flags passed to the program, or the secrets.json file.
// It prioritizes the flags.
func getSecrets() (secrets, error) {
	// If the clientID and secret are passed as flags, use those.
	if clientID != "" && secret != "" {
		log.Default().Println("Using flags")
		return secrets{ClientID: clientID, Secret: secret}, nil
	}
	// return an error if the json file does not exist.
	if _, err := os.Stat(jsonFile); os.IsNotExist(err) {
		return secrets{}, fmt.Errorf("no client id or secret flags sent, no %v file found", jsonFile)
	}

	// Otherwise, read the secrets.json file.
	var s secrets
	file, err := os.ReadFile(jsonFile)
	if err != nil {
		return s, err
	}
	err = json.Unmarshal(file, &s)
	if err != nil {
		return s, err
	}
	return s, nil
}

// interruptlog is a goroutine that will listen for SIGINT and SIGTERM signals,
// and will run exitSeq when it recieves one.
func interruptlog() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-sigs
		log.Default().Printf("Recieved %v signal", sig)
		exitSeq(nil)
	}()
}

// exitSeq is a function that will log the time the server ran for, and exit the program.
func exitSeq(e error) {
	if e != nil {
		log.Default().Printf("Error: %v", e)
	}
	log.Default().Println("Server stopped")
	log.Default().Printf("Server ran for %s", time.Since(ServerstartTime))
	os.Exit(0)
}
