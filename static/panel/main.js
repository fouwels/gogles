var ticks = 0

function render() {
    ticks++
    document.getElementById("ticker").innerHTML = "Healthkeeper v0.1 | " + ticks
}

setInterval(render, 100)