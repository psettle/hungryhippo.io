package hhclientmanager

import (
	"fmt"

	"hungryhippo.io/go-src/hhdatabase"

	"github.com/bitly/go-simplejson"
)

func createNewPlayerResponse(player *hhdatabase.Player) (*simplejson.Json, error) {
	return simplejson.NewJson([]byte(`{
		"type" : "NewPlayerResponse",
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
