var messageTypeVal = 0
var messageType = {
  newPlayerRequest      : messageTypeVal++, //A new player requests to join the game 
  newPlayerResponse     : messageTypeVal++, //Server acknowledging new player request, provided initial condition
  
  positionUpdateRequest : messageTypeVal++, //A player asks to be moved to a new location
  positionUpdateMessage : messageTypeVal++, //Server tells player about the position of all players

  consumeFruitRequest  : messageTypeVal++, //A player asks to consume an existing fruit
	consumeFruitMessage  : messageTypeVal++, //server notifies clients that a fruit has died
	newFruitMessage      : messageTypeVal++, //the server has generated a new fruit

	consumePlayerRequest  : messageTypeVal++, //a player asks to consume another player
	consumePlayerResponse : messageTypeVal++, //the server accepts/denies the consumption request
	playerDeathMessage    : messageTypeVal++ //server notififies a player that they have died, and must submit a new newPlayerRequest
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

function validateNewPlayerResponse(message) {
  if(!('id' in message)) {
    return false
  }

  if(!('location' in message)) {
    return false
  }

  if(!('direction' in message.location)) {
    return false
  }

  if(!('centre' in message.location)) {
    return false
  }

  if(!('x' in message.location.centre)) {
    return false
  }

  if(!('y' in message.location.centre)) {
    return false
  }

  if(!('points' in message)) {
    return false
  }

  return true
}

function createPositionUpdateMessage(clientID, newX, newY, newDirection) {
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