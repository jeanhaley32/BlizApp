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
	apiURL   = "https://us.api.blizzard.com/hearthstone/cards/"
	tokenURL = "https://oauth.battle.net/token"
	// limit set for pagination.
	locale = "en_US"
)

var pageLimit = 10

// criteria type is used to define a map of query parameters for the API call.
type criteria map[string]any

// client is the struct that holds the secrets, apikey, and criteria for the API call.
type client struct {
	secrets      *secrets
	apiKey       string
	apiKeyExpiry time.Time
	criteria     criteria
}

// Card holds just the data we need for a single Hearthstone card.
type Card struct {
	ID         int    `json:"id"`
	ClassID    Class  `json:"classId"`
	CardTypeID Type   `json:"cardTypeId"`
	CardSetID  Set    `json:"cardSetId"`
	RarityID   Rarity `json:"rarityId"`
	ManaCost   int    `json:"manaCost"`
	Name       string `json:"name"`
	Text       string `json:"text"`
	Image      string `json:"image"`
}

// secrets is the struct that holds the clientID and secret for the API call.
type secrets struct {
	ClientID string `json:"clientid"`
	Secret   string `json:"secret"`
}

// Card Response is the upper level struct that holds the cards and page count.
// page count is used for pagination.
type CardsResponse struct {
	Cards     []Card `json:"cards"`
	PageCount int    `json:"pageCount"`
}

// GetCard is the function that makes the API call and returns a slice of cards.
func (c *client) getCard() ([]Card, error) {
	var CardPages []CardsResponse
	pages := 1
	for i := 1; i <= pages; i++ {
		url := apiURL + "?locale=" + locale + "&access_token=" + c.apiKey + "&page=" + fmt.Sprintf("%v", i)
		// concatenate the criteria to the url.
		url += concatCriteria(&c.criteria)
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return nil, err
		}
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()
		cardPage := CardsResponse{}
		err = json.NewDecoder(resp.Body).Decode(&cardPage)
		if err != nil {
			return nil, err
		}
		CardPages = append(CardPages, cardPage)
		if i == pageLimit {
			break
		}
		pages = cardPage.PageCount
	}
	log.Default().Printf("[page limit %v] - Received %v pages of %v\n",
		pageLimit, len(CardPages), CardPages[0].PageCount)
	if len(CardPages) == 0 {
		return nil, fmt.Errorf("no cards returned from the API")
	}
	var r CardsResponse
	// Flatten the pages into a single slice.
	for i, page := range CardPages {
		r.Cards = append(r.Cards, page.Cards...)
		if i == len(CardPages)-1 {
			log.Default().Printf("Pulled %v cards\n", len(r.Cards))
		}
	}
	return r.Cards, nil
}

func (c *client) getAPIKey() error {
	if c.apiKey != "" && time.Now().Before(c.apiKeyExpiry) {
		log.Default().Printf("API Key is still valid")
		return nil
	}
	log.Default().Println("Getting API Key")
	if c.secrets.ClientID == "" || c.secrets.Secret == "" {
		return fmt.Errorf("ClientID or Secret is not set")
	}

	type Response struct {
		Access_token string `json:"access_token"`
		Expires_in   int    `json:"expires_in"`
	}

	authdata := c.secrets.ClientID + ":" + c.secrets.Secret

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
		return err
	}

	defer resp.Body.Close()

	var r Response

	err = json.NewDecoder(resp.Body).Decode(&r)
	if err != nil {
		return err
	}
	c.apiKey = r.Access_token
	c.apiKeyExpiry = time.Now().Add(time.Duration(r.Expires_in) * time.Second)
	log.Default().Printf("API Key expires in %v\n", time.Until(c.apiKeyExpiry))
	return nil
}

// cardPicker is the entry point for the client server that speaks to the blizzard api.
// It returns a slice of 10 cards that match the criteria passed to it.
func (c *client) CardPicker() ([]Card, error) {
	err := c.getAPIKey()
	if err != nil {
		return nil, err
	}
	cards, err := c.getCard()
	if err != nil {
		return nil, err
	}
	rand.Shuffle(len(cards), func(i, j int) {
		cards[i], cards[j] = cards[j], cards[i]
	})
	cards = cards[:10]
	sort.SliceStable(cards, func(i, j int) bool {
		return cards[i].ID < cards[j].ID
	})
	log.Default().Printf("Pulled and sorted %v cards\n", len(cards))
	return cards, nil
}

// concatCriteria takes a criteria struct and turns them into a string to be appended to the url.
func concatCriteria(c *criteria) string {
	var result string
	for k, v := range *c {
		// if the value is a slice, append each value to the url.
		// This allows us to append multiple values to the same query parameter.
		if reflect.TypeOf(v).Kind() == reflect.Slice {
			result += "&" + k + "="
			for _, i := range v.([]any) {
				result += fmt.Sprintf("%v", i) + ","
			}
			// trim the trailing comma.
			// This doesn't effect the efficacy of the url, but it looks nicer.
			result = strings.TrimSuffix(result, ",")
			continue
		} else {
			// append the search criteria to the url. format k=v
			result += "&" + k + "=" + fmt.Sprintf("%v", v)
		}
	}
	return result
}
