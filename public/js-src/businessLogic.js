let scoreboard = new Scoreboard();

// Invoked upon page refresh
$(window).on('beforeunload', function() {

});

$(document).ready( function() {
    scoreboard.renderView();
});

function processNewIdUpdate(id) {
    console.log("NewId Update Received:", id);
}

function processGamestateUpdate(players, playerCount, fruits, fruitCount) {
    console.log("Gamestate Update Received:", players, playerCount, fruits, fruitCount);
    scoreboard.update(players);
}