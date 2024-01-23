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

// Construct the site using the secrets and criteria.
func constructSite(s secrets, c criteria) []byte {
	// Site is the html page that will be returned to the client.
	site := []byte(`<!DOCTYPE html>`)
	css := []byte(`
	<style>
	html, body{
		height: 100%;
		width: 100%;
		font-family: 'Roboto', sans-serif;
		background-color: hsl(264, 100%, 99%);
		overflow:auto
		}
	h1 {
		font-size: clamp(1.5rem, 7vw, 3rem);
		text-align: left;
	}
	h2{
		font-size: clamp(1.3rem, 5vw, 2rem);
		text-align: left;
	}
	h3 {
		font-size: clamp(1rem, 3vw, 1.5rem);
		text-align: left;
	}
	p {
		font-size: clamp(.5rem, 2vw, 1rem);
		text-align: center;
	}
	.cards-container {
		background-color: white;
		align-items: center;
		max-width: 90%;
		max-height: auto;
		display: grid;
		grid-template-rows: repeat(.8fr, 1fr);
		overflow: auto;
		margin: 5px auto;
		border-radius: 20px;
		box-shadow: 0px 0px 10px rgb(181, 186, 191);
	}
	.card {
		position: relative;
		background-color: white;
		align-items: left;
		max-width: 80%;
		max-height: auto;
		display: grid;
		grid-template-rows: repeat(.8fr, 1fr);
		overflow: auto;
		margin: 5px auto;
		border-radius: 20px;
		box-shadow: 0px 0px 10px rgb(181, 186, 191);
	}
	.card-image img {
		max-width: 80%;
		max-height: auto;
		grid-row-start: 1;
		display: grid;
		width: 100%;
		border-radius: 5px;
	}
	.card-body {
		grid-column-start: 2;
		display:table;
		margin: 2rem 2rem 2rem 2rem;
		background-color: white;
	}
	.card-body #name {
		margin-left: 3rem;
		grid-row-start: 2;
		display:table;
		margin: 2rem 2rem 2rem 2rem;
		background-color: white;
	.card-body #info {
		margin-left: 3rem;
		margin-right: 3rem;
		font-size: clamp(.20rem, 1rem + .1vw, 3rem);
		text-align: left;
	}
	</style>
	`)
	site = append(site, css...)
	site = append(site, []byte(`<div class="cards-container">`)...)
	// Source defines a card representation, and is used as a template for the raymond template engine.
	source := `<div class="card">
	<div class="card-image">
	<img src="{{Image}}" alt="{{Name}}">
	</div>
	<div class="card-body">
	<h id=name>Name: {{Name}}</h>
	<p id=info>Type: {{CardTypeID}}</p>
	<p id=info>Class: {{ClassID}}</p>
	<p id=info>Set: {{CardSetID}}</p>
	<p id=info>Rarity: {{RarityID}}</p>
	</div>
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
