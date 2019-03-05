var BusinessLogic = (function() {
publicMethods = {} //no one calls business logic
    
//This is basically a promises system, theres probably a cleaner way
//to do this
var dependencyCount = 1;
SpriteDrawing.ready(readyYet);

function readyYet() {
    dependencyCount--;
    if(dependencyCount <= 0) {
        onReady()
    }
}

function onReady() {
    Movement.subscribe(onDirectionChanged)
    AppServer.subscribe(processGamestateUpdate)

    var fruit1 = SpriteDrawing.Fruit.drawFruit(0.5, 0.5, 0.05);
    var fruit2 = SpriteDrawing.Fruit.drawFruit(0.1, 0.1, 0.25);

    $("[name='nickname']").change(function() {
        //send a request to join the game to the server
        AppServer.sendNewPlayerRequest($(this).val())
    })
}

function onDirectionChanged(dx, dy) {

}

function processGamestateUpdate(players, playerCount, fruits, fruitCount) {
    //TODO: draw current gamestate from update
}

return publicMethods
})();