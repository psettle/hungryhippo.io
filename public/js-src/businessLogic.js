$(document).ready( function() {
    var fruit1 = drawFruit(app, app.screen.width / 2, app.screen.height / 2, 0.05);
    var fruit2 = drawFruit(app, 10, 10, 0.25);
});

function processNewIdUpdate(id) {
    console.log("NewId Update Received:", id)
}

function processGamestateUpdate(players, playerCount, fruits, fruitCount) {
    console.log("Gamestate Update Received:", players, playerCount, fruits, fruitCount)
}