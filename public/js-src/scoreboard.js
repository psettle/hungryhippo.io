// Scoreboard
class Scoreboard {
    // Updates the scoreboard
    update(players) {
        if (players === null) return;
        // Find the current player
        let clientID = AppServer.getClientID();
        let myIndex = players.findIndex(function(player) {
            return player.id === clientID;
        });
        // Find the top ten players
        players = this.sortByKey(players, 'points').slice(0, 10);
        // Clear the table
        $("#scoreboard").find("tr:gt(0)").remove();
        players.forEach(function(player, index) {
            if (player.id === clientID) {
                $("#scoreboard").append(
                    '<tbody>' +
                        '<tr>' +
                            `<td><strong>${index + 1}.</strong></td>` +
                            `<td><strong>${player.nickname}</strong></td>` +
                            `<td><strong>${player.points}</strong></td>`+
                        '</tr>' +
                    '</tbody>'
                );
            } else {
                $("#scoreboard").append(
                    '<tbody>' +
                        '<tr>' +
                            `<td>${index + 1}.</td>` +
                            `<td>${player.nickname}</td>` +
                            `<td>${player.points}</td>`+
                        '</tr>' +
	                '</tbody>'
                );
            }
        });
        if (myIndex > 9) {
            $("#scoreboard").append(
            	'<tbody>' +
                    '<tr>' +
                        '<td></td>' +
                    '</tr>' +
	                '<tr>' +
                        `<td><strong>${myIndex + 1}</strong></td>` +
                        `<td><strong>${players[myIndex].nickname}</strong></td>` +
                        `<td><strong>${players[myIndex].points}</strong></td>`+
                    '</tr>' +
	            '</tbody>'
            );
        }
    }
    // Renders the scoreboard without data
    renderView() {
        let table = $(
	        '<div class="d-flex">' +
	            '<div class="p-2">' +
	                '<table class="table table-striped" id="scoreboard">' +
	                    '<thead>' +
	                        '<tr>' +
	                            '<th>#</th>' +
	                            '<th>Nicknames</th>' +
	                            '<th>Points</th>' +
	                        '</tr>' +
	                    '</thead>' +
	                '</table>' +
	            '</div>' +
	        '</div>'
        );
        $("body").append(table);
    }
    // Sorts the array in descending order
    sortByKey(array, key) {
        return array.sort(function(a, b) {
            let x = a[key]; let y = b[key];
            return ((x > y) ? -1 : ((x < y) ? 1 : 0));
        });
    }
}