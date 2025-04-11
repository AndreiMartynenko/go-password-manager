package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

type PasswordEntry struct {
	Service  string `json:"service"`
	Username string `json:"username"`
	Password string `json:"password"`
}

// in memory slice to hold entries
var passwords []PasswordEntry

const dataFile = "passwords.enc"

var encryptionKey = []byte("32-byte-long-key-1234567890abcdef123456") // Replace later!

// --- Encryption/Decryption Helpers ---

func encrypt(data []byte) ([]byte, error) {
	block, err := aes.NewCipher(encryptionKey)
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}
	return gcm.Seal(nonce, nonce, data, nil), nil
}

func decrypt(data []byte) ([]byte, error) {
	block, err := aes.NewCipher(encryptionKey)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonceSize := gcm.NonceSize()
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	return gcm.Open(nil, nonce, ciphertext, nil)
}

// --- Save/Load Functions ---
func savePassword(password PasswordEntry) error {
	data, err := json.Marshal(password)
	if err != nil {
		return err
	}

	encrypted, err := encrypt(data)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(dataFile, encrypted, 0644)
}

func loadPasswords() error {
	if _, err := os.Stat(dataFile); os.IsNotExist(err) {
		return nil // No file yet (first run)
	}
	encrypted, err := ioutil.ReadFile(dataFile)
	if err != nil {
		return err
	}
	return json.Unmarshal(encrypted, &passwords)
}

func main() {
	// 1. Print a menu: "1. Add password", "2. Get password", "3. Exit"
	// 2. Use fmt.Scanln() to read user input
	// 3. Exit when user chooses "3"
	fmt.Println("🔒 Go Password Manager")

	// Load existing passwords on startup
	if err := loadPasswords(); err != nil {
		fmt.Println("Warning: Could not load passwords.", err)
	}

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
			var service, username, password string
			fmt.Print("Enter service: ")
			fmt.Scanln(&service)
			fmt.Print("Enter username: ")
			fmt.Scanln(&username)
			fmt.Print("Enter password: ")
			fmt.Scanln(&password)

			if err := savePassword(PasswordEntry{service, username, password}); err != nil {
				fmt.Println("Error saving password:", err)
			}

			// Append to passwords slice
			passwords = append(passwords, PasswordEntry{
				Service:  service,
				Username: username,
				Password: password,
			})
			fmt.Println("Password saved for", service)
		case 2:
			var service string
			fmt.Print("Enter service: ")
			fmt.Scanln(&service)
			for _, entry := range passwords {
				if entry.Service == service {
					fmt.Println("Service: \nUsername: \nPassword: \n", entry.Service, entry.Username, entry.Password)
					break
				}
				fmt.Println("No password saved for", service)
		case 3:
			if err := savePasswords(); err != nil {
				fmt.Println("Error saving passwords:", err)
				}
				fmt.Println("Existing...")
				return
			}
		}
	}
}
