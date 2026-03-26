package main

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	pwd := "123456789"

	result, err := bcrypt.GenerateFromPassword([]byte(pwd), 4)
	if err != nil {
		fmt.Println("error", err)
	}

	fmt.Println("result: ", string(result))
	err = bcrypt.CompareHashAndPassword(result, []byte(pwd))
	if err != nil {
		fmt.Println("not matched")
	} else {
		fmt.Println("matched")
	}
}
