$(document).ready( function() {
    let app = new PIXI.Application({width: 500, height: 500});
    app.renderer.backgroundColor = 0x061639;
    app.renderer.view.style.position = "absolute";
    app.renderer.view.style.display = "block";
    app.renderer.autoDensity = true;
    app.renderer.resize(window.innerWidth, window.innerHeight);
    document.body.appendChild(app.view);

    var fruit1 = drawFruit(app.screen.width / 2, app.screen.height / 2);
    var fruit2 = drawFruit(10, 10);

    function drawFruit(x, y)
    {
        var watermelon = PIXI.Sprite.from("js/images/watermelon.png");
        watermelon.scale.set(0.05, 0.05);
        watermelon.position.set(x, y);
        app.stage.addChild(watermelon);
        return watermelon;
    }

    function eraseFruit(fruit)
    {
        app.stage.removeChild(fruit);
    }
});