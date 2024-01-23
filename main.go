// Main file for the blizh.go web server.
package main

import (
	"encoding/json"
	"net/http"

	"github.com/aymerick/raymond"
	"goji.io"
	"goji.io/pat"
)

var (
	clientID = "b24a67ce2b23404f9a25e8ec528e0f4b"
	secret   = "WEKTIKw0O0LTMEcPQgLHNhaGgD6Abe8s"
	params   = map[string]any{
		"sort":     "ID:asc",
		"manaCost": 7,
		"rarity":   legendary,
		"class":    []any{warlock, druid},
	}
)

func main() {
	mux := goji.NewMux()
	mux.HandleFunc(pat.Get("/"), func(w http.ResponseWriter, r *http.Request) {
		w.Write(constructSite())
	})
	http.ListenAndServe("localhost:8080", mux)

}

// PrettyStruct takes an interfact of json formated data, and returns a pretty printed string representation.
func PrettyStruct(data interface{}) (string, error) {
	val, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		return "", err
	}
	return string(val), nil
}

func constructSite() []byte {
	site := []byte(`<!DOCTYPE html>`)
	source := `<div class="card">
	<img src="{{Image}}" alt="{{Name}}">
  </div>`

	cards, err := cardPicker(secrets{
		clientID: clientID,
		secret:   secret,
	}, params)
	if err != nil {
		panic(err)
	}

	tpl, err := raymond.Parse(source)
	if err != nil {
		panic(err)
	}

	for _, card := range cards {
		result, err := tpl.Exec(card)
		if err != nil {
			panic(err)
		}
		site = append(site, []byte(result)...)
	}
	return site
}
