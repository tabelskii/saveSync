package main

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"github.com/sandipmavani/hardwareid"
	"net"
	"os"
	"path/filepath"
)

func sendInt64(conn net.Conn, value int64) {
	valueConverted := make([]byte, 8)
	binary.PutVarint(valueConverted, value)
	conn.Write(valueConverted)
}

func writePartiton(conn net.Conn, data []byte) {
	length := len(data)
	sendInt64(conn, int64(length))
	conn.Write(data)
}

func sendFile(conn net.Conn, path string, bufferSize int) {
	file, err := os.Open(path)

	if err != nil {
		fmt.Println(err)

	}
	defer file.Close()
	fileInfo, _ := file.Stat()
	fileSize := fileInfo.Size()

	conn.Write([]byte{0})
	fmt.Println("Write zero byte")
	writePartiton(conn, []byte(path))
	fmt.Println("Write path: ", path)
	sendInt64(conn, fileSize)
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

func main() {
	hardwareId, _ := hardwareid.ID()

	conn, _ := net.Dial("tcp", "127.0.0.1:9000")
	writePartiton(conn, []byte("kirill"))
	writePartiton(conn, []byte("1234"))
	writePartiton(conn, []byte(hardwareId))

	response := make([]byte, 1)
	conn.Read(response)
	isAuthorized := response[0] == 1

	if !isAuthorized {
		return
	}
	fname := "./client/data_to_send"
	fname, _ = filepath.Abs(fname)
	sendFile(conn, fname, 3)

}
