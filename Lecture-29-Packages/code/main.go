package main

import (
	"fmt"

	"github.com/Suyash-Kamath/Golang-Coders-Gyan/code/auth"
	"github.com/Suyash-Kamath/Golang-Coders-Gyan/code/user"
	"github.com/fatih/color"
)



func main(){
	auth.LoginWithCredentials("Suyash", "123b349587b")
	session := auth.GetSession()
	fmt.Println("Session is ", session)

	user:= user.User{
		// Name: "Suyash",
		Email: "suyash@example.com",
	}
	// fmt.Println("User is ", user)

	color.Green(user.Email)



}
