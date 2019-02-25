package hhdatabase

import (
	"github.com/mediocregopher/radix.v2/redis"
	uuid "github.com/satori/go.uuid"
)

//Item is an interface to a collection of related fields in redis
//Each item is identified with a UUID
type Item interface {
	//getWatchKey returns a key used as a 'watch variable' in redis to prevent illegal transactions
	getWatchKey() string

	//getMembersKey returns a key to the set of valid item instances
	getMembersKey() string

	//getValueKey returns a key to the values of item instances
	getValueKey() string

	//create returns a pointer to an empty item intialized with the provided uuid
	create(uuid *uuid.UUID) Item

	///listKeys returns a list of item database keys for members
	listKeys() []string

	//mapKeys returns a map of database key => item value
	mapKeys() map[string]interface{}

	//fromList parses a list of item values into members.
	//list is provided in the same order as listKeys() returns
	fromList(list []string) (Item, error)

	//fromMap parses a map of database key => item value into members
	fromMap(entries map[string]string) (Item, error)
}

//Watch puts a watch on a item entry, critical for load -> edit -> save operations,
//requires a close call after the transaction is complete.
func Watch(i Item, conn *redis.Client) error {
	reply := conn.Cmd("WATCH", i.getWatchKey())
	return reply.Err
}

//Load gets a single item from the database, returns true if the entry existed, false on not exists or error
func Load(i Item, conn *redis.Client) (Item, bool, error) {
	var connI interface {
		Cmd(cmd string, args ...interface{}) *redis.Resp
	}

	reply, err := connI.Cmd("HMGET", i.getValueKey(), i.listKeys()).List()
	if err != nil {
		return i, false, err
	}

	//copy result to object
	i, err = i.fromList(reply)
	if err != nil {
		//error indicates the list was invalid, and since we know the command didn't fail, the item must not exist
		return i, false, nil
	}

	//no error
	return i, true, nil
}

//WatchMany starts a watch on a set of item entries
//
// returns a client with which the watch started, can be used in LoadPlayers to begin a bulk operation
// must be closed using endOperation(conn) to free up the allocated connection
//
//if items == nil, a watch will be put on all items
func WatchMany(itemtype Item, items []*Item, conn *redis.Client) error {
	//start with a watch on the item set
	reply := conn.Cmd("WATCH", itemtype.getMembersKey())
	if reply.Err != nil {
		return reply.Err
	}

	var idStrings []string
	if items == nil {
		//ids was nil, load all members into a list
		var err error
		idStrings, err = conn.Cmd("SMEMBERS", itemtype.getMembersKey()).List()
		if err != nil {
			return err
		}
	} else {
		//else put all provided ids into a list
		for _, element := range items {
			idStrings = append(idStrings, (*element).getWatchKey())
		}
	}

	//put a watch on all members
	reply = conn.Cmd("WATCH", idStrings)
	if reply.Err != nil {
		return reply.Err
	}

	return nil
}

//LoadMany returns
//a set of item objects for each id in ids
//a exists boolean for each id in ids
//a general error state for the operation
//
//if ids == nil, all items in the database will be returned
//if conn != nil, it will be used to execute load commands, this is to maintain a watch operation, the operation will be terminated on error
func LoadMany(itemtype Item, items []*Item, conn *redis.Client) ([]*Item, []bool, error) {
	var err error

	//get ids to load
	var idStrings []string
	if items == nil {
		//ids was nil, load all members into a list
		idStrings, err = conn.Cmd("SMEMBERS", itemtype.getMembersKey()).List()
		if err != nil {
			return nil, nil, err
		}
	} else {
		//else put all provided ids into a list
		for _, element := range items {
			idStrings = append(idStrings, (*element).getWatchKey())
		}
	}

	//fetch item data
	var retItems []*Item
	var exists []bool

	if items == nil {
		//ids was nil, load all item data into a map, then let each item parse their own members
		entries, allErr := conn.Cmd("HGETALL", itemtype.getValueKey()).Map()
		if allErr != nil {
			return nil, nil, allErr
		}

		for _, id := range idStrings {
			uuid, uuidErr := uuid.FromString(id)
			if uuidErr != nil {
				return nil, nil, uuidErr
			}

			i := itemtype.create(&uuid)
			i, err = i.fromMap(entries)

			if err == nil {
				retItems = append(retItems, &i)
				exists = append(exists, true)
			}
			//(Don't add invalid item entries)
		}
	} else {
		//ids was not nil, load all ids into a list, then let each item parse their own members
		var keys []string

		for _, id := range idStrings {
			uuid, uuidErr := uuid.FromString(id)
			if uuidErr != nil {
				return nil, nil, uuidErr
			}

			i := itemtype.create(&uuid)
			retItems = append(retItems, &i)

			keys = append(keys, i.listKeys()...)
		}

		list, allErr := conn.Cmd("HMGET", itemtype.getValueKey(), keys).List()
		if allErr != nil {
			return nil, nil, allErr
		}

		loc := 0
		for _, i := range retItems {
			(*i), err = (*i).fromList(list[loc : loc+len((*i).listKeys())])
			exists = append(exists, err == nil)

			loc += len((*i).listKeys())
		}
	}

	return retItems, exists, nil
}

//Save saves a item entry to the database
//Note: should be used with watches/transactions, which are managed seperately on the provided conn
func Save(i Item, conn *redis.Client) error {
	//put the item id into the set of item
	response := conn.Cmd("SADD", i.getMembersKey(), i.getWatchKey())
	if response.Err != nil {
		return response.Err
	}

	//save the item entry
	response = conn.Cmd("HMSET", i.getValueKey(), i.mapKeys())
	if response.Err != nil {
		return response.Err
	}

	//write to the item watch, to prevent concurrent access
	response = conn.Cmd("SET", i.getWatchKey(), "")
	if response.Err != nil {
		return response.Err
	}

	//no errors
	return nil
}

//Delete removes a item entry from the database
//
//Note: should be used with watches/transactions, which are managed seperately on the provided conn
//Note: error = nil does NOT imply successful deletion, the item may have already been deleted
//  	to check if an item was deleted, watch the item, load it to confirm it exists, then delete in a transaction
func Delete(i *Item, conn *redis.Client) error {
	item := *i

	//remove the item from the set of items
	response := conn.Cmd("SREM", item.getMembersKey(), item.getWatchKey())
	if response.Err != nil {
		return response.Err
	}

	//delete the item entry
	response = conn.Cmd("HDEL", item.getValueKey(), item.listKeys())
	if response.Err != nil {
		return response.Err
	}

	//delete the item watch, to prevent concurrent access
	response = conn.Cmd("DEL", item.getWatchKey())
	if response.Err != nil {
		return response.Err
	}

	//no errors
	return nil
}
