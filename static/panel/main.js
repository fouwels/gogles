var ticks = 0

function render() {
    ticks++

    $.getJSON("/api", function(data) {
        $("#api_response").text(JSON.stringify(data, null, 2));
    })

    document.getElementById("ticker").innerHTML = "Healthkeeper v0.1 | " + ticks
}

setInterval(render, 100)