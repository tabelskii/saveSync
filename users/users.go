package users

import (
	"fmt"
	"path"
	"path/filepath"
	"saveSync/config"
)

type UserProfile struct {
	Name       string
	Authorized bool
	HardwareId string
	Id         uint
}

func (user *UserProfile) GetFolder() string {
	fullPath, err := filepath.Abs(path.Join(config.BasicFolder, user.Name))
	if err != nil {
		fmt.Println(err)
	}
	return fullPath
}
