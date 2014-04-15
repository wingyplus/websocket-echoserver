var sock;
var roomName = document.getElementById("roomNameInput").value;
function submitMessage() {
    var message = document.getElementById("messageInput").value;

    console.log("Websocket - status: " + sock.readyState);
    if (sock.readyState == 1) {
        sock.send(JSON.stringify({ event: "ECHO", roomName: roomName, message: message }));
    }
}
function connect() {
    roomName = document.getElementById("roomNameInput").value;
    if (sock.readyState == 1) {
        sock.send(JSON.stringify({event: "ADD", roomName: roomName}));
    }
}
window.onload = function() {
    sock = new WebSocket("ws://localhost:12345/echo");
    sock.onopen = function(m) {
        console.log("CONNECTION opened..." + sock.readyState);
    }
    sock.onmessage = function(m) {
        console.log("MESSAGE RECEIVE: " + m.data)
    }
    sock.onclose  = function(m) {
        sock.send(JSON.stringify({event: "CLOSE", message: "CONNECTION CLOSE"}));
        console.log("CONNECTION CLOSE");
    }
};
