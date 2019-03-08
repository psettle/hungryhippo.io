var BusinessLogic = (function() {
    var publicMethods = {} //no one calls business logic
    
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
        AppServer.subscribe(processGamestateUpdate)
        initNicknameTextbox()

        SpriteDrawing.Fruit.drawFruit(0.25, 0.25, 0.1)
    }

    var gamestate = {
        local : {
            player: null,           //the player associated with this client
            playerSprite: null,     //the sprite of the player assocaited with this client
            scale: 0.2,             //the scale of the local player
        }
    }

    function processGamestateUpdate(players, fruits) {
        //find outself in the gamestate
        gamestate.local.player = findLocalPlayer(players)

        //check if we've joined the game
        if(gamestate.local.player == null) {
            return;
        }

        //draw ourselves, if we haven't yet, it is easy because
        //- we are always the same size
        //- we are always in the middle of the screen
        //- we always face the cursor
        if(gamestate.local.playerSprite == null) {
            //draw self in middle of window
            gamestate.local.playerSprite = SpriteDrawing.Player.drawPlayer(0.5, 0.5, gamestate.local.scale)
            //register to receive direction updates
            Movement.subscribe(onDirectionChanged)

            //apply the first movement event so the hippo starts the right way
            var d = Movement.getDirection()
            onDirectionChanged(d.dx, d.dy)
        }    
    }

     
    function onDirectionChanged(dx, dy) {
        const speedScale = 10.0

        dx *= speedScale
        dy *= speedScale

        SpriteDrawing.Player.setLocalSpeed(gamestate.local.playerSprite, dx, dy)
    }

    function findLocalPlayer(players) {
        clientID = AppServer.getClientID()

        for(var i = 0; i < players.length; ++i) {
            var player = players[i]

            if(player.id == clientID) {
                return player
            }
        }

        return null
    }

    function initNicknameTextbox() {
        //grab relevant elements
        var input = $('.search-form');
        var search = $('input')
        var button = $('button');

        input.on('keyup', function (e) {
            //treat enter as a click on the button
            if( e.keyCode == 13) {
                button.trigger('click')
            }
        })
        
        button.on('click', function(e) {
            nickname = search.val()

            if (nickname == "") {
                //no name in the field
                return
            }

            search.val("")

            //send the new player request
            AppServer.sendNewPlayerRequest(nickname)
            //hide the nickname box
            input.removeClass('active');
        })
        search.on('focus', function() {
            input.addClass('focus');
        })

        search.on('blur', function() {
            search.val().length != 0 ? 
                input.addClass('focus') :
                input.removeClass('focus');
        })

        //trigger an 'open' animation
        input.addClass('active');
        search.focus();
    }

    return publicMethods
})();