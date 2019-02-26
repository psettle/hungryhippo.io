package hhclientmanager

import (
	"fmt"
	"log"
	"strconv"

	"hungryhippo.io/go-src/hhdatabase"

	"github.com/bitly/go-simplejson"
	uuid "github.com/satori/go.uuid"
)

const (
	newPlayerRequest  = iota //A new player requests to join the game
	newPlayerResponse = iota //Server acknowledging new player request, provided initial condition

	positionUpdateRequest = iota //A player asks to be moved to a new location
	positionUpdateMessage = iota //Server tells player about the position of all players

	consumeFruitRequest = iota //A player asks to consume an existing fruit
	consumeFruitMessage = iota //server notifies clients that a fruit has died
	newFruitMessage     = iota //the server has generated a new fruit

	consumePlayerRequest  = iota //a player asks to consume another player
	consumePlayerResponse = iota //the server accepts/denies the consumption request
	playerDeathMessage    = iota //server notififies a player that they have died, and must submit a new newPlayerRequest
)

func createNewPlayerResponse(player *hhdatabase.Player) (*simplejson.Json, error) {
	return simplejson.NewJson([]byte(`{
		"type" : ` + fmt.Sprintf("%d", newPlayerResponse) + `,
		"data" : {
			"id" : "` + player.ID.String() + `",
			"points" : 0,
			"location" : {
				"centre": {
					"x": ` + fmt.Sprintf("%f", player.Location.Centre.X) + `,
					"y": ` + fmt.Sprintf("%f", player.Location.Centre.Y) + `
				},
				"direction": ` + fmt.Sprintf("%f", player.Location.Direction) + `
			}
		}
	}`))
}

func createPositionUpdateMessage() (*simplejson.Json, error) {
	conn, connErr := hhdatabase.BeginOperation()
	if connErr != nil {
		return nil, connErr
	}
	defer hhdatabase.EndOperation(conn)

	players, exists, err := hhdatabase.LoadMany(hhdatabase.Player{}, nil, conn)
	if err != nil {
		return nil, err
	}

	message, jsonErr := simplejson.NewJson([]byte(`{
		"type" : ` + strconv.Itoa(positionUpdateMessage) + `,
		"data" : {
			"count" : ` + strconv.Itoa(len(players)) + `
		}
	}`))

	//shouldn't fail, it's a static type
	if jsonErr != nil {
		log.Panic(err)
	}

	var playerEntries []*simplejson.Json

	for i := range players {
		item := *players[i]
		player := item.(hhdatabase.Player)
		exist := exists[i]

		if exist {
			playerEntries = append(playerEntries, playerToSimplejson(&player))
		}
	}

	message.Get("data").Set("players", playerEntries)

	return message, nil
}

func createFruitConsumptionMessage(playerid *uuid.UUID, fruitid *uuid.UUID) (*simplejson.Json, error) {
	return simplejson.NewJson([]byte(`{
		"type" : ` + fmt.Sprintf("%d", consumeFruitMessage) + `,
		"data" : {
			"consumer_id" : "` + playerid.String() + `",
			"consumed_id" : "` + fruitid.String() + `"
		}
	}`))
}

func createNewFruitMessage(fruit *hhdatabase.Fruit) (*simplejson.Json, error) {
	return simplejson.NewJson([]byte(`{
		"type" : ` + fmt.Sprintf("%d", newFruitMessage) + `,
		"data" : {
			"id" : "` + fruit.ID.String() + `",
			"position": {
				"x": ` + fmt.Sprintf("%f", fruit.Position.X) + `,
				"y": ` + fmt.Sprintf("%f", fruit.Position.Y) + `
			},
		}
	}`))
}

func createConsumePlayerResponse(playerid *uuid.UUID, score int) (*simplejson.Json, error) {
	return simplejson.NewJson([]byte(`{
		"type" : ` + fmt.Sprintf("%d", consumePlayerResponse) + `,
		"data" : {
			"id" : "` + playerid.String() + `",
			"points":` + strconv.Itoa(score) + `,
		}
	}`))
}

func createPlayerDeathMessage(playerid *uuid.UUID) (*simplejson.Json, error) {
	return simplejson.NewJson([]byte(`{
		"type" : ` + fmt.Sprintf("%d", playerDeathMessage) + `,
		"data" : {
			"id" : "` + playerid.String() + `"
		}
	}`))
}

func playerToSimplejson(player *hhdatabase.Player) *simplejson.Json {
	json, err := simplejson.NewJson([]byte(`{
		"id" : "` + player.ID.String() + `",
		"points" : 0,
		"location" : {
			"centre": {
				"x": ` + fmt.Sprintf("%f", player.Location.Centre.X) + `,
				"y": ` + fmt.Sprintf("%f", player.Location.Centre.Y) + `
			},
			"direction": ` + fmt.Sprintf("%f", player.Location.Direction) + `
		}
	}`))

	//shouldn't fail, it's a static type
	if err != nil {
		log.Panic(err)
	}

	return json
}
