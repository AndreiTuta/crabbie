package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/calini/crabbie/pkg/strings"
	"github.com/gin-gonic/gin"
)

const CODE_LENGTH = 4

var activeGames map[string]Game
var activePlayers map[string]Player

func main() {
	activeGames = make(map[string]Game)

	r := gin.Default()
	r.GET("/game/", GetNewGame)
	r.GET("/game/:game_type/:code", GetGame)
	r.POST("/game/:game_type/:code/user/:user_name", CreateNewUser)

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")

}

func GetNewGame(c *gin.Context) {

	code := strings.GenerateCode(CODE_LENGTH)

	// get a default, suite based deck of 52 cards
	var d = Deck{map[string][]Card{}}
	d.populate()
	activeGames[code] = Game{code, "1", []Player{}, d, []Turn{}}
	response := fmt.Sprintf("Here's your new game: %s, %s", activeGames[code].Code, activeGames[code].GameType)
	c.String(http.StatusAccepted, response)
}

func CreateNewUser(c *gin.Context) {

	game, found := activeGames[c.Param("code")]

	if found {
		player := Player{c.Param("user_name"), []string{}}
		game.Players = append(game.Players, player)
		activeGames[game.Code] = game
		c.String(http.StatusOK, player.Name)
	} else {
		c.String(http.StatusNotFound, "Could not add user for game!")
	}
}

func GetGame(c *gin.Context) {

	game, found := activeGames[c.Param("code")]
	if found {
		playTurn(game.Deck)
		e, err := json.Marshal(game)
		if err != nil {
			fmt.Println(err)
			return
		}
		c.String(http.StatusOK, string(e))
	} else {
		c.String(http.StatusNotFound, "No game was found!")
	}
}

type Player struct {
	Name string   `json:"player_name"`
	Card []string `json:"player_cards"`
}

type Turn struct {
	Id            int               `json:turn_number`
	PlayerActions map[string][]Card `json:"player_actions"`
}

func (t Turn) playTurn(d Deck) Deck {
	for suite, values := range d.Cards {
		if suite == "Spades" {
			for index, card := range values {
				if card.Rank == 15 {
					fmt.Println("Removing card index: ", index)
					d.removeCard(card)
				}
			}
		}
	}
	return d
}

type Card struct {
	Rank int    `json:"card_rank"`
	Suit string `json:"card_suit"`
}

type Deck struct {
	Cards map[string][]Card `json:"cards_in_deck"`
}

func (c Deck) populate() {

	var ranks = []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}
	var suits = []string{"Clovers ", "Diamonds", "Hearts", "Spades"}

	for sec_index, suit := range suits {
		var cards = c.Cards[suit]
		for index, rank := range ranks {
			var card = Card{rank, suit}
			cards = append(cards, card)
			fmt.Println("Rank/suit indexes: ", index, sec_index)

			// fmt.Println(cards)
		}
		c.Cards[suit] = cards

	}
	// fmt.Println(c.Cards)
}

func (d Deck) removeCard(c Card) {
	var vals []Card = d.Cards[c.Suit]
	for i := range vals {
		if vals[i].Rank == c.Rank {
			fmt.Print("Removing card ", i)
			if i == 0 {
				vals = vals[i:]
			} else if i == len(vals)-1 {
				vals = vals[:i]
			} else {
				vals[i] = Card{}
			}
			fmt.Println("Updating cards after removing ", c, vals)
		}
	}
	d.Cards[c.Suit] = vals
}

type Game struct {
	Code     string   `json:"game_code"`
	GameType string   `json:"game_type"`
	Players  []Player `json:"players_in_game"`
	Deck     Deck     `json:"deck_cards"`
	Turns    []Turn   `json:game_turns`
}
