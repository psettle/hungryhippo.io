package hhdatabase

import (
	"github.com/mediocregopher/radix.v2/redis"
	uuid "github.com/satori/go.uuid"
)

//Fruit defines a fruit score, id and position
type Fruit struct {
	ID       uuid.UUID
	Position Coordinate
	conn     *redis.Client
}

//CreateFruit initializes an empty fruit with an ID, it can then be loaded, saved, etc.
func CreateFruit(id *uuid.UUID) *Fruit {
	p := new(Fruit)
	p.ID = *id
	return p
}

//getWatchKey returns a key used as a 'watch variable' in redis to prevent illegal transactions
func (f Fruit) getWatchKey() string {
	return ""
}

//getMembersKey returns a key to the set of valid item instances
func (f Fruit) getMembersKey() string {
	return ""
}

//getValueKey returns a key to the values of item instances
func (f Fruit) getValueKey() string {
	return ""
}

//create returns a pointer to an empty item intialized with the provided uuid
func (f Fruit) create(uuid *uuid.UUID) Item {
	return *CreateFruit(uuid)
}

///listKeys returns a list of item database keys for members
func (f Fruit) listKeys() []string {
	return nil
}

//mapKeys returns a map of database key => item value
func (f Fruit) mapKeys() map[string]interface{} {
	return nil
}

//fromList parses a list of item values into members.
//list is provided in the same order as listKeys() returns
func (f Fruit) fromList(list []string) (Item, error) {
	return f, nil
}

//fromMap parses a map of database key => item value into members
func (f Fruit) fromMap(entries map[string]string) (Item, error) {
	return f, nil
}
