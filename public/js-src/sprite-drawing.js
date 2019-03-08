var SpriteDrawing = (function() {
    init = {
        ready: false,
        readyCallbacks: [],
        dependencies: 2,
        readyYet: function() {
            init.dependencies--;
            if(init.dependencies <= 0) {
                onReady()
            }
        }
    }

    publicMethods = {
        ready: function(cb) {
            if(init.ready) {
                cb()
            } else {
                init.readyCallbacks.push(cb)
            }
        },
        Fruit : {
            drawFruit: function(x, y, scale) {
                return drawFruit(x, y, scale)
            },
            eraseFruit: function(fruit) {
                eraseFruit(fruit)
            }
        },
        Player : {
            drawPlayer: function(x, y, scale) {
                return drawPlayer(x, y, scale)
            },
            erasePlayer: function(player) {
                erasePlayer(player)
            }
        },
        Sprite: {
            setDirection: function(sprite, dx, dy) {
                setDirection(sprite, dx, dy)
            }
        }
    }

    var app = new PIXI.Application({width: 500, height: 500});
    var loader = PIXI.loader;
    //setup textures we need, must finish before API functions below are called
    loader
        .add("img/watermelon.png")
        .add("img/hippo.png")
        .load(init.readyYet);

    //init for pixi.js, must be ran before api functions below are called
    $(document).ready(init.readyYet);

    function onReady() {
        app.renderer.backgroundColor = 0x061639;
        app.renderer.view.style.position = "absolute";
        app.renderer.view.style.display = "block";
        app.renderer.autoDensity = true;
        app.renderer.resize(window.innerWidth, window.innerHeight);
        $("body").append(app.view);
        //finished initializing everything, tell listeners we are ready
        for(var i = 0; i <  init.readyCallbacks.length; ++i) {
            init.readyCallbacks[i]();
        }
    }

    function setDirection(sprite, dx, dy) {
        //figure out what direction (dx, dy) defines
        var a = Math.atan(dy / dx)
        if(dx < 0) {
            a += Math.PI
        }

        //PIXI treats up as the default direction, instead of right like atan
        a += Math.PI / 2

        //set the sprite to that rotation
        sprite.rotation = a
    }

    function drawFruit(x, y, scale) {
        x *= app.screen.width;
        y *= app.screen.height;

        var watermelon = new PIXI.Sprite(
            loader.resources["img/watermelon.png"].texture
            );
        watermelon.scale.set(scale, scale);
        watermelon.position.set(x, y);
        app.stage.addChild(watermelon);
        return watermelon;
    }

    function eraseFruit(fruit) {
        app.stage.removeChild(fruit);
    }

    function drawPlayer(x, y, scale) {
        x *= app.screen.width;
        y *= app.screen.height;

        var texture = loader.resources["img/hippo.png"].texture

        var hippo = new PIXI.Sprite(texture);

        hippo.pivot.x = texture.width / 2
        hippo.pivot.y = texture.height / 2

        hippo.scale.set(scale, scale);
        hippo.position.set(x, y);
        app.stage.addChild(hippo);
        return hippo;
    }

    function erasePlayer(player) {
        app.stage.removeChild(fruit);
    }

    return publicMethods;
})();