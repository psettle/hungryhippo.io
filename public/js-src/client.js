//open a websocket connection to the server that sent this document
const ws = new WebSocket("ws://" + window.location.hostname + "/ws")
ws.addEventListener('message', event => {
  onWebsocketReceive(JSON.parse(event.data))
});

let clientID = null;

//when the document has loaded, attach a listener to the nickname input field
$(document).ready(function() {
  // TODO: replace with a button and bind the method that would let the player to join the game
  $("[name='nickname']").change(function() {
    //send a request to join the game to the server
    sendNewPlayerRequest($(this).val())
  })
});

//callback for messages from the server
function onWebsocketReceive(message) {
  //check for top level validity
  if(!('type' in message) || !('data' in message)) {
    console.log("Invalid message received");
    return
  }
  //decode message type
  switch(message.type) {
    case messageType.newPlayerResponse:
      onNewPlayerResponse(message.data);
      break;
    case messageType.gamestateUpdateMessage:
      onGamestateUpdateMessage(message.data);
      break;
    default:
      console.log("Unknown message received:" + message.type)
  }
}

function onNewPlayerResponse(message) {
  clientID = message.id;
  processNewIdUpdate(clientID)
}

function onGamestateUpdateMessage(message) {
  processGamestateUpdate(message.players.elements, message.players.count, message.fruits.elements, message.fruits.count)
}

//generates and sends a request to join the game to the server
function sendNewPlayerRequest(nickname) {
  let message = createNewPlayerRequest(nickname);
  ws.send(message)
}

//send a local position update event to the server
function sendPositionUpdateRequest(newX, newY, newDirection) {
  if(clientID == null)
  {
    console.log("Position update without id");
    return
  }

  let message = createPositionUpdateRequest(clientID, newX, newY, newDirection);
  ws.send(message)
}

//send a request to consume a fruit to the server
function sendFruitConsumptionRequest(fruitID) {
  if(clientID == null)
  {
    console.log("Fruit consume without ID");
    return
  }

  let message = createFruitConsumptionRequest(clientID, fruitID);
  ws.send(message)
}

//send a request to consume a player to the server
function sendPlayerConsumptionRequest(consumedID) {
  if(clientID == null)
  {
    console.log("Player consume without ID");
    return
  }

  let message = createPlayerConsumptionRequest(clientID, consumedID);
  ws.send(message)
}