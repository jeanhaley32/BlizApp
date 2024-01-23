package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
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
	pageLimit = 1
	locale    = "en_US"
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
	default:
		return "unknown"
	}
}

type CardsResponse struct {
	Cards []Card `json:"cards"`
}

// Query Hearthstone api for cards matching a set criteria.
func (c *client) GetCard() ([]Card, error) {
	var CardPages []CardsResponse
	for i := 1; i <= pageLimit; i++ {
		url := apiURL + "?locale=" + locale + "&access_token=" + c.apiKey + "&page=" + fmt.Sprintf("%v", i)
		for k, v := range c.criteria {
			if reflect.TypeOf(v).Kind() == reflect.Slice {
				url += "&" + k + "="
				for _, i := range v.([]any) {
					url += fmt.Sprintf("%v", i) + ","
				}
				url = strings.TrimSuffix(url, ",")
				continue
			} else {
				// append the search criteria to the url. format k=v
				url += "&" + k + "=" + fmt.Sprintf("%v", v)
			}
		}
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
		fmt.Println(resp)
		cardPage := CardsResponse{}
		err = json.NewDecoder(resp.Body).Decode(&cardPage)
		if err != nil {
			return nil, err
		}
		CardPages = append(CardPages, cardPage)
	}
	if len(CardPages) == 0 {
		return nil, fmt.Errorf("no cards returned from the API")
	}

	var r CardsResponse
	for _, page := range CardPages {
		r.Cards = append(r.Cards, page.Cards...)
	}
	return r.Cards, nil
}

func (c *client) GetAPIKey() error {

	if c.apiKey != "" && time.Now().Before(c.apiKeyExpiry) {
		return fmt.Errorf("API Key is still valid")
	}

	if c.secrets.clientID == "" || c.secrets.secret == "" {
		return fmt.Errorf("ClientID or Secret is not set")
	}

	type Response struct {
		Access_token string `json:"access_token"`
		Expires_in   int    `json:"expires_in"`
	}

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
	return nil
}

func cardPicker(s secrets, c criteria) ([]Card, error) {
	client := client{
		secrets:  &s,
		criteria: c,
	}
	client.GetAPIKey()
	cards, err := client.GetCard()
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
	return cards, nil
}
