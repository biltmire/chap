
var timerId = 0;

// Send redirect url to person in queue when opponent enters game
function addInitListener(status,redirect){
  // Create WebSocket connection.
  socket = createQueueSocket()
  socket.addEventListener('open', function (event) {
    socket.send(JSON.stringify({flag: status,url: redirect}))
    socket.close()
  });
}

//Create and add the GameListener which makes moves in the game when message is
//received
function addGameListener(color){
  gamesock = createGameSocket(color)
  setTimeout(keepAlive(),10000)

  gamesock.addEventListener('message', function (event) {
    var msg = JSON.parse(event.data)
    if(msg['flags'] != "keepAlive"){
      var fen = msg['fen']
      delete msg['fen']
      game.move(msg)
      board.position(fen,true)
      updateStatus()
    }
    else {
      console.log('keep alive')
    }
  });

  gamesock.onclose = function(event) {
    cancelKeepAlive();
  }
}

//Keep websocket alive given a timer
function keepAlive() {
    var timeout = 30000;
    if (gamesock.readyState == gamesock.OPEN) {
        var message = { 'color': "", 'from': "", 'to': "", 'flags': "keepAlive", 'piece': "p", 'san': "e4", 'fen': "" }
        gamesock.send(JSON.stringify(message));
    }
    timerId = setTimeout(keepAlive, timeout);
}

//Kill websocket
function cancelKeepAlive() {
    if (timerId) {
        clearTimeout(timerId);
    }
}

//Create and add the message listener for the queue
function addQueueListener(){
  // Create WebSocket connection.
  socket = createQueueSocket()
  socket.addEventListener('message', function (event) {
    var msg = JSON.parse(event.data)
      if(msg.flag == 'start') {
        parser = new URL(window.location.href)
        player_name = parser.searchParams.get("player")
        //Redirect the person in queue to a game with current player param
        socket.close()
        location.href = msg.url + "&player="+player_name
      }
  });
}

//Create the gamesocket connection for the game
function createGameSocket(color){
  stem = window.location.href.split('/')[3]
  host = window.location.hostname
  id = stem.substring(8,17)
  const gamesock = new WebSocket('wss://'+host+'/gamesock?id='+id+'&color='+color)
  return gamesock
}

//Create the socket connection for the queue
function createQueueSocket(){
  hostname = window.location.hostname
  const socket = new WebSocket('wss://'+hostname+'/ws');
  return socket
}
