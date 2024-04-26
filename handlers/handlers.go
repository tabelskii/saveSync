package handlers

import (
	"fmt"
	"saveSync/handlers/basicHandler"
	"saveSync/handlers/createFile"
	"saveSync/server"
	"saveSync/users"
)

func auth(connection server.Connection) users.UserProfile {
	user := string(connection.ReadPartition())
	pass := string(connection.ReadPartition())
	hardwareId := string(connection.ReadPartition())
	authorized := pass == "1234"
	fmt.Println(user, pass, hardwareId, authorized)
	return users.UserProfile{user, authorized, hardwareId}

}

func HandleNewConnection(connenction server.Connection) {
	handlers := [1]basicHandler.Handler{&createFile.Handler{}}

	defer connenction.Close()
	user := auth(connenction)
	if !user.Authorized {
		connenction.WriteFalse()
		return
	}
	connenction.WriteTrue()
	operation := connenction.ReadByte()

	handler := handlers[operation]
	handler.Handle(connenction, &user)
}
