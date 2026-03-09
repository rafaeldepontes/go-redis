package main

import (
	"errors"
	"io"
	"log"
	"net"
	"strings"
)

func main() {
	log.Println("[INFO] Server running on port: 6379")

	l, err := net.Listen("tcp", ":6379")
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
			if errors.Is(err, io.EOF) {
				break
			}
			log.Fatalln("[ERROR] Something went wrong reading the client:", err)
		}

		if value.typ != "array" {
			log.Println("[WARN] Invalid request, expected array")
			continue
		}

		if len(value.array) == 0 {
			log.Println("[WARN] Invalid request, expected array size > 0")
			continue
		}

		command := strings.ToUpper(value.array[0].bulk)
		args := value.array[1:]

		writer := NewWriter(conn)

		handler, ok := Handlers[command]
		if !ok {
			log.Println("[WARN] Invalid command: ", command)
			_ = writer.Write(Value{typ: "string", str: ""})
			continue
		}

		result := handler(args)
		_ = writer.Write(result)
	}
}
