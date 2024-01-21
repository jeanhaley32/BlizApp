// Main file for the blizh.go web server.
package main

import "fmt"

var (
	clientID = "b24a67ce2b23404f9a25e8ec528e0f4b"
	secret   = "WEKTIKw0O0LTMEcPQgLHNhaGgD6Abe8s"
)

func main() {
	// Set the Criteria for this challenge.
	// Cards will be sorted by ID in ascending order.
	// Cards will have a manaCost of 7.
	// Cards will be legendary.
	// Cards will be either warlock or druid.
	criteria := map[string]any{
		"sort":     "ID:asc",
		"manaCost": 7,
		"rarity":   legendary,
		"class":    []any{warlock, druid},
	}

	cards := cardPicker(secrets{
		clientID: clientID,
		secret:   secret,
	}, criteria)

	for i, card := range cards {
		fmt.Printf("NUM: %d ID: %d | Name: %s | Class: %s | Rarity: %s | Mana Cost: %d\n",
			i+1, card.ID, card.Name, card.ClassID, card.RarityID, card.ManaCost)
	}
	// mux := goji.NewMux()
	// mux.HandleFunc(pat.Get("/"), blizh)
	// http.ListenAndServe("localhost:8000", mux)

}
