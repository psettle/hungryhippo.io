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

    pub = {
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
            },
            setScale: function(sprite, scale) {
                setScale(sprite, scale)
            },
            setPosition: function(sprite, x, y) {
                setPosition(sprite, x, y)
            },
            //register a handler for position updates
            //(sprite, dx, dy), in window size units, independent of player perspective
            setGamePositionHandler: function(sprite, cb) {
                sprite.updateGamePos = cb
            }
        },
        Collision: {
            checkForCollision: function(sprite1, sprite2) {
                return checkForCollision(sprite1, sprite2)
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
        .add("img/water.png")
        .load(init.readyYet);

    //init for pixi.js, must be ran before api functions below are called
    $(document).ready(init.readyYet);

    var bump = new Bump(app.renderer);

    var backgroundGroup = new PIXI.display.Group(0, false);
    var fruitGroup = new PIXI.display.Group(1, false);
    var hippoGroup = new PIXI.display.Group(2, false);
    app.stage = new PIXI.display.Stage();
    //Don't reorder these because the zIndex value given to the contructor of the groups
    //doesn't actually do anything as of this version of pixi-display and relies on this ordering
    app.stage.addChild(new PIXI.display.Layer(backgroundGroup));
    app.stage.addChild(new PIXI.display.Layer(fruitGroup));
    app.stage.addChild(new PIXI.display.Layer(hippoGroup));

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

        //apply movement to moving objects
        for(var i = 0; i < state.movingSprites.length; ++i) {
            var sprite = state.movingSprites[i]

            translateSprite(sprite, sprite.dx * delta, sprite.dy * delta)
        }

        // apply map translation to all objects, so they appear static on the map while not moving
        for(var i = 0; i < state.sprites.length; ++i) {
            var sprite = state.sprites[i]

            sprite.position.x += dxBackground * delta
            sprite.position.y += dyBackground * delta
        }

        //move the background to create perception of movement
        state.background.tilePosition.x += dxBackground * delta
        state.background.tilePosition.y += dyBackground * delta
    }

    function setDirection(sprite, dx, dy) {
        //set the sprite to that direction
        sprite.rotation = PositionManager.toAngle(dx, dy)

        //add another 90 deg cause pixi uses weird angles
        sprite.rotation += Math.PI / 2
    }

    function setSpeed(sprite, dx, dy) {
        sprite.dx = dx
        sprite.dy = dy

        if(state.movingSprites.indexOf(sprite) === -1) {
            state.movingSprites.push(sprite)
        }
    }

    function setScale(sprite, scale) {
        sprite.scale.set(scale, scale)
    }

    function setPosition(sprite, x, y) {
        x *= app.screen.width;
        y *= app.screen.height;

        sprite.position.x = x
        sprite.position.y = y
    }

    function translateSprite(sprite, dx, dy) {
        sprite.position.x += dx
        sprite.position.y += dy

        if('updateGamePos' in sprite) {
            //callback needs dx, dy in terms of window size, they are currently in pixels
            dx /= app.screen.width
            dy /= app.screen.height

            //adjust game position if sprite cares
            sprite.updateGamePos(sprite, dx, dy)
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
        watermelon.parentGroup = fruitGroup;
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

        hippo.parentGroup = hippoGroup;
        app.stage.addChild(hippo);

        state.sprites.push(hippo)

        return hippo;
    }

    function erasePlayer(player) {
        app.stage.removeChild(player);
    }

    function backgroundInit() {
        // create a texture from an image path
        var texture = loader.resources["img/water.png"].texture

        var stage = new PIXI.Container();

        // create a tiling sprite
        state.background = new PIXI.extras.TilingSprite(texture, app.screen.width, app.screen.height);
        state.background.tileScale.set(2, 2)
        setBackgroundSpeed(0, 0);
        stage.addChild(state.background);

        state.background.parentGroup = backgroundGroup;
        app.stage.addChild(stage)
    }

    function checkForCollision(sprite1, sprite2) {
        return bump.hit(sprite1, sprite2)
    }

    return pub;
})();
