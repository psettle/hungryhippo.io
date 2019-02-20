//open a websocket connection to the server that sent this document
const ws = new WebSocket("ws://" + window.location.hostname + "/ws")
ws.addEventListener('message', event => {
  onWebsocketReceive(JSON.parse(event.data))
});

var clientID = null

//when the document has loaded, attach a listener to the nickname input field
$(document).ready(function() {
  $("[name='nickname']").change(function() {
    //send a request to join the game to the server
    sendNewPlayerRequest($(this).val())
  })
})

//generates and sends a request to join the game to the server
function sendNewPlayerRequest(nickname) {
  var message = createNewPlayerRequest(nickname)
  ws.send(message)
}

//callback for messages from the server
function onWebsocketReceive(message) {
  //check for top level validity
  if(!('type' in message) || !('data' in message)) {
    console.log("Invalid message received")
    return
  }
  
  //decode message type
  switch(message.type) {
    case messageType.newPlayerResponse:
      onNewPlayerResponse(message.data)
      break;
    case messageType.positionUpdateMessage:
      onPositionUpdateMessage(message.data)
      break;
    case messageType.consumeFruitMessage:
      onConsumeFruitMessage(message.data)
      break;
    case messageType.newFruitMessage:
      onNewFruitMessage(message.data)
      break
    default:
      console.log("Unknown message received:" + message.type)
  }
}

function onNewPlayerResponse(message) {
  //validate message integrety
  if(!validateNewPlayerResponse(message)) {
    console.log("Invlaid NewPlayerResponse message received" + JSON.stringify(message))
    return
  }

  clientID = message.id

  //draw initial position, etc.
}

function onPositionUpdateMessage(message) {
  console.log("Position update message: " + JSON.stringify(message))
}

function onConsumeFruitMessage(message) {
  console.log("Fruit Consumed: " + JSON.stringify(message))
}

function onNewFruitMessage(message) {
  console.log("New Fruit: " + JSON.stringify(message))
}

//send a local position update event to the server
function sendPositionUpdateMessage(newX, newY, newDirection) {
  if(clientID == null)
  {
    console.log("Position update without id")
    return
  }

  message = createPositionUpdateMessage(clientID, newX, newY, newDirection)
  ws.send(message)
}

//send a request to consume a fruit to the server
function sendFruitConsumptionRequest(fruitID) {
  if(clientID == null)
  {
    console.log("Fruit consume without ID")
    return
  }

  message = createFruitConsumptionRequest(clientID, fruitID)
  ws.send(message)
}