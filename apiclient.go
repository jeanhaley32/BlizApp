package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"reflect"
	"sort"
	"strings"
	"time"
)

const (
	// The URL for the Blizzard API.
	apiURL = "https://us.api.blizzard.com/hearthstone/cards/"
	// The URL for the Blizzard API token.
	tokenURL = "https://oauth.battle.net/token"
	// Define Page Limit
	pageLimit = 2
	// Define Locale
	locale = "en_US"
)

type secrets struct {
	clientID string
	secret   string
}

type criteria map[string]any

type client struct {
	secrets      *secrets
	apiKey       string
	apiKeyExpiry time.Time
	criteria     criteria
}

// Card struct for the Hearthstone API.
type Card struct {
	ID         int    `json:"id"`
	ClassID    class  `json:"classId"`
	CardTypeID int    `json:"cardTypeId"`
	CardSetID  int    `json:"cardSetId"`
	RarityID   rarity `json:"rarityId"`
	ManaCost   int    `json:"manaCost"`
	Name       string `json:"name"`
	Text       string `json:"text"`
	Image      string `json:"image"`
}

type rarity int
type class int

// Creating this enum logic allows me to easily set filter parameters,
// and also to easily print those values. I tend to use this pattern alot
// in my code, and find alot of value in it.
const (
	legendary rarity = 5
	druid     class  = 2
	warlock   class  = 9
)

func (r rarity) String() string {
	switch r {
	case legendary:
		return "legendary"
	default:
		return "unknown"
	}
}

func (c class) String() string {
	switch c {
	case druid:
		return "druid"
	case warlock:
		return "warlock"
	case 12:
		return "warlock"
	default:
		return "unknown"
	}
}

// CardsResponse struct for the Hearthstone API.
type CardsResponse struct {
	Cards []Card `json:"cards"`
}

// Obtain 10 cards from the Hearthstone API.
func (c *client) GetCard() []Card {
	var CardPages []CardsResponse
	// There is logic here to handle multiple pages of results. From what I can tell, the criteria set for this
	// challenge will only return 1 pages of results.
	// I want to write logic that will actually scale the number of pages returned based on information recieved from the API.
	// But that is outside the scope of this challenge.
	for i := 1; i <= pageLimit; i++ {
		url := apiURL + "?locale=" + locale + "&access_token=" + c.apiKey + "&page=" + fmt.Sprintf("%v", i)
		// append the search creteria to the url.
		for k, v := range c.criteria {
			// Check if value is a slice.
			if reflect.TypeOf(v).Kind() == reflect.Slice {
				// append the slice to the url, in format k=v1,v2,v3...
				url += "&" + k + "="
				for _, i := range v.([]any) {
					url += fmt.Sprintf("%v", i) + ","
				}
				// remove the trailing comma.
				url = strings.TrimSuffix(url, ",")
				continue
			} else {
				// append the search criteria to the url. format k=v
				url += "&" + k + "=" + fmt.Sprintf("%v", v)
			}
		}
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			log.Fatal(err)
		}
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			log.Fatal(err)
		}

		defer resp.Body.Close()
		cardPage := CardsResponse{}
		// Read the response body to a string
		err = json.NewDecoder(resp.Body).Decode(&cardPage)
		if err != nil {
			log.Fatal(err)
		}

		CardPages = append(CardPages, cardPage)
	}
	// Decode the JSON response
	var r CardsResponse
	for _, page := range CardPages {
		r.Cards = append(r.Cards, page.Cards...)
	}
	// return the cards.
	return r.Cards
}

// Constructing a Get request for client credentials.
func (c *client) GetAPIKey() {

	// Check if the API Key is still valid.
	if c.apiKey != "" && time.Now().Before(c.apiKeyExpiry) {
		fmt.Println("API Key is still valid.")
		return
	}

	// Check if the clientID and secret are set.
	if c.secrets.clientID == "" || c.secrets.secret == "" {
		fmt.Println("ClientID or Secret is not set.")
		return
	}

	// Define the response struct.
	type Response struct {
		Access_token string `json:"access_token"`
		Expires_in   int    `json:"expires_in"`
	}

	// Encode the clientID and clientSecret for the request.
	authdata := c.secrets.clientID + ":" + c.secrets.secret

	formData := url.Values{
		"grant_type": {"client_credentials"},
	}

	req, _ := http.NewRequest("POST", tokenURL,
		strings.NewReader(formData.Encode()))

	req.Header.Add(("Content-Type"), "application/x-www-form-urlencoded")
	req.Header.Add("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(authdata)))

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	var r Response

	err = json.NewDecoder(resp.Body).Decode(&r)
	if err != nil {
		log.Fatal(err)
	}
	c.apiKey = r.Access_token
	c.apiKeyExpiry = time.Now().Add(time.Duration(r.Expires_in) * time.Second)

}

// Requests a block of cards from the Hearthstone API, and returns 10 random cards.
func cardPicker(s secrets, c criteria) []Card {
	// Pass the clientID and secret to the client struct.
	// Passing a criteria map that is optional, but can be used to filter the
	// results of the API call. This is not optional for this project.
	client := client{
		secrets:  &s,
		criteria: c,
	}
	client.GetAPIKey()
	cards := client.GetCard()
	// Shuffle the cards.
	rand.Shuffle(len(cards), func(i, j int) {
		cards[i], cards[j] = cards[j], cards[i]
	})
	// grab the first 10 cards.
	cards = cards[:10]
	// organize cards by ID.
	sort.Slice(cards, func(i, j int) bool {
		return cards[i].ID < cards[j].ID
	})

	return cards
}
