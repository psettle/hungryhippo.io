var messageTypeVal = 0
var messageType = {
  newPlayerRequest      : messageTypeVal++, //A new player requests to join the game 
  newPlayerResponse     : messageTypeVal++, //Server acknowledging new player request, provided initial condition
  
  gamestateUpdateMessage : messageTypeVal++, //Server tells player about the position of all players and fruits

  positionUpdateRequest : messageTypeVal++, //A player asks to be moved to a new location
  consumeFruitRequest  : messageTypeVal++, //A player asks to consume an existing fruit
	consumePlayerRequest  : messageTypeVal++, //a player asks to consume another player
}

function createNewPlayerRequest(nickname) {
  return JSON.stringify(
    { 
      type: messageType.newPlayerRequest,
      data: {
        "nickname" : nickname
      }
    }
  )
}

function createPositionUpdateRequest(clientID, newX, newY, newDirection) {
  return JSON.stringify(
    {
      type: messageType.positionUpdateRequest,
      data: {
        id: clientID,
        location: {
          centre: {
            x: newX,
            y: newY
          },
          direction: newDirection
        },
      }
    }
  )
}

function createFruitConsumptionRequest(clientID, fruitID) {
  return JSON.stringify(
    {
      type: messageType.consumeFruitRequest,
      data: {
        client_id: clientID,
        fruit_id: fruitID
      }
    }
  )
}

function createPlayerConsumptionRequest(consumerID, consumedID) {
  return JSON.stringify(
    {
      type: messageType.consumePlayerRequest,
      data: {
        consumer_id: consumerID,
        consumed_id: consumedID
      }
    }
  )
}