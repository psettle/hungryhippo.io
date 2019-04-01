//component for managing other player positions
var PlayerManager = (function() {
    var pub = {
        playersUpdated: function(players) {
            playersUpdated(players)
        },
        //fetch local game position
        //{x, y, dir}
        //null if local player doesn't exist
        getLocalPosition: function() {
            return getLocalPosition()
        },

        getLocalScore: function() {
            return getLocalScore()
        },

        getLocalSprite: function() {
            return getLocalSprite()
        },

        getLocallyConsumedPlayers: function() {
            return getLocallyConsumedPlayers()
        },

        resetLocallyConsumedPlayers: function() {
            resetLocallyConsumedPlayers()
        }
    }

    function playerRecord() {
        return {
            sprite: null,
            isLocal: false,
            dbRecord: null,
            //change in position since last update (for local player)
            change: {
                x: 0,
                y: 0
            },
            //position on last update (for other players)
            prevPos: {
                x: 0,
                y: 0
            }
        }
    }

    var local = null
    var playerRecords = {}
    var locallyConsumedPlayers = []

    function playersUpdated(players) {
        //first, we need an accurate local copy for doing relative rendering
        updateLocal(players)
        //stop if local doesn't exist yet
        if(local == null) {
            return
        }

        //delete dead players, add new players
        fixPlayerSet(players)

        //draw all players at the correct scale relative to the local player
        setPlayerScale()

        //set all players to face and move the correct direction
        setPlayerDirection()

        //set all players to the correct positions relative to the local player
        setPlayerPosition()

        checkForLocallyConsumedPlayers()
    }

    function updateLocal(players) {
        clientID = AppServer.getClientID()

        for(var i = 0; i < players.length; ++i) {
            var player = players[i]

            if(player.id == clientID) {
                if(local == null) {
                    createLocalPlayer(player)
                }

                //reset change since we are updating
                local.dbRecord = player
                var dbPos = local.dbRecord.location.centre
                dbPos.x += local.change.x
                dbPos.y += local.change.y
                dbPos = PositionManager.trimGamePosition(dbPos)
                local.change = {
                    x: 0,
                    y: 0
                }
                return
            }
        }
    }

    function fixPlayerSet(players) {
        //assume all players are dead
        for(var id in playerRecords) {
            playerRecords[id].dead = true
        }

        //foreach player provided, check if they exist in the set of players
        //if they do, update their db record
        //if the don't create a new sprite for them
        //
        //set all players !dead that are touched here
        for(var i = 0; i < players.length; ++i) {
            var player = players[i]
            var id = player.id
            
            if(id in playerRecords) {
                //calculate the 'previous position', i.e. where we expect that player to be given their
                //most recent trajectory
                playerRecords[id].prevPos = playerRecords[id].dbRecord.location.centre
                playerRecords[id].prevPos.x += playerRecords[id].change.x
                playerRecords[id].prevPos.y += playerRecords[id].change.y

                playerRecords[id].prevPos = PositionManager.trimGamePosition(playerRecords[id].prevPos)

                //update their 'true' position in their database record
                playerRecords[id].dbRecord = player

                //compute the difference between where we expected them to be and where they actually are
                playerRecords[id].change = {
                    x: playerRecords[id].prevPos.x - playerRecords[id].dbRecord.location.centre.x,
                    y: playerRecords[id].prevPos.y - playerRecords[id].dbRecord.location.centre.y
                }

                playerRecords[id].change = PositionManager.trimDirectionVector(playerRecords[id].change)
            } else {
                createPlayer(player)
            }

            playerRecords[id].dead = false
        }

        //foreach player in the records, check if they are dead
        //if they died, remove them
        for(var id in playerRecords) {
            if(playerRecords[id].dead) {
                SpriteDrawing.Player.erasePlayer(playerRecords[id].sprite)
                delete playerRecords[id]
            }
        }
    }

    function setPlayerScale() {
        var localSize = local.dbRecord.points + 1

        for(var id in playerRecords) {
            var player = playerRecords[id]

            if(player.isLocal) {
                continue
            }

            var relativeScale = (player.dbRecord.points + 1) / localSize
            var absoluteScale = relativeScale * PositionManager.playerScale

            SpriteDrawing.Sprite.setScale(player.sprite, absoluteScale)
        }
    }

    function setPlayerPosition() {
        var localPos = getLocalPosition()

        for(var id in playerRecords) {
            var player = playerRecords[id]

            if(player.isLocal) {
                continue
            }

            //set them to where they are (except now accounting)
            var currentPos = JSON.parse(JSON.stringify(player.prevPos))
            // currentPos.x -= player.change.x
            // currentPos.y -= player.change.y

            var mapPos = PositionManager.gameToScreen(localPos, currentPos)

            SpriteDrawing.Sprite.setPosition(player.sprite, mapPos.x, mapPos.y)
        }
    }

    function setPlayerDirection() {
        for(var id in playerRecords) {
            var player = playerRecords[id]

            if(player.isLocal) {
                continue
            }

            //we need to decompose the direction of the player into dx, dy such that abs([dx, dy]) == 1
            var dx = Math.cos(player.dbRecord.location.direction)
            var dy = Math.sin(player.dbRecord.location.direction)

            //apply the player speed metric
            dx *= PositionManager.speed
            dy *= PositionManager.speed
            

            //apply some speed to move them from where we thought they would be to where they are
            dExpectationX = -player.change.x
            dExpectationY = -player.change.y

            //pixi runs at 60 fps equivalent, with dx, dy in pixels/frame
            //we expect these updates every 250 ms, therefore to do a smooth transition:
            dExpectationX = (dExpectationX / (60 * 0.25))
            dExpectationY = (dExpectationY / (60 * 0.25))

            dx += dExpectationX
            dy += dExpectationY            

            //set direction and speed
            SpriteDrawing.Sprite.setDirection(player.sprite, dx, dy)
            SpriteDrawing.Sprite.setSpeed(player.sprite, dx, dy)
        }
    }

    function getLocalPosition() {
        if(local == null) {
            return null
        }

        //update position based on how far we went since the last update
        var position = {
            x: local.dbRecord.location.centre.x + local.change.x,
            y: local.dbRecord.location.centre.y + local.change.y
        }

        var dir = Movement.getDirection()
        var angle = PositionManager.toAngle(dir.dx, dir.dy)
        position = PositionManager.trimGamePosition(position)

        return {
            x: position.x,
            y: position.y,
            dir: angle
        }
    }

    function getLocalScore() {
        if(local == null) {
            return null
        }

        return local.dbRecord.points
    }

    function getLocalSprite() {
        if(local == null) {
            return null
        }

        return local.sprite
    }

    function getLocallyConsumedPlayers() {
        return locallyConsumedPlayers
    }

    function resetLocallyConsumedPlayers() {
        locallyConsumedPlayers = []
    }

    function createLocalPlayer(player) {
        //draw self in middle of window
        var sprite = SpriteDrawing.Player.drawPlayer(0.5, 0.5, PositionManager.playerScale)

        //attach data to the sprite
        sprite.record = playerRecord()

        //create local reference to that data
        local = sprite.record
        local.sprite = sprite
        
        local.isLocal = true
        local.dbRecord = player
        playerRecords[clientID] = local

        //setup gamepos handling
        SpriteDrawing.Sprite.setGamePositionHandler(sprite, onGamePositionUpdate)

        //register to receive direction updates
        Movement.subscribe(onDirectionChanged)

        //apply the first movement event so the hippo starts the right way
        var d = Movement.getDirection()
        onDirectionChanged(d.dx, d.dy)
    }

    function createPlayer(player) {
        //create the player sprite (position & sclae do not matter, they will be adjusted later)
        var sprite = SpriteDrawing.Player.drawPlayer(0.0, 0.0, 1.0)

        //attach data to the sprite
        sprite.record = playerRecord()

        //create local reference to that data
        var record = sprite.record
    
        //populate local record
        record.sprite = sprite     
        record.isLocal = false
        record.dbRecord = player
        record.prevPos = record.dbRecord.location.centre

        playerRecords[player.id] = record

        //setup gamepos handling
        SpriteDrawing.Sprite.setGamePositionHandler(sprite, onGamePositionUpdate)
    }

    //handler for game position updates to all player sprites
    function onGamePositionUpdate(sprite, dx, dy) {
        var record = sprite.record

        //dx, dy are in terms of window size
        //1 window size would be gamestate.mapSize / gamestate.mapScale game units
        dx *= PositionManager.mapSize / PositionManager.mapScale
        dy *= PositionManager.mapSize / PositionManager.mapScale

        record.change.x += dx
        record.change.y += dy
    }

    //handler for direction change events from the user
    function onDirectionChanged(dx, dy) {
        //normalize
        speed = Math.sqrt(dx * dx + dy * dy)
        if(speed == 0) {
            dx = 0
            dy = 0
        } else {
            dx /= speed
            dy /= speed
        }

        //apply local speed
        dx *= PositionManager.speed
        dy *= PositionManager.speed

        SpriteDrawing.Player.setLocalSpeed(local.sprite, dx, dy)
    }

    function checkForLocallyConsumedPlayers() {
        for(var id in playerRecords) {
            var player = playerRecords[id]
            if(player.isLocal) {
                continue
            }
            var playerScore = player.dbRecord.points
            var playerSprite = player.sprite
            var playerConsumed = SpriteDrawing.Collision.checkForCollision(local.sprite, playerSprite)

            if (playerConsumed && playerScore <= local.dbRecord.points) {
                locallyConsumedPlayers.push(id)
            }
        }
    }

    return pub
})();