$(document).ready( function() {
    let app = new PIXI.Application({width: 500, height: 500});
    app.renderer.backgroundColor = 0x061639;
    app.renderer.view.style.position = "absolute";
    app.renderer.view.style.display = "block";
    app.renderer.autoDensity = true;
    app.renderer.resize(window.innerWidth, window.innerHeight);
    document.body.appendChild(app.view);

    var fruit1 = drawFruit(app, app.screen.width / 2, app.screen.height / 2, 0.05);
    var fruit2 = drawFruit(app, 10, 10, 0.25);
});