var FruitManager = (function() {
    var pub = {
        fruitsUpdated: function(fruits) {
            fruitsUpdated(fruits)
        }
    }

    var fruitRecords = {}

    function fruitRecord() {
        return {
            sprite: null,
            dbRecord: null
        }
    }

    function fruitsUpdated(fruits) {
        //If local player isn't loaded yet then don't do anything
        if(PlayerManager.getLocalScore() == null) {
            return
        }

        //only use current fruit status
        fixFruitSet(fruits)

        //draw all fruit at the correct scale relative to the local player
        setFruitScale()

        setFruitPosition()
    }

    function fixFruitSet(fruits) {
        //Set all fruits to consumed to start
        for (var id in fruitRecords) {
            fruitRecords[id].consumed = true
        }

        //foreach fruit provided, check if they exist in the fruitsRecord
        //if they don't create a new sprite for them
        //
        //set all fruits !consumed that are touched here        
        for(var i = 0; i < fruits.length; ++i) {
            var fruit = fruits[i]
            var id = fruit.id

            if(!(id in fruitRecords)) {
                createFruit(fruit)
            }

            fruitRecords[id].consumed = false
        }
 
        //foreach fruit in the records, check if they have been consumed
        //if they have, remove them
        for(var id in fruitRecords) {
            if(fruitRecords[id].consumed) {
                SpriteDrawing.Fruit.eraseFruit(fruitRecords[id].sprite)
                delete fruitRecords[id]
            }
        }
    }

    function setFruitScale() {
        var localPlayerSize = PlayerManager.getLocalScore()

        for(var id in fruitRecords) {
            var fruit = fruitRecords[id]

            var scale = PositionManager.fruitScale / (localPlayerSize + 1)
            SpriteDrawing.Sprite.setScale(fruit.sprite, scale)
        }
    }

    function setFruitPosition() 
    {
        var localPlayerPosition = PlayerManager.getLocalPosition()

        for(var id in fruitRecords) {
            var fruit = fruitRecords[id]

            var fruitPosition = JSON.parse(JSON.stringify(fruit.dbRecord.position))
            var mapPosition = PositionManager.gameToScreen(localPlayerPosition, fruitPosition)

            SpriteDrawing.Sprite.setPosition(fruit.sprite, mapPosition.x, mapPosition.y)
        }
    }

    function createFruit(fruit) {
        //create the fruit sprite (position & sclae do not matter, they will be adjusted later)
        var sprite = SpriteDrawing.Fruit.drawFruit(0.0, 0.0, 1.0)

        //attach data to the sprite
        sprite.record = fruitRecord()

        //create local reference to that data
        var record = sprite.record
    
        //populate local record
        record.sprite = sprite  
        record.dbRecord = fruit 

        fruitRecords[fruit.id] = record
    }

    return pub
})();