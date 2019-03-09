//component for converting types of game coordinates
var PositionManager = (function() {
    var pub = {
        gameToScreen: function(localPlayerPos, targetPos) {
            return gameToScreen(localPlayerPos, targetPos)
        },
        trimGamePosition: function(position) {
            return trimGamePosition(position)
        },
        toAngle(dx, dy) {
            return toAngle(dx, dy)
        },
        mapSize: 1000,              //size of the map in 'game units'
        mapScale: 1,               //how many screens it takes to cross the map
        speed: 5,
        playerScale: 0.1,
    }

    function gameToScreen(localPlayerPos, targetPos) {
        //compute difference from local player position
        var dx = targetPos.x - localPlayerPos.x
        var dy = targetPos.y - localPlayerPos.y

        //normalize to map coordinates
        dx /= (pub.mapSize / ( 2 * pub.mapScale))
        dy /= (pub.mapSize / ( 2 * pub.mapScale))

        return {
            x: dx,
            y: dy
        }
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