package main

import ("time"
        "math/rand")

//Struct that contains information on the game
type Game struct {
	PlayerMap map[string]string
	GameID string
}

//Does a lookup on the host and returns what color they are in the board if
//they are in a game
func hostLookup(host string) string{
	for game := range game_list {
		for player := range game_list[game].PlayerMap {
			if(host == player){
					return game_list[game].PlayerMap[host]
			}
		}
	}
	return ""
}

//Give two host names, create a player map dictionary and store it in a new
//game object which is appended to game_list
func gameManager(player_one string, player_two string){
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
	new_game := Game{new_map, "abcdefg"}
	game_list = append([]Game{new_game},game_list...)
}
