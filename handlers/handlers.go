package handlers

import (
	"fmt"
	"saveSync/db"
	"saveSync/handlers/basicHandler"
	"saveSync/handlers/createFile"
	"saveSync/server"
	"saveSync/users"
)

func auth(connection server.Connection) users.UserProfile {
	userName := string(connection.ReadPartition())
	password := string(connection.ReadPartition())
	user := db.User{}
	database := db.GetDB()
	database.Where("name = ?", userName).First(&user)
	authorized := user.CheckPassword(password)
	hardwareId := string(connection.ReadPartition())

	if authorized {
		machine := db.Machine{HardWareId: hardwareId, UserID: int(user.ID)}
		db.GetOrCreate(&machine)
	}

	fmt.Println(user, password, hardwareId, authorized)
	return users.UserProfile{userName, authorized, hardwareId, user.ID}

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
