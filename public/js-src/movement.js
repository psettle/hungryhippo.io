var Movement = (function() {
    var pub = {
        //subscribe to receive up to date direction
        //(dx, dy)
        subscribe: function(cb) {
            subscribe(cb)
        },
        //fetch the current direction from the centre of the screen to the 
        //mouse
        getDirection: function() {
            return getDirection()
        }
    } 

    var dependencyCount = 2;
    SpriteDrawing.ready(readyYet);
    $(document).ready(readyYet);

    function readyYet() {
        dependencyCount--;
        if(dependencyCount <= 0) {
            onReady()
        }
    }

    function onReady() {
        $("[name='mainbody']").mousemove(function(e) {
            position.mousex = e.pageX
            position.mousey = e.pageY
            publish()
        })
    }

    position = {
        mousex: getWidth()/ 2,
        mousey: getHeight() / 2,
    }

    var subscribers = []

    function subscribe(cb) {
        subscribers.push(cb)
    }

    function publish() {
        var d = getDirection()

        for(var i = 0; i < subscribers.length; ++i) {
            subscribers[i](d.dx, d.dy)
        }
    }

    function getDirection() {
        var w = getWidth() / 2
        var h = getHeight() / 2

        var dx = (position.mousex - w) / w;
        var dy = (position.mousey - h) / h;

        return {
            dx: dx,
            dy: dy
        }
    }

    function getWidth() {
        return $(window).width();
    }

    function getHeight() {
        return $(window).height();
    }

    return pub
})();