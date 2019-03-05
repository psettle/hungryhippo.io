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
    }

    var app = new PIXI.Application({width: 500, height: 500});
    var loader = PIXI.loader;
    //setup textures we need, must finish before API functions below are called
    loader
        .add("img/watermelon.png")
        .load(init.readyYet);

    //init for pixi.js, must be ran before api functions below are called
    $(document).ready(init.readyYet);

    function onReady() {
        app.renderer.backgroundColor = 0x061639;
        app.renderer.view.style.position = "absolute";
        app.renderer.view.style.display = "block";
        app.renderer.autoDensity = true;
        app.renderer.resize(window.innerWidth, window.innerHeight);
        document.body.appendChild(app.view);
        //finished initializing everything, tell listeners we are ready
        for(var i = 0; i <  init.readyCallbacks.length; ++i) {
            init.readyCallbacks[i]();
        }
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

    return publicMethods;
})();