//Singleton client manager component
var AppServer = (function() {
    var pub = {
        //subscribe for gamestate updates
        //(players, players.count, fruits, fruits.count)
        subscribe: function (cb) {
            subscribe(cb)
        },
        sendNewPlayerRequest: function (nickname) {
            sendNewPlayerRequest(nickname)
        },
        sendFruitConsumptionRequest: function (fruitID) {
            sendFruitConsumptionRequest(fruitID)
        },
        sendPositionUpdateRequest: function (newX, newY, newDirection) {
            sendPositionUpdateRequest(newX, newY, newDirection)
        },
        sendPlayerConsumptionRequest: function (consumedID) {
            sendPlayerConsumptionRequest(consumedID)
        },
        getClientID: function() {
            return clientID
        }
    }

    var clientID = null
    
    var ws = null

    var server = null

    $.get("/ws", function(data) {
        server = JSON.parse(data)
        ws = new WebSocket("ws://" + window.location.hostname + ":" + server.port + "/ws")
        ws.addEventListener('message', event => {
            onWebsocketReceive(JSON.parse(event.data))
        });

        ws.addEventListener('error', event => {
            console.log('websocket error')
            onWebsocketClose()
        });

        ws.addEventListener("close", event => {
            console.log('websocket closed')
            onWebsocketClose()
        });
    })

    //callback for messages from the server
    function onWebsocketReceive(message) {
        //check for top level validity
        if(!('type' in message) || !('data' in message)) {
        console.log("Invalid message received")
        return
        }
        
        //decode message type
        switch(message.type) {
        case MessageCreator.messageType.newPlayerResponse:
            onNewPlayerResponse(message.data)
            break;
        case MessageCreator.messageType.gamestateUpdateMessage:
            onGamestateUpdateMessage(message.data)
            break;
        default:
            console.log("Unknown message received:" + message.type)
        }
    }

    function onWebsocketClose() {
        //grab a new app server, while telling the load balancer that an app server died
        $.get("/ws?dead-app-server-ip=" + server.ip + "&dead-app-server-port=" + server.port, function(data) {
            server = JSON.parse(data)
            ws = new WebSocket("ws://" +  window.location.hostname + ":" + server.port + "/ws?rejoin-clientid=" + clientID)
            ws.addEventListener('message', event => {
                onWebsocketReceive(JSON.parse(event.data))
            });
    
            ws.addEventListener("close", event => {
                onWebsocketClose()
            })
        })
    }
    
    var subscribers = []
    function subscribe(cb) {
        subscribers.push(cb)
    }

    function onNewPlayerResponse(message) {
        clientID = message.id
    }
        
    function onGamestateUpdateMessage(message) {
        //send update to subscribers
        for(var i = 0; i < subscribers.length; ++i) {
            subscribers[i]
                (
                message.players.elements,
                message.fruits.elements
                );
        }
    }
        
    //generates and sends a request to join the game to the server
    function sendNewPlayerRequest(nickname) {
        var message = MessageCreator.createNewPlayerRequest(nickname)
        ws.send(message)
    }
        
    //send a local position update event to the server
    function sendPositionUpdateRequest(newX, newY, newDirection) {
        if(clientID == null) {
            console.log("Position update without id")
            return
        }

        var message = MessageCreator.createPositionUpdateRequest(clientID, newX, newY, newDirection)
        ws.send(message)
    }
        
    //send a request to consume a fruit to the server
    function sendFruitConsumptionRequest(fruitID) {
        if(clientID == null) {
            console.log("Fruit consume without ID")
            return
        }

        var message = MessageCreator.createFruitConsumptionRequest(clientID, fruitID)
        ws.send(message)
    }
        
    //send a request to consume a player to the server
    function sendPlayerConsumptionRequest(consumedID) {
        if(clientID == null) {
            console.log("Player consume without ID")
            return
        }

        var message = MessageCreator.createPlayerConsumptionRequest(clientID, consumedID)
        ws.send(message)
    }

    return pub;
})();