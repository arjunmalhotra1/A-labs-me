// Sample program to show that you cannot take the address
// of an element in a map.
package main

// player represents someone playing our game.
type player struct {
	name  string
	score int
}

func main() {

	// Declare a map with initial values using a map literal cosntruction .
	players := map[string]player{
		"anna":  {"Anna", 42},
		"jacob": {"Jacob", 21},
	}

	// Trying to take the address of a map element fails.
	// We are trying to read the value of the map and takes it's address.
	// But because this value is not going to be places in memory directly.
	// There is no variable going to be associated with it.
	// Compiler says - "Cannot take the address of players["anna"]""
	// Unless it is stored some where.
	// Lines 32-34 are the correct way to do this.
	anna := &players["anna"]
	anna.score++

	// ./example4.go:23:10: cannot take the address of players["anna"]

	// Instead take the element, modify it, and put it back.
	player := players["anna"]
	player.score++
	players["anna"] = player
}
