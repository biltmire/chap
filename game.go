package main

import ("time"
        "math/rand"
        "fmt"
        "github.com/gorilla/websocket")

//Charset for generating random strings
const charset = "abcdefghijklmnopqrstuvwxyz" +
  "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

//Random seed for generating random strings
var seededRand *rand.Rand = rand.New(
  rand.NewSource(time.Now().UnixNano()))

//Does a lookup on the host and returns what color they are in the board if
//they are in a game
func hostLookup(game_id string,host string) string{
	return game_list[game_id].PlayerMap[host]
}

//Give two host names, create a new game object to store in game_list. Return
//the id of the new game
func gameManager(player_one string, player_two string) string{
	generator := rand.New(rand.NewSource(time.Now().UnixNano()))
	new_map := make(map[string]string)
	x := generator.Float64()
	if(x > 0.5) {
		new_map[player_one] = "white"
		new_map[player_two] = "black"
	} else {
		new_map[player_one] = "black"
		new_map[player_two] = "white"
	}
  id := String(9)
  empty_slice := [2]*websocket.Conn{nil,nil}
	new_game := Game{new_map,empty_slice,make(chan Message)}
  game_list[id] = &new_game
  fmt.Println(new_map)
  fmt.Println(len(game_list))
  return id
}

func StringWithCharset(length int, charset string) string {
  b := make([]byte, length)
  for i := range b {
    b[i] = charset[seededRand.Intn(len(charset))]
  }
  return string(b)
}

func String(length int) string {
  return StringWithCharset(length, charset)
}
