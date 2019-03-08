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
            },
            setLocalSpeed: function(sprite, dx, dy) {
                setDirection(sprite, dx, dy)
                setBackgroundSpeed(-dx, -dy)
                setSpeed(sprite, dx, dy)
            }
        },
        Sprite: {
            setDirection: function(sprite, dx, dy) {
                setDirection(sprite, dx, dy)
            },
            setSpeed: function(sprite, dx, dy) {
                setSpeed(sprite, dx, dy)
            }
        }
    }

    var state = {
        sprites: [],
        background: null,
        movingSprites: [],
    }

    var app = new PIXI.Application({width: 1000, height: 1000});
    var loader = PIXI.loader;
    //setup textures we need, must finish before API functions below are called
    loader
        .add("img/watermelon.png")
        .add("img/hippo.png")
        .add("img/swampbig.png")
        .load(init.readyYet);

    //init for pixi.js, must be ran before api functions below are called
    $(document).ready(init.readyYet);

    function onReady() {
        app.renderer.view.style.position = "absolute";
        app.renderer.view.style.display = "block";
        app.renderer.autoDensity = true;
        app.renderer.resize(window.innerWidth, window.innerHeight);
        $("body").append(app.view);

        backgroundInit()

        app.ticker.add(function(delta) {
            speedUpdates(delta)
        })

        //finished initializing everything, tell listeners we are ready
        for(var i = 0; i <  init.readyCallbacks.length; ++i) {
            init.readyCallbacks[i]();
        }
    }

    function speedUpdates(delta) {
        var dxBackground = state.background.dx
        var dyBackground = state.background.dy

        for(var i = 0; i < state.movingSprites.length; ++i) {
            state.movingSprites[i].position.x -= dxBackground * delta
            state.movingSprites[i].position.y -= dyBackground * delta
        }

        for(var i = 0; i < state.sprites.length; ++i) {
            state.sprites[i].position.x += dxBackground * delta
            state.sprites[i].position.y += dyBackground * delta
        }

        state.background.tilePosition.x += dxBackground * delta
        state.background.tilePosition.y += dyBackground * delta
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

    function setSpeed(sprite, dx, dy) {
        sprite.dx = dx
        sprite.dy = dy

        if(state.movingSprites.indexOf(sprite) === -1) {
            state.movingSprites.push(sprite)
        }
    }

    function setBackgroundSpeed(dx, dy) {
        state.background.dx = dx
        state.background.dy = dy
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

        state.sprites.push(watermelon)

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

        state.sprites.push(hippo)

        return hippo;
    }

    function erasePlayer(player) {
        app.stage.removeChild(player);
    }

    function backgroundInit() {
        // create a texture from an image path
        var texture = loader.resources["img/swampbig.png"].texture

        var stage = new PIXI.Container();

        // create a tiling sprite
        state.background = new PIXI.extras.TilingSprite(texture, app.screen.width, app.screen.height);
        state.background.tileScale.set(2, 2)
        setBackgroundSpeed(0, 0)
        stage.addChild(state.background);

        app.stage.addChild(stage)
    }

    return publicMethods;
})();