// Scoreboard
class Scoreboard {
    // Updates the scoreboard
    update(players) {
        if (players === null) return;
        // find the current player
        let clientID = AppServer.getClientID()
        let me = players.find(function(player) {
            return player.id === clientID;
        });
        // find the top ten players
        players = this.sortByKey(players, 'points').slice(0, 10);
        // clear the table
        $("#scoreboard").find("tr:gt(0)").remove();
        let inTopTen = false;
        players.forEach(function(player) {
            if (player.id === clientID) {
                inTopTen = true;
                $("#scoreboard").append(
                  '<tr>' +
                    `<td><strong>${player.nickname}</strong></td>` +
                    `<td><strong>${player.points}</strong></td>`+
                  '</tr>');
            } else {
                $("#scoreboard").append(
                  '<tr>' +
                    `<td>${player.nickname}</td>` +
                    `<td>${player.points}</td>`+
                  '</tr>');
            }
        });
        if (!inTopTen) {
            $("#scoreboard").append(
              '<tr>' +
                '<td></td>' +
              '</tr>' +
              '<tr>' +
                `<td><strong>${me.nickname}</strong></td>` +
                `<td><strong>${me.points}</strong></td>`+
              '</tr>');
        }
    }

    renderView() {
        let table = $("<table id='scoreboard' style ='position: absolute; top: 8px; right:16px; font-size: 18px; color: red;'></table>");
        let header = $("<tr></tr>");
        let nicknames = $("<th></th>").text("Nicknames");
        let points = $("<th></th>").text("Points");
        header.append(nicknames, points);
        table.append(header);
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