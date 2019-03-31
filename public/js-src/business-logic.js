var BusinessLogic = (function() {
    let scoreboard = null;
    var pub = {} //no one calls business logic
    
    //This is basically a promises system, theres probably a cleaner way to do this
    var dependencyCount = 1;
    SpriteDrawing.ready(readyYet);

    function readyYet() {
        dependencyCount--;
        if(dependencyCount <= 0) {
            onReady()
        }
    }

    function onReady() {
        AppServer.subscribe(processGamestateUpdate)
        initNicknameTextbox()

        scoreboard = new Scoreboard();
        scoreboard.renderView();
    }

    var gamestate = {
        local : {
            player: null,           //the player associated with this client
            change: {               //local changes in position that haven't been sent to the server
                x: 0,
                y: 0
            },
            playerSprite: null,     //the sprite of the player assocaited with this client
            scale: 0.2,             //the scale of the local player
        },
        players : {

        }
    }

    function processGamestateUpdate(players, fruits) {
        //send an update based on changes since the last update
        updateRemotePosition()

        scoreboard.update(players);

        //manage player sprites using new info
        PlayerManager.playersUpdated(players)

        FruitManager.fruitsUpdated(fruits)

        sendConsumptionRequests()
    }

    function updateRemotePosition() {
        //send a request to update position based on how far we went since the last update
        var pos = PlayerManager.getLocalPosition()

        if(pos == null) {
            return
        }

        AppServer.sendPositionUpdateRequest(pos.x, pos.y, pos.dir)
    }

    function initNicknameTextbox() {
        //grab relevant elements
        let input = $('.search-form');
        let search = $('input');
        let button = $('button');

        input.on('keyup', function (e) {
            //treat enter as a click on the button
            if( e.keyCode === 13) {
                button.trigger('click')
            }
        });
        
        button.on('click', function(e) {
            let nickname = search.val();

            if (nickname === '') return; // no name in the field

            search.val('');

            //send the new player request
            AppServer.sendNewPlayerRequest(nickname);
            //hide the nickname box
            input.removeClass('active');
        });
        search.on('focus', function() {
            input.addClass('focus');
        });

        search.on('blur', function() {
            search.val().length !== 0 ? input.addClass('focus') : input.removeClass('focus');
        });

        //trigger an 'open' animation
        input.addClass('active');
        search.focus();
    }

    function sendConsumptionRequests() {
        var locallyConsumedFruits = FruitManager.getLocallyConsumedFruit()
        for (var i = 0; i < locallyConsumedFruits.length; i++) {
            AppServer.sendFruitConsumptionRequest(locallyConsumedFruits[i])
        }
        FruitManager.resetLocallyConsumedFruit()

        var locallyConsumedPlayers = PlayerManager.getLocallyConsumedPlayers()
        for (var j = 0; j < locallyConsumedPlayers.length; j++) {
            console.log(locallyConsumedPlayers[j])
            AppServer.sendPlayerConsumptionRequest(locallyConsumedPlayers[j])
        }
        PlayerManager.resetLocallyConsumedPlayers()
    }

    return pub
})();