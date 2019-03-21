// Scoreboard
class Scoreboard {
    // Updates the scoreboard
    update(players) {
        if (players === null) return;
        // find the current player
        let clientID = AppServer.getClientID();
        console.log('client id: ' + clientID);
        let myIndex = players.findIndex(function(player) {
            return player.id === clientID;
        });
        // find the top ten players
        players = this.sortByKey(players, 'points').slice(0, 10);
        // clear the table
        $("#scoreboard").find("tr:gt(0)").remove();
        players.forEach(function(player, index) {
            if (player.id === clientID) {
                $("#scoreboard").append(
                    '<tr>' +
                        `<td><strong>${index + 1}.</strong></td>` +
                        `<td><strong>${player.nickname}</strong></td>` +
                        `<td><strong>${player.points}</strong></td>`+
                    '</tr>');
            } else {
                $("#scoreboard").append(
                    '<tr>' +
                        `<td>${index + 1}.</td>` +
                        `<td>${player.nickname}</td>` +
                        `<td>${player.points}</td>`+
                    '</tr>');
            }
        });
        if (myIndex > 9) {
            $("#scoreboard").append(
                '<tr>' +
                    '<td></td>' +
                '</tr>' +
                '<tr>' +
                    `<td><strong>${myIndex + 1}</strong></td>` +
                    `<td><strong>${players[myIndex].nickname}</strong></td>` +
                    `<td><strong>${players[myIndex].points}</strong></td>`+
                '</tr>');
        }
    }
    // Render the scoreboard without data
    renderView() {
        let table = $(
            `<div class="container">` +
                `<div class="row">` +
                    `<div class="col-8"></div>` +
                    `<div class="col-4">` +
                        `<table class="table table-striped" id="scoreboard">` +
                            `<thead>` +
                                `<tr>` +
                                    `<th>#</th>` +
                                    `<th>Nicknames</th>` +
                                    `<th>Points</th>` +
                                `</tr>` +
                            `</thead>` +
                        `</table>` +
                    `</div>` +
                `</div>` +
            `</div>`);
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