function drawFruit(app, x, y, scale)
{
    var watermelon = PIXI.Sprite.from("js/images/watermelon.png");
    watermelon.scale.set(scale, scale);
    watermelon.position.set(x, y);
    app.stage.addChild(watermelon);
    return watermelon;
}

function eraseFruit(app, fruit)
{
    app.stage.removeChild(fruit);
}