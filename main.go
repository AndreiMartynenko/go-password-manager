package main

import "fmt"

type PasswordEntry struct {
	Service  string `json:"service"`
	Username string `json:"username"`
	Password string `json:"password"`
}

// in memory slice to hold entries
var passwords []PasswordEntry

func main() {
	// 1. Print a menu: "1. Add password", "2. Get password", "3. Exit"
	// 2. Use fmt.Scanln() to read user input
	// 3. Exit when user chooses "3"
	fmt.Println("🔒 Go Password Manager")

	for {
		//Print Menu
		fmt.Println("Menu Password Manager")
		fmt.Println("\n1. Add password")
		fmt.Println("2. Get password")
		fmt.Println("3. Exit")
		fmt.Println(">")

		//Read user input
		var choice int
		_, err := fmt.Scanln(&choice)
		if err != nil {
			fmt.Println("Invalid input. Try again.")
			continue
		}

		//Handle choices
		switch choice {
		case 1:
			fmt.Println("Add password selected")
		case 2:
			fmt.Println("Get password selected")
		case 3:
			fmt.Println("Exiting...")
			return
		default:
			fmt.Println("Invalid input. Try again.")
		}
	}
}
