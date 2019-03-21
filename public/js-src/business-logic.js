var BusinessLogic = (function() {
    var publicMethods = {}; //no one calls business logic
    let scoreboard = null;
    
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
        Movement.subscribe(onDirectionChanged);
        AppServer.subscribe(processGamestateUpdate);
        initNicknameTextbox()

        scoreboard = new Scoreboard();
        scoreboard.renderView();

        var fruit1 = SpriteDrawing.Fruit.drawFruit(0.5, 0.5, 0.05);
        var fruit2 = SpriteDrawing.Fruit.drawFruit(0.1, 0.1, 0.25);  
    }

    function onDirectionChanged(dx, dy) {
        //TODO: draw local direction from cursor
        //console.log(dx, dy)
    }

    function processGamestateUpdate(players, fruits) {
        //TODO: draw current gamestate from update
        console.log(players);
        scoreboard.update(players);
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

    return publicMethods;
})();