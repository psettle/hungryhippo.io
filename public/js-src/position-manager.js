//component for converting types of game coordinates
var PositionManager = (function() {
    var pub = {
        gameToScreen: function(localPlayerPos, targetPos) {
            return gameToScreen(localPlayerPos, targetPos)
        },
        trimGamePosition: function(position) {
            return trimGamePosition(position)
        },
        trimDirectionVector: function(vector) {
            return trimDirectionVector(vector)
        },
        toAngle(dx, dy) {
            return toAngle(dx, dy)
        },
        mapSize: 1000,              //size of the map in 'game units'
        mapScale: 1,               //how many screens it takes to cross the map
        speed: 5,
        playerScale: 0.1,
        fruitScale: 0.05
    }

    function gameToScreen(localPlayerPos, targetPos) {
        //compute difference from local player position


        var t = JSON.parse(JSON.stringify(targetPos))
        var l = JSON.parse(JSON.stringify(localPlayerPos))

        t = trimGamePosition(t)
        l = trimGamePosition(l)

        t.x += pub.mapSize
        t.y += pub.mapSize

        var d = {
            x: t.x - l.x,
            y: t.y - l.y
        }

        d = trimDirectionVector(d)
        
        //normalize to map coordinates
        d.x /= (pub.mapSize / ( 2 * pub.mapScale))
        d.y /= (pub.mapSize / ( 2 * pub.mapScale))

        d.x += 1
        d.y += 1

        d.x /= 2
        d.y /= 2

        return d
    }

    function trimGamePosition(position) {
        //apply wraparound effects such that x, y are bounded in [0, gamestate.mapSize]
        while(position.x < 0) {
            position.x += pub.mapSize
        }

        while(position.x > pub.mapSize) {
            position.x -= pub.mapSize
        }

        while(position.y < 0) {
            position.y += pub.mapSize
        }

        while(position.y > pub.mapSize) {
            position.y -= pub.mapSize
        }

        return position
    }

    function trimDirectionVector(vector) {
        while(vector.x < -pub.mapSize / 2) {
            vector.x += pub.mapSize
        }

        while(vector.x > pub.mapSize / 2) {
            vector.x -= pub.mapSize
        }

        while(vector.y < -pub.mapSize / 2) {
            vector.y += pub.mapSize
        }

        while(vector.y > pub.mapSize / 2) {
            vector.y -= pub.mapSize
        }

        return vector
    }

    function toAngle(dx, dy) {
        var angle = 0
        if(dx != 0) {
            angle = Math.atan(dy / dx)
            if(dx < 0) {
                angle += Math.PI
            }
        }

        if(angle < 0) {
            angle += 2 * Math.PI
        }

        return angle
    }
    return pub
})();