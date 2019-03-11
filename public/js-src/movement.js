window.onload = function () {
    let currentX = 0;
    let currentY = 0;

    document.addEventListener("keydown", function onEvent(event) {
        if (event.key === "ArrowLeft") {
            currentX--;
            sendPositionUpdateMessage(currentX, currentY, 0)
        }
        if (event.key === "ArrowRight") {
            currentX++;
            sendPositionUpdateMessage(currentX, currentY, 0)
        }
        if (event.key === "ArrowDown") {
            currentY--;
            sendPositionUpdateMessage(currentX, currentY, 0)
        }
        if (event.key === "ArrowUp") {
            currentY++;
            sendPositionUpdateMessage(currentX, currentY, 0)
        }
    });
}
