package hhclientmanager

import (
	"fmt"
	"log"
	"strconv"

	"hungryhippo.io/go-src/hhdatabase"

	"github.com/bitly/go-simplejson"
)

const (
	newPlayerRequest      = iota
	newPlayerResponse     = iota
	positionUpdateRequest = iota
	positionUpdateMessage = iota
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
	players, exists, err := hhdatabase.LoadPlayers(nil, nil)
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
		player := players[i]
		exist := exists[i]

		if exist {
			playerEntries = append(playerEntries, playerToSimplejson(player))
		}
	}

	message.Get("data").Set("players", playerEntries)

	return message, nil
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
