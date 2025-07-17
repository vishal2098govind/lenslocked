package main

import (
	"fmt"
	"os"

	"golang.org/x/crypto/bcrypt"
)

func main() {

	switch os.Args[1] {
	case "hash":
		// hash password and return
		if len(os.Args) < 3 {
			fmt.Println("Invalid args for hash.")
			return
		}
		hash(os.Args[2])
	case "compare":
		// compare password with hash
		if len(os.Args) < 4 {
			fmt.Println("Invalid args for compare.")
			return
		}
		compare(os.Args[2], os.Args[3])
	default:
		fmt.Println("Invalid command")
	}
}

func hash(password string) {
	hbytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Printf("error hashing: %v\n", password)
		return
	}
	fmt.Println(string(hbytes))
}

func compare(password, hashValue string) {
	err := bcrypt.CompareHashAndPassword([]byte(hashValue), []byte(password))
	if err != nil {
		fmt.Printf("Incorrect password %q\n", password)
		return
	}
	fmt.Println("Password is correct")
}
