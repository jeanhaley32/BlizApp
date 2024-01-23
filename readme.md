# BlizApp
BlizApp is my solution to Blizzard's take-home project. Part of their Application process for SRE. 

I chose `Option 1: Software Engineering`, as I've been working on several active Golang projects of my own, and I feel like this would be a logical use of fresh information. Along with a chance to learn something new.  

# About

The project's primary goal is to communicate with the Hearthstone API to obtain several cards with set criteria and display those cards, sorted by card ID in a web app. 

# How to use this
- BlizApp, by default, will run a web app on `localhost:8080`, which can be accessed via the browser.
- It will also log entries to the terminal for certain actions.

> IMPORTANT: to run BlizApp, you must pass a client id and a secret. This can be done in two ways.

*BlizApp `--help` information*
``` bash
  -clientid string
    	clientID for the blizzard api
  -json string
    	json file containing the clientID and secret (default "secrets.json")
  -secret string
    	secret for the blizzard API
```
1. Pass the values via the `clientid` and `secret` flags ``` go run . --clientid=<client id> --secret=<secret>```
2. You can also create a json file named `secrets.json` and append the `client id` and `secret` here.
``` json
{
  "clientid": "<client id>",
  "secret" :  "<secret>"
}
```

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
    - Handle incoming GET requests from clients. 
    - Using a secret and Client ID, maintain an API Key with a half-life of 24 hours.
    - Obtain a deck of ten cards that meet the criteria listed in the criteria section above.
    - Generate a web view of these cards, sorting them by Card ID
        - Each `Card` must display the card's `image`, `Name`, `Type`, `Rarity`, `Set,` and `Class.`
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
 2. `API Client` - Negotiate connection with Hearthstone API.
 3. `Site renderer` - Construct a Site based on the deck received from API, and a preconstructed template. 

 ### Web Server
The Webserver in this solution is a simple multiplexer `Goji` That receives an incoming GET request and returns a rendered webpage containing ten ID-sorted cards that meet the passed `criteria`

### API Access Key
The API Client has the `.GetAPIKey` method. \
This method
- Check if an API key exists, and if it's still valid If those checks fail -> use a stored secret and client ID to obtain a new API key

> Secret and client ID are passed via Command Line flags or read from secrets.json, a JSON file stored at the code's root directory, with the following construction.

Since secret keys are stored server-side, along with pages being prerendered before going to the user, there isn't an opportunity for them to be revealed( I am open to being wrong about this)

Passing them as flags is useful for containerization, as long as you pass those values carefully, and store them in an encrypted fashion. 

Storing them as unencrypted JSON files is the least secure (how I stored these during local testing); a future implementation of this would have me going down a long rabbit hole and coming out the other end with a much more elegant solution.

> I understand these methods are not industry standard and likely contain security holes. This is only for use within the scope of this project. 

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
Our API Client contains a secrets struct that holds the `clientID` and `Secret,` along with the API key, its expiration time, and `criteria`, a key value map of `string` to `any` type. 

> For this project, the instantiation of a criteria object is called params, and contains the search criteria set by the project objective. 
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

 ### Pagination
  The server implementation I created supports paginations. It does this by reading the `pages` field provided by our initial call to the `hearthstone` API and modifying the loops target to this int. 
  This is not necessary for the bounds of this project, as with the filtered restrictions, we only ever see a max of 16 cards, which will only ever fit on one page. But I noticed the pagination and had fun
  solving the problem. I tested this with a much broader filter and verified I was able to process multiple pages; I added a limit (set to ten) because 140 pages take far too long. 

  Since pagination was not really necessary, I stopped myself from further improving on this. But I can think of a few ideas to get more cards in a single request and reduce pagination requests to one or two. I think you might be able to 
  make a request for multiple pages by adding in comma-separated values, that's how I was able to gather a mix of card classes in one request. But I'm not going to mess with this further. 
 
 ### Server-Side Renderer
  - All ten pre-sorted cards are rendered into a page that is then sent to the client. 
  - This page organizes the cards into separate containers, displaying an image for the card on the left, and the data on the right. 
  - CSS for this solution was recycled and modified from a previous project I worked on. 
  > This is probably the aspect of this I know the least amount, I feel there needs to be cleanup done on this code to remove superfluous lines. But, it does the job I need it to do for this project. 
  - Card data is stringified via an enum library defined in `enums.go` 

## Action Plan
resources: [Getting Started](https://develop.battle.net/documentation/guides/getting-started), [API Guides](https://develop.battle.net/documentation/hearthstone/guides)

| Task | Status |
|------|--------|
| > Research Hearthstone API Documentation| Done | 
| 1. Learn Endpoints and authentication requirements | Done |
| 2. Determine how to securely store the API secret for communication | Done |
| > Implement Web Server |  done |
| 1. Use `Goji` to handle `HTTP` `GET` Requests| done  |
| > Implement API CliClient | done |
| 1. Dependant on the primary`research` task above. | done |
| > Create Server-Sider Renderer  | done |
| 1. Design web template to use when rendering `Deck` | done |
| 2. append `cards` in `deck` to template, and construct the rendered site | done |


---
