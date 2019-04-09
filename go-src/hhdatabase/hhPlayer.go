package hhdatabase

import (
	"strconv"

	uuid "github.com/satori/go.uuid"
)

//Coordinate defines a 2D position
type Coordinate struct {
	X float64
	Y float64
}

//Location defines a player position and direction
type Location struct {
	Centre    Coordinate
	Direction float64
}

//Player defines a player score, id and position
type Player struct {
	ID       uuid.UUID
	Name     string
	Score    int
	Location Location
}

//CreatePlayer initializes an empty player with an ID, it can then be loaded, saved, etc.
func CreatePlayer(id *uuid.UUID) *Player {
	p := new(Player)
	p.ID = *id
	return p
}

//getWatchKey returns a key used as a 'watch variable' in redis to prevent illegal transactions
func (p Player) getWatchKey() string {
	return p.ID.String()
}

//getMembersKey returns a key to the set of valid item instances
func (p Player) getMembersKey() string {
	return "players"
}

//getValueKey returns a key to the values of item instances
func (p Player) getValueKey() string {
	return "player_data"
}

//create returns a pointer to an empty item intialized with the provided uuid
func (p Player) create(uuid *uuid.UUID) Item {
	player := CreatePlayer(uuid)

	return *player
}

//listKeys as defined by Item.listKeys
func (p Player) listKeys() []string {
	id := p.ID.String()

	return []string{
		id + ":name",
		id + ":score",
		id + ":x",
		id + ":y",
		id + ":direction"}
}

//mapKeys as defined by Item.mapKeys
func (p Player) mapKeys() map[string]interface{} {
	id := p.ID.String()

	return map[string]interface{}{
		id + ":name":      p.Name,
		id + ":score":     p.Score,
		id + ":x":         p.Location.Centre.X,
		id + ":y":         p.Location.Centre.Y,
		id + ":direction": p.Location.Direction}
}

//fromList as defined by Item.fromList
func (p Player) fromList(list []string) (Item, error) {
	var err error

	p.Name = list[0]
	p.Score, err = strconv.Atoi(list[1])
	if err != nil {
		return p, err
	}

	p.Location.Centre.X, err = strconv.ParseFloat(list[2], 64)
	if err != nil {
		return p, err
	}

	p.Location.Centre.Y, err = strconv.ParseFloat(list[3], 64)
	if err != nil {
		return p, err
	}

	p.Location.Direction, err = strconv.ParseFloat(list[4], 64)
	if err != nil {
		return p, err
	}

	return p, nil
}

//fromMap as defined by Item.fromMap
func (p Player) fromMap(entries map[string]string) (Item, error) {
	var err error
	id := p.ID.String()

	p.Name = entries[id+":name"]
	p.Score, err = strconv.Atoi(entries[id+":score"])
	if err != nil {
		return p, err
	}
	p.Location.Centre.X, err = strconv.ParseFloat(entries[id+":x"], 64)
	if err != nil {
		return p, err
	}
	p.Location.Centre.Y, err = strconv.ParseFloat(entries[id+":y"], 64)
	if err != nil {
		return p, err
	}
	p.Location.Direction, err = strconv.ParseFloat(entries[id+":direction"], 64)
	if err != nil {
		return p, err
	}

	return p, nil
}
