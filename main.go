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

var passwords []PasswordEntry

const dataFile = "passwords.enc"

var encryptionKey = []byte("32-byte-long-key-1234567890abcdef123456")

func encrypt(data []byte) ([]byte, error) {
	block, err := aes.NewCipher(encryptionKey)
	if err != nil {
		return nil, err
	}
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

func savePasswords() error { // Changed from savePassword to savePasswords
	data, err := json.Marshal(passwords) // Save entire slice, not single password
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
		return nil
	}
	encrypted, err := ioutil.ReadFile(dataFile)
	if err != nil {
		return err
	}
	decrypted, err := decrypt(encrypted) // Added missing decryption step
	if err != nil {
		return err
	}
	return json.Unmarshal(decrypted, &passwords)
}

func main() {
	fmt.Println("🔒 Go Password Manager")

	if err := loadPasswords(); err != nil {
		fmt.Println("Warning: Could not load passwords:", err)
	}

	for {
		fmt.Println("\nMenu Password Manager")
		fmt.Println("1. Add password")
		fmt.Println("2. Get password")
		fmt.Println("3. Exit")
		fmt.Print("> ")

		var choice int
		_, err := fmt.Scanln(&choice)
		if err != nil {
			fmt.Println("Invalid input. Try again.")
			continue
		}

		switch choice {
		case 1:
			var service, username, password string
			fmt.Print("Enter service: ")
			fmt.Scanln(&service)
			fmt.Print("Enter username: ")
			fmt.Scanln(&username)
			fmt.Print("Enter password: ")
			fmt.Scanln(&password)

			passwords = append(passwords, PasswordEntry{
				Service:  service,
				Username: username,
				Password: password,
			})

			if err := savePasswords(); err != nil { // Save all passwords after adding
				fmt.Println("Error saving:", err)
			} else {
				fmt.Println("Password saved for", service)
			}

		case 2:
			var service string
			fmt.Print("Enter service: ")
			fmt.Scanln(&service)
			found := false
			for _, entry := range passwords {
				if entry.Service == service {
					fmt.Printf("\nService: %s\nUsername: %s\nPassword: %s\n",
						entry.Service, entry.Username, entry.Password)
					found = true
					break
				}
			}
			if !found {
				fmt.Println("No password saved for", service)
			}

		case 3:
			if err := savePasswords(); err != nil {
				fmt.Println("Error saving passwords:", err)
			}
			fmt.Println("Exiting...")
			return

		default:
			fmt.Println("Invalid choice. Try again.")
		}
	}
}
