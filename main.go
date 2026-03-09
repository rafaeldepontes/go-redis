package main

import (
	"errors"
	"io"
	"log"
	"net"
)

func main() {
	log.Println("[INFO] Server running on port: 6379")

	l, err := net.Listen("tcp", "6379")
	if err != nil {
		log.Fatalln("[ERROR] Cannot acess port 6379:", err)
	}

	conn, err := l.Accept()
	if err != nil {
		log.Println("[ERROR] Could not accept connection:", err)
		return
	}
	defer conn.Close()

	for {
		resp := NewResp(conn)
		value, err := resp.Read()
		if err != nil {
			if errors.Is(io.EOF, err) {
				break
			}
			log.Fatalln("[ERROR] something went wrong reading the client:", err)
		}

		_ = value

		writer := NewWriter(conn)
		writer.Write(Value{typ: "string", str: "OK"})
	}
}
