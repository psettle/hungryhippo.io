package hhclientmanager

import (
	"fmt"
	"log"
	"strconv"

	"hungryhippo.io/go-src/hhdatabase"

	"github.com/bitly/go-simplejson"
)

const (
	newPlayerRequest  = iota //A new player requests to join the game
	newPlayerResponse = iota //Server acknowledging new player request, provided initial condition

	gamestateUpdateMessage = iota //Server tells player about the position of all players and fruits

	positionUpdateRequest = iota //A player asks to be moved to a new location
	consumeFruitRequest   = iota //A player asks to consume an existing fruit
	consumePlayerRequest  = iota //a player asks to consume another player
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

func createGamestateUpdateMessage() (*simplejson.Json, error) {
	conn, connErr := hhdatabase.BeginOperation()
	if connErr != nil {
		return nil, connErr
	}
	defer hhdatabase.EndOperation(conn)

	//Note: This load is not transactional, so it might miss new players or include recently dead players very briefly
	//      Scores for players may also be incorrect very briefly
	players, playerExists, playerErr := hhdatabase.LoadMany(hhdatabase.Player{}, nil, conn)
	if playerErr != nil {
		return nil, playerErr
	}

	//Note: This load is not transactional, so it might miss new fruits or include recently consumed fruits very briefly
	//      Scores for players may also fail to reflect recently consumed fruits
	fruits, fruitExists, fruitErr := hhdatabase.LoadMany(hhdatabase.Fruit{}, nil, conn)
	if fruitErr != nil {
		return nil, fruitErr
	}

	message, jsonErr := simplejson.NewJson([]byte(`{
		"type" : ` + strconv.Itoa(gamestateUpdateMessage) + `,
		"data" : {
			"players" : {
				"count" : ` + strconv.Itoa(len(players)) + `
			},
			"fruits" : {
				"count" : ` + strconv.Itoa(len(fruits)) + `
			}		
		}
	}`))

	//shouldn't fail, it's a static type
	if jsonErr != nil {
		log.Panic(playerErr)
	}

	var playerEntries []*simplejson.Json
	for i := range players {
		item := *players[i]
		player := item.(hhdatabase.Player)
		exist := playerExists[i]

		if exist {
			playerEntries = append(playerEntries, playerToSimplejson(&player))
		}
	}

	var fruitEntries []*simplejson.Json
	for i := range fruits {
		item := *fruits[i]
		fruit := item.(hhdatabase.Fruit)
		exist := fruitExists[i]

		if exist {
			fruitEntries = append(fruitEntries, fruitToSimplejson(&fruit))
		}
	}

	message.Get("data").Get("players").Set("elements", playerEntries)
	message.Get("data").Get("fruits").Set("elements", fruitEntries)

	return message, nil
}

func playerToSimplejson(player *hhdatabase.Player) *simplejson.Json {
	json, err := simplejson.NewJson([]byte(`{
		"id" : "` + player.ID.String() + `",
		"points" : "` + fmt.Sprintf("%d", player.Score) + `",
		"nickname" : "` + player.Name + `",
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

func fruitToSimplejson(fruit *hhdatabase.Fruit) *simplejson.Json {
	json, err := simplejson.NewJson([]byte(`{
		"id" : "` + fruit.ID.String() + `",
		"position" : {
			"x": ` + fmt.Sprintf("%f", fruit.Position.X) + `,
			"y": ` + fmt.Sprintf("%f", fruit.Position.Y) + `
		}
	}`))

	//shouldn't fail, it's a static type
	if err != nil {
		log.Panic(err)
	}

	return json
}
