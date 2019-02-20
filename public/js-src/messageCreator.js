var messageTypeVal = 0
var messageType = {
  newPlayerRequest      : messageTypeVal++,
	newPlayerResponse     : messageTypeVal++,
  positionUpdateRequest : messageTypeVal++,
  positionUpdateMessage : messageTypeVal++
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