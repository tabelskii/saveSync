package db

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var dB *gorm.DB

func GetDB() *gorm.DB {
	var err error
	if dB == nil {
		dB, err = gorm.Open(postgres.New(postgres.Config{
			DSN:                  "host=localhost user=sync password=1234 dbname=save_sync port=5432", // data source name, refer https://github.com/jackc/pgx
			PreferSimpleProtocol: true,                                                                // disables implicit prepared statement usage. By default pgx automatically uses the extended protocol
		}), &gorm.Config{})
		for _, model := range []interface{}{&User{}, &File{}, &FileHistory{},
			&Machine{}, &Folder{}} {
			dB.AutoMigrate(model)
		}
		fmt.Println(dB, err)
	}
	return dB
}

type User struct {
	ID       uint `gorm:"primarykey"`
	Name     string
	Password string
	Files    []File
	Machines []Machine
}

func (user *User) SetPassword(password string) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		fmt.Println(err)
	}
	user.Password = string(hashed)
}

func (user *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		fmt.Println(err)
		return false
	}
	return true
}

type File struct {
	ID      uint `gorm:"primarykey"`
	Hash    string
	Path    string
	UserID  int
	History []FileHistory
}

type FileHistory struct {
	ID     uint `gorm:"primarykey"`
	FileID int
	Hash   string
}

func (FileHistory) TableName() string {
	return "files_history"
}

type Machine struct {
	ID         uint `gorm:"primarykey"`
	UserID     int
	HardWareId string
}

type Folder struct {
	ID        uint `gorm:"primarykey"`
	Path      string
	MachineID int
	Files     []File `gorm:"many2many:folder_files;"`
}
