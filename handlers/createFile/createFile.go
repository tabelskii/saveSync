package createFile

import (
	"crypto/rand"
	"fmt"
	"os"
	"path"
	"saveSync/server"
)
import "saveSync/users"

type Handler struct {
}

func (handler *Handler) Handle(connection server.Connection, profile *users.UserProfile) {
	originalPath := string(connection.ReadPartition())
	fmt.Println(originalPath)
	fullPath := profile.GetFolder()
	os.MkdirAll(fullPath, 0760)
	fnameB := make([]byte, 16)
	rand.Read(fnameB)
	fname := fmt.Sprintf("%X-%X-%X-%X-%X", fnameB[:4], fnameB[4:6], fnameB[6:8], fnameB[8:10], fnameB[10:])
	fullPath = path.Join(fullPath, fname)
	connection.ReceiveFile(fullPath)
}
