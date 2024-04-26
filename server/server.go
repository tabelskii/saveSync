package server

import (
	"encoding/binary"
	"fmt"
	"net"
	"os"
)

type Connection struct {
	server *Server
	conn   net.Conn
}

func (connection Connection) Close() {
	connection.conn.Close()
}

func (connection Connection) Read(data []byte) {
	connection.conn.Read(data)
}

func (connection Connection) Write(data []byte) {
	connection.conn.Write(data)
}

func (connection Connection) WriteTrue() {
	connection.Write([]byte{1})
}

func (connection Connection) WriteFalse() {
	connection.Write([]byte{0})
}

func (connection Connection) ReadByte() byte {
	data := make([]byte, 1)
	connection.Read(data)
	return data[0]
}

func (connection Connection) ReadInt64() int64 {
	intBytes := make([]byte, 8)
	connection.Read(intBytes)
	fmt.Println(intBytes)
	value, _ := binary.Varint(intBytes)
	return value
}

func (connection Connection) ReadPartition() []byte {
	length := connection.ReadInt64()
	dataBytes := make([]byte, length)
	connection.Read(dataBytes)
	return dataBytes
}

func (connection Connection) ReceiveFile(savePath string) {
	size := connection.ReadInt64()
	fmt.Println("File size: ", size)
	chunks := int(size / 2)
	file, err := os.OpenFile(savePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	buffer := make([]byte, 2)
	for i := 0; i < chunks; i++ {
		connection.Read(buffer)
		fmt.Println("Received: ", string(buffer))
		b, err := file.Write(buffer)
		if err != nil {
			fmt.Println(err)
			fmt.Println(b)
		}
	}
	file.Close()

}

type Server struct {
	instance *net.Listener
}

func (server *Server) Run(address string) {
	instance, err := net.Listen("tcp", address)
	if err != nil {
		panic(err)
	}
	server.instance = &instance

}

func (server Server) Close() {
	instance := *server.instance
	instance.Close()
}

func (server *Server) Accept() (Connection, error) {
	instance := *server.instance
	conn, err := instance.Accept()
	connection := Connection{server, conn}
	return connection, err
}
