package main

import (
	"flag"
	"fmt"
	"saveSync/db"
	"saveSync/handlers"
	"saveSync/server"
)

func run(address string) {
	app := server.Server{}

	app.Run(address)
	db.GetDB()

	for {
		connection, err := app.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}
		go handlers.HandleNewConnection(connection)

	}
}

func newUser(userName string, password string) {
	dataBase := db.GetDB()
	newUser := db.User{Name: userName}
	newUser.SetPassword(password)
	tx := dataBase.Create(&newUser)
	fmt.Println(tx)
}

func main() {
	var address string
	var action string
	var userName string
	var password string
	flag.StringVar(&address, "address", ":8080", "server address")
	flag.StringVar(&action, "action", "run", "action")
	flag.StringVar(&userName, "username", "", "username")
	flag.StringVar(&password, "password", "", "password")
	flag.Parse()

	switch action {
	case "run":
		run(address)
	case "create-user":
		newUser(userName, password)
	}

}
