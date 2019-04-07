package hhclientmanager

import (
	"math"

	"hungryhippo.io/go-src/hhdatabase"
)

const xBoardWidth = 1000.0
const yBoardWidth = 1000.0
const maxDirection = math.Pi * 2

//Insert a new player and a new fruit into the database
//
//returns true if the player and fruit were created, false otherwise
func createNewPlayer(player *hhdatabase.Player, fruit *hhdatabase.Fruit) (bool, error) {
	//we need to save the player, so begin an operation
	conn, err := hhdatabase.BeginOperation()
	if err != nil {
		return false, err
	}
	defer hhdatabase.EndOperation(conn)

	//put a watch on the player to avoid conflicts
	err = hhdatabase.Watch(player, conn)
	if err != nil {
		return false, err
	}

	//put a watch on the fruit to avoid conflicts
	err = hhdatabase.Watch(fruit, conn)
	if err != nil {
		return false, err
	}

	//load the player to check if it exists already
	var exists bool
	_, exists, err = hhdatabase.Load(player, conn)
	if err != nil {
		return false, err
	}
	if exists {
		return false, nil
	}

	//load the fruit to check if it exists already
	_, exists, err = hhdatabase.Load(fruit, conn)
	if err != nil {
		return false, err
	}
	if exists {
		return false, nil
	}

	//start a transaction to prepare for a save
	err = hhdatabase.Multi(conn)
	if err != nil {
		return false, err
	}

	//queue the save operation on the player
	err = hhdatabase.Save(player, conn)
	if err != nil {
		return false, err
	}

	//queue the save operation on the fruit
	err = hhdatabase.Save(fruit, conn)
	if err != nil {
		return false, err
	}

	//execute the queued operation
	var applied bool
	applied, err = hhdatabase.Exec(conn)
	if err != nil {
		return false, err
	}

	return applied, nil
}

func deletePlayer(player *hhdatabase.Player) (bool, error) {
	//we need to delete the player, so begin an operation
	conn, err := hhdatabase.BeginOperation()
	if err != nil {
		return false, err
	}
	defer hhdatabase.EndOperation(conn)

	//infinite loop for retries
	//will exit if
	//- the player and a fruit is deleted successfully
	//- only the player is deleted successfully (there are no fruits to delete)
	//- the player no longer exists
	//- a redis operation fails (implies the server is not available)
	for {
		//start a watch so the delete is valid
		err = hhdatabase.Watch(player, conn)
		if err != nil {
			return false, err
		}

		//also need a watch on the fruit type so if we fail to get one, we are sure it's because there are none left
		err = hhdatabase.WatchType(hhdatabase.Fruit{}, conn)
		if err != nil {
			return false, err
		}

		//check that the player exists so we can return applied correctly
		//(applied will  be true on delete if the player never existed, which gives a poor impression to the caller)
		_, exists, loadErr := hhdatabase.Load(player, conn)
		if loadErr != nil {
			return false, loadErr
		}

		if !exists {
			return false, nil
		}

		//fetch a random member of fruits
		var item hhdatabase.Item
		item, exists, err = hhdatabase.LoadRandom(hhdatabase.Fruit{}, conn)
		fruit := item.(hhdatabase.Fruit)

		//try to delete the player and fruit
		err = hhdatabase.Multi(conn)
		if err != nil {
			return false, err
		}

		err = hhdatabase.Delete(player, conn)
		if err != nil {
			return false, err
		}

		//only try to delete the fruit if we are sure it exists
		if exists {
			err = hhdatabase.Delete(fruit, conn)
			if err != nil {
				return false, err
			}
		}

		var applied bool
		applied, err = hhdatabase.Exec(conn)
		if err != nil {
			return false, err
		}

		if applied {
			return true, nil
		}
	}
}

//Update a player position in the database
//
//returns true if the position was updated, false otherwise
func updatePlayerPosition(player *hhdatabase.Player, newX float64, newY float64, newDirection float64) (bool, error) {
	//We need to update the player, start an operation
	conn, err := hhdatabase.BeginOperation()
	if err != nil {
		return false, err
	}
	defer hhdatabase.EndOperation(conn)

	//infinite loop for retries
	//will exit if
	//- the position update is successfully applied
	//- the position update is deemed to be illegal
	//- a redis operation fails (implies the server is not available)
	for {
		//put a watch on the player
		err = hhdatabase.Watch(player, conn)
		if err != nil {
			return false, err
		}

		//Load the player
		var item hhdatabase.Item
		var exists bool
		item, exists, err = hhdatabase.Load(player, conn)
		if err != nil {
			return false, err
		}

		if !exists {
			//player doesn't exist
			return false, nil
		}

		//TODO: validate that movement was legal... (player didn't collide/get collided with)

		//save the new player position
		loadedPlayer := item.(hhdatabase.Player)
		player = &loadedPlayer

		player.Location.Centre.X = newX
		player.Location.Centre.Y = newY
		player.Location.Direction = newDirection

		err = hhdatabase.Multi(conn)
		if err != nil {
			return false, err
		}

		err = hhdatabase.Save(player, conn)
		if err != nil {
			return false, err
		}

		var applied bool
		applied, err = hhdatabase.Exec(conn)
		if err != nil {
			return false, err
		}

		if applied {
			return true, nil
		}

		//(retry if change wasn't applied)
	}
}

func consumeFruit(player *hhdatabase.Player, fruit *hhdatabase.Fruit, newFruit *hhdatabase.Fruit) (bool, error) {
	//start an operation for the consumption
	conn, err := hhdatabase.BeginOperation()
	if err != nil {
		return false, err
	}
	defer hhdatabase.EndOperation(conn)

	//Infinite loop for applying fruit consumption, ends if
	//- redis operation fails (implies database is not accessible)
	//- operation is deemed invalid (player or fruit missing)
	//- fruit successfully consumed & new fruit created
	for {

		//put watches on relevant items
		err = hhdatabase.Watch(fruit, conn)
		if err != nil {
			return false, err
		}

		err = hhdatabase.Watch(newFruit, conn)
		if err != nil {
			return false, err
		}

		err = hhdatabase.Watch(player, conn)
		if err != nil {
			return false, err
		}

		var exists bool
		_, exists, err = hhdatabase.Load(fruit, conn)
		if err != nil {
			return false, err
		}

		if !exists {
			//fruit doesn't exist, must have been consumed already
			return false, nil
		}

		var item hhdatabase.Item
		item, exists, err = hhdatabase.Load(player, conn)
		if err != nil {
			return false, err
		}

		if !exists {
			//player doesn't exist, may have died already
			return false, nil
		}

		var playerItem hhdatabase.Player
		playerItem = item.(hhdatabase.Player)
		player = &playerItem

		player.Score++ //all fruits are worth one

		//start the transaction now that we know player and fruit state
		err = hhdatabase.Multi(conn)
		if err != nil {
			return false, err
		}

		err = hhdatabase.Save(player, conn)
		if err != nil {
			return false, err
		}

		err = hhdatabase.Delete(fruit, conn)
		if err != nil {
			return false, err
		}

		err = hhdatabase.Save(newFruit, conn)
		if err != nil {
			return false, err
		}

		var applied bool
		applied, err = hhdatabase.Exec(conn)
		if err != nil {
			return false, err
		}

		if applied {
			return true, nil
		}

		//(Else retry operation)
	}
}

func consumePlayer(consumer *hhdatabase.Player, consumed *hhdatabase.Player) (bool, error) {
	//start an operation for the consumption
	conn, err := hhdatabase.BeginOperation()
	if err != nil {
		return false, err
	}
	defer hhdatabase.EndOperation(conn)

	//Infinite loop for applying player consumption, ends if
	//- redis operation fails (implies database is not accessible)
	//- operation is deemed invalid (player missing)
	//- player successfully consumed
	for {

		//put watches on relevant items
		err = hhdatabase.Watch(consumer, conn)
		if err != nil {
			return false, err
		}

		err = hhdatabase.Watch(consumed, conn)
		if err != nil {
			return false, err
		}

		var exists bool
		var item hhdatabase.Item
		var playerItem hhdatabase.Player

		item, exists, err = hhdatabase.Load(consumer, conn)
		if err != nil {
			return false, err
		}

		if !exists {
			//player doesn't exist, must have been consumed already
			return false, nil
		}

		playerItem = item.(hhdatabase.Player)
		consumerItem := playerItem

		item, exists, err = hhdatabase.Load(consumed, conn)
		if err != nil {
			return false, err
		}

		if !exists {
			//player doesn't exist, must have been consumed already
			return false, nil
		}

		playerItem = item.(hhdatabase.Player)
		consumedItem := playerItem

		//TODO: validate that consumption is allowed

		//apply the consumption
		consumerItem.Score += consumedItem.Score

		//start the transaction now that we know player and fruit state
		err = hhdatabase.Multi(conn)
		if err != nil {
			return false, err
		}

		//update the consumer
		err = hhdatabase.Save(consumerItem, conn)
		if err != nil {
			return false, err
		}

		//delete the consumed
		err = hhdatabase.Delete(consumedItem, conn)
		if err != nil {
			return false, err
		}

		var applied bool
		applied, err = hhdatabase.Exec(conn)
		if err != nil {
			return false, err
		}

		if applied {
			//copy consumer/consumed over in case caller want the new scores
			*consumer = consumerItem
			*consumed = consumedItem

			return true, nil
		}

		//(Else retry operation)
	}
}
