let app = new PIXI.Application({width: 500, height: 500});

//init for pixi.js, must be ran before api functions below are called
$(document).ready( function() {
    app.renderer.backgroundColor = 0x061639;
    app.renderer.view.style.position = "relative";
    app.renderer.view.style.display = "block";
    app.renderer.autoDensity = true;
    app.renderer.resize(window.innerWidth, window.innerHeight);
    document.body.appendChild(app.view);

    let fruit1 = drawFruit(app, app.screen.width / 2, app.screen.height / 2, 0.05);
    let fruit2 = drawFruit(app, 10, 10, 0.25);
});


function drawFruit(app, x, y, scale) {
    let watermelon = PIXI.Sprite.from("js/images/watermelon.png");
    watermelon.scale.set(scale, scale);
    watermelon.position.set(x, y);
    app.stage.addChild(watermelon);
    return watermelon;
}

function eraseFruit(app, fruit) {
    app.stage.removeChild(fruit);
}