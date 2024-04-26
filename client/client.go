package client

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"net"
	"os"
)

func SendInt64(conn net.Conn, value int64) {
	valueConverted := make([]byte, 8)
	binary.PutVarint(valueConverted, value)
	conn.Write(valueConverted)
}

func WritePartiton(conn net.Conn, data []byte) {
	length := len(data)
	SendInt64(conn, int64(length))
	conn.Write(data)
}

func SendFile(conn net.Conn, path string, bufferSize int) {
	file, err := os.Open(path)

	if err != nil {
		fmt.Println(err)

	}
	defer file.Close()
	fileInfo, _ := file.Stat()
	fileSize := fileInfo.Size()

	conn.Write([]byte{0})
	fmt.Println("Write zero byte")
	WritePartiton(conn, []byte(path))
	fmt.Println("Write path: ", path)
	SendInt64(conn, fileSize)
	fmt.Println("Write file size:", fileSize)

	reader := bufio.NewReader(file)

	err = nil
	size := 2

	readed := 1
	for readed > 0 {
		buffer := make([]byte, size)
		readed, err = reader.Read(buffer)
		if err != nil {
			fmt.Println(err)
			break
		}
		fmt.Println("Write ", string(buffer))
		conn.Write(buffer)
	}

}
