package main

import (
	"flag"
	"fmt"
	"github.com/sandipmavani/hardwareid"
	"net"
	"path/filepath"
	"saveSync/client"
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

func clientSendFile() {
	hardwareId, _ := hardwareid.ID()

	conn, _ := net.Dial("tcp", "127.0.0.1:8080")
	client.WritePartiton(conn, []byte("kirill"))
	client.WritePartiton(conn, []byte("1234"))
	client.WritePartiton(conn, []byte(hardwareId))

	response := make([]byte, 1)
	conn.Read(response)
	isAuthorized := response[0] == 1

	if !isAuthorized {
		return
	}
	fname := "./client/data_to_send"
	fname, _ = filepath.Abs(fname)
	client.SendFile(conn, fname, 3)
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
	fmt.Println(action)
	switch action {
	case "run":
		run(address)
	case "create-user":
		newUser(userName, password)
	case "send-file":
		clientSendFile()
	}

}
