// Create WebSocket connection.
hostname = window.location.hostname
const socket = new WebSocket('wss://'+hostname+'/ws');

// Send redirect url to person in queue when opponent enters game
function addInitListener(status,redirect){
  socket.addEventListener('open', function (event) {
    socket.send(JSON.stringify({flag: status,url: redirect}))
    socket.close()
  });
}
