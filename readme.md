# BlizApp
BlizApp is my solution to Blizzard's take-home project as part of their take-home application process.  

I chose `Option 1: Software Engineering`, as I've been working on several active Golang projects of my own, and I feel like this would be a logical use of fresh information. 

# About

The primary goal of the project is to communicate with the `Hearthstone API` to obtain several cards with set criteria and display those cards, sorted by card ID, in a web app. 

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
    - Host an API secret used for communication with the Heartstone API
    - Obtain a `deck` of `ten cards`` that meet the criteria listed above
    - Generate an HTML + CSS `View` of these cards, sorting them by Card ID
        - Each `Card` must display the card's `image`, `Name`, `Type`, `Rarity`, `Set`, and `Class`
    - Pass this HTML and CSS code to the user. 

> In order to conceptualize the data i'm going to be handling, I want to define some struct concepts.
``` Golang

type url, ctype, rarity, class string

type card struct {
    id     uint32
    Image  url
    Name   string
    Type   ctype
    Rarity rarity
    Set    string
    Class  class
}

type deck []card
```
> We can define a custom type for any property that can only be one of a static number of values. Using enums can help when we need to create requests for certain data. I may choose to not do this based on how hard this will be to parse in the case that I am marshalling and de-marshalling json. 


 ## Proposed Solutions
 Let's break this up into constituent components. 
 1. `Web Server` - Handle Get Requests.
 2. `API Client` - Negotiate connection with `Hearthstone` API.
 3. `Site Builder` - Construct a Site based on the `deck` received from API, and a preconstructed template. 

 ### Web Server
> We need to choose a framework for negotiation client connections.
- `Goji`  HTTP Request Multiplexer
    - I have chosen to Use `Goji` to handle these requests.
    - My primary reasons were ease of use and availability of documentation.
- The user path will be simple. \
  <- Receive Get request from Client.  \
  -> Construct an `ID` sorted `Deck` of `cards` that meet the appropriate criteria from the Heartstone API  server. \
  O Construct a view using the 

  ### API Access Key
   - Our web server needs an API key to request data from the Hearthstone API
   - An API key is generated by providing a `client ID`` and `secret`` and lasts for 24 hours.
   -  The following flow should occur for us to use and maintain a proper API key
    - Upon startup, the server should obtain a new `API Token`
    - Upon receiving a new `Get` request from a client, we should verify if the token has expired, and obtain a new one if so. 
    - This way we don't have to update at every pass, and we don't have to micromanage our token to make sure it is always fresh.
    > This is closer to a solution I'd provide given more time to put inot this project. But this time i'm going to pull client id and password from a local file, and grab an API key once.  
   


 ### API Client
 > **Objective** - I'd like to generate a random deck of cards that meets the above criteria every time the webpage is loaded. 
 
 *This _may be _inefficient_ based on_ how long the API client takes to negotiate with _the Heartstone server and receive_ cards. It may be useful to generate a deck ahead of the client, pass in the pre-cached deck, and generate a new one. This would help to more quickly provide the user with a deck, this is stipulated on the premise that only one or two people will use this at a time. _An influx of users would_ likely cause this idea to fail. **Figuring out a more efficient method of providing clients with `decks` is **out of **the **scope**** of this** project.***

ATM Negotiating with the Hearthstone API is an area of least knowledge. My priority should in figuring out how to negotiate with the Heartstone API and obtaining the cards that I need to accomplish the set objective. 
 
 ### Server-Side Renderer
  - We are going to render the final page that goes out to the user based on some preconstructed HTML and CSS Code. 
  - For the cards, I am going to try and modify a solution I created previously for my website (not currently online).

## Action Plan
resources: [Getting Started](https://develop.battle.net/documentation/guides/getting-started)](https://develop.battle.net/documentation/guides/getting-started), [API Guides](https://develop.battle.net/documentation/hearthstone/guides)

| Task | Status |
|------|--------|
| > Research Hearthstone API Documentation| Done | 
| 1. Learn Endpoints and authentication requirements | Done |
| 2. Determine how to securely store the API secret for communication | Done |
| > Implement Web Server |  Working |
| 1. Use `Goji` to handle `HTTP` `Get` Requests| Pending  |
| > Implement API CliClient | Done |
| 1. Dependant on the primary`research` task above. | Pending |
| > Create Server-Sider Renderer  | Pending |
| 1. Design web template to use when rendering `Deck` | Pending |
| 2. append `cards` in `deck` to template, and construct the rendered site | Pending |


---
