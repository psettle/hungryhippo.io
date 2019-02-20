package hhclientmanager

import (
	"fmt"

	"github.com/bitly/go-simplejson"
	uuid "github.com/satori/go.uuid"
)

func createNewPlayerResponse(clientID *uuid.UUID, xpos float64, ypos float64, direction float64) (*simplejson.Json, error) {
	return simplejson.NewJson([]byte(`{
		"type" : "NewPlayerResponse",
		"data" : {
			"id" : "` + clientID.String() + `",
			"points" : 0,
			"location" : {
				"centre": {
					"x": ` + fmt.Sprintf("%f", xpos) + `,
					"y": ` + fmt.Sprintf("%f", ypos) + `
				},
				"direction": ` + fmt.Sprintf("%f", direction) + `
			}
		}
	}`))
}
