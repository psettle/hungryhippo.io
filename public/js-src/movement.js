
window.onload = function () {
    var currentX = 0;
    var currentY = 0;

    document.addEventListener("keydown", function onEvent(event) {
        if (event.key === "ArrowLeft") {
            currentX--;
            updatePositionString();
        }
        if (event.key === "ArrowRight") {
            currentX++;
            updatePositionString();
        }
        if (event.key === "ArrowDown") {
            currentY--;
            updatePositionString();
        }
        if (event.key === "ArrowUp") {
            currentY++;
            updatePositionString();
        }
    });

    function updatePositionString() {
        document.getElementById("positionIndicator").innerHTML = 
            "You are at position (" + currentX + ", " + currentY + ")";
    }
}
