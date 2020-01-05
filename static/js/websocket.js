// Create WebSocket connection.
const socket = new WebSocket('ws://localhost:3000/ws');



// Send redirect url to person in queue when opponent enters game
function addInitListener(status,redirect){
  socket.addEventListener('open', function (event) {
    socket.send(JSON.stringify({flag: status,url: redirect}))
    socket.close()
  });
}
