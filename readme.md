# BlizApp
BlizApp is my solution to Blizzard's take-home project as part of their take-home application process.  

I chose `Option 1: Software Engineering`, as I've been working on several active Golang projects of my own, and I feel like this would be a logical use of fresh information. 

# About

The primary goal of the project is to communicate with the `Hearthstone API` to obtain several cards with set criteria and display those cards, sorted by card ID in a web app. 

## Criteria
- Obtain `10` cards
- either class `Druid` or `Warlock`
- `Mana` of at least `seven`
- Legendary Rarity. 
- Sort by `Card ID` in a table, and display the following in a human-readable table. 
    - `Card Image`
    - `Name`
    - `Type`
    - `Rarity`
    - `Set`
    - `Class`


 ### Abstraction
This application needs to do several things. 
- Host a Web Server that returns a formatted website with `ten` ID-sorted Hearthstone cards. 
    - Handle Get requests for `/`
    - Using a `secret` and `Client ID`, maintain an `API Key` with a halflife of `24` hours.
    - Obtain a `deck` of `ten cards` that meet the criteria listed in the `criteria` section above.
    - Generate a web view of these cards, sorting them by Card ID
        - Each `Card` must display the card's `image`, `Name`, `Type`, `Rarity`, `Set`, and `Class`
    - Pass this HTML and CSS code to the user. 

> Here is the final struct used to house needed card data received from the Hearthstone API
``` Golang
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

type CardsResponse struct {
	Cards     []Card `json:"cards"`
	PageCount int    `json:"pageCount"`
}

```


 ## Proposed Solutions
 Let's break this up into constituent components. 
 1. `Web Server` - Handle Get Requests.
 2. `API Client` - Negotiate connection with `Hearthstone` API.
 3. `Site renderer` - Construct a Site based on the `deck` received from API, and a preconstructed template. 

 ### Web Server
The Webserver in this solution is a simple multiplexer `Goji` That receives an incoming GET request and returns a rendered webpage containing ten ID-sorted cards that meet the passed `criteria`

### API Access Key
the API Client has the`.`GetAPIKey` method.\
This method
Check if an API key exists, and if it's still valid.
2. Uses a stored `secret` and `client ID` to obtain a new `API key` if it fails those checks. 

> secret and client ID are either passed via Command Line flags, or read from `secrets.json`, a json file stored at the codes root directory, with the following construction.

I understand that this solution may not be an industry standard or ideal one. 

Since secret keys are stored server-side, along with pages being prerendered before going to the user, I don't believe there is an opportunity for them to be revealed. 

Passing them as flags is useful for containerization, as long as you pass those values carefully, and store them in an encrypted fashion. 

Storing them as JSON files is probably the least secure (how I am doing it), a future implementation of this would have me going down a long rabbit hole, and coming out the other end with a much more elegant solution. 

``` json
{
  "clientid": "<client id>",
  "secret" :  "<secret>"
}
```

 ### API Client
```go
type criteria map[string]any


type client struct {
	secrets      *secrets
	apiKey       string
	apiKeyExpiry time.Time
	criteria     criteria
}

type secrets struct {
	ClientID string `json:"clientid"`
	Secret   string `json:"secret"`
}

```
Our API Client contains a `secrets` struct, that holds the `clientID` and `Secret`, along with the `API key` and its `expiration time`, and `criteria`, a `key` `value` map of `string` to `any` type. 

> For this project, the instantiation of a `criteria` object is called `params`, and contains the search criteria set by the project objective. 
``` go
params = map[string]any{
    "sort":     "ID:asc",
    "manaCost": 7,
    "rarity":   legendary,
    "class":    []any{warlock, druid},
}
```
the client struct's `GetCard()` method is used to construct the appropriate URLs to obtain cards from the Hearthstone API server.

- This method returns a slice of cards with the information that we need to meet this project's objectives. 
- We construct the `outbound URL` and append all of the necessary search parameters from the `criteria` object, this function supports multiple values for a single search parameter by using reflection to see when an entry is a slice and appends both values to the URL as comma-separated entries. 
 
 ### Server-Side Renderer
 This is the portion of the assignment I was least familiar with and learned a lot from. If I were to put more work into this, I'd make the containers for the cards more uniform and cleaner. It could use some work, but for the sake of pushing this through to the end, I'm labeling this extra work as out of scope for this assignment. 
- I adapted CSS I had pre-written for a different project to format and display the card information used here.  
- I used the `aymerick/raymond` golang `handlebars` library to dynamically construct entries for each card.
> I chose this `handlebars` solution because it solved what I was initially going to construct by hand with a very simple and easy execution.
## Action Plan
resources: [Getting Started](https://develop.battle.net/documentation/guides/getting-started), [API Guides](https://develop.battle.net/documentation/hearthstone/guides)

| Task | Status |
|------|--------|
| > Research Hearthstone API Documentation| Done | 
| 1. Learn Endpoints and authentication requirements | Done |
| 2. Determine how to securely store the API secret for communication | Done |
| > Implement Web Server |  done |
| 1. Use `Goji` to handle `HTTP` `Get` Requests| done  |
| > Implement API CliClient | done |
| 1. Dependant on the primary`research` task above. | done |
| > Create Server-Sider Renderer  | done |
| 1. Design web template to use when rendering `Deck` | done |
| 2. append `cards` in `deck` to template, and construct the rendered site | done |


---
