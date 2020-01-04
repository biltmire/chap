// Create WebSocket connection.
const socket = new WebSocket('ws://localhost:3000/ws');

// Connection opened
function addInitListener(msg){
  socket.addEventListener('open', function (event) {
    socket.send(
      JSON.stringify({flags: msg})
    )});
}
