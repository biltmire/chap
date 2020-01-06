// Create WebSocket connection.
hostname = window.location.hostname
console.log('ws://'+hostname+':3000/ws')
const socket = new WebSocket('ws://'+hostname+':3000/ws');

// Send redirect url to person in queue when opponent enters game
function addInitListener(status,redirect){
  socket.addEventListener('open', function (event) {
    socket.send(JSON.stringify({flag: status,url: redirect}))
    socket.close()
  });
}
