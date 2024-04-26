package createFile

import (
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"path"
	"saveSync/db"
	"saveSync/server"
)
import "saveSync/users"

type Handler struct {
}

func (handler *Handler) Handle(connection server.Connection, profile *users.UserProfile) {
	originalPath := string(connection.ReadPartition())
	//originalFolder := path.Dir(originalPath)
	//originalName := path.Base(originalPath)
	fmt.Println(originalPath)
	fullPath := profile.GetFolder()
	os.MkdirAll(fullPath, 0760)
	fnameB := make([]byte, 16)
	rand.Read(fnameB)
	fname := fmt.Sprintf("%X-%X-%X-%X-%X", fnameB[:4], fnameB[4:6], fnameB[6:8], fnameB[8:10], fnameB[10:])
	fullPath = path.Join(fullPath, fname)
	connection.ReceiveFile(fullPath)

	f, _ := os.Open(fullPath)
	defer f.Close()
	h := sha256.New()
	io.Copy(h, f)
	hashSum := string(h.Sum(nil))
	fmt.Println(hashSum)
	machine := db.Machine{HardWareId: profile.HardwareId}
	db.GetOrCreate(&machine)
	database := db.GetDB()
	database.Where(&machine).First(&machine)
	fmt.Println(&machine)
	if machine.ID == 0 {
		machine.UserID = int(profile.Id)
		database.Create(machine)
	}
}
