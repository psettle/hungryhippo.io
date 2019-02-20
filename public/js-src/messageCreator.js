function createNewPlayerRequest(nickname) {
  return JSON.stringify(
    { 
      type: "NewPlayerRequest",
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
      type: "PositionUpdateRequest",
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