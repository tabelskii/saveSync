package basicHandler

import "saveSync/server"
import "saveSync/users"

type Handler interface {
	Handle(server.Connection, *users.UserProfile)
}
