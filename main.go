package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
)

func main() {
	log.Println("[INFO] Server running on port: 6379")

	l, err := net.Listen("tcp", ":6379")
	if err != nil {
		log.Fatalln("[ERROR] Cannot acess port '6379':", err)
	}

	aof, err := NewAof("database.aof")
	if err != nil {
		log.Fatalln("[ERROR] Could not read or create AOF:", err)
	}
	defer aof.Close()

	err = aof.Read(func(value Value) {
		command := strings.ToUpper(value.array[0].bulk)
		args := value.array[1:]

		handler, ok := Handlers[command]
		if !ok {
			fmt.Println("Invalid command: ", command)
			return
		}

		handler(args)
	})
	if err != nil {
		log.Fatalln("[ERROR] Something went wrong while reading the AOF:", err)
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

		if command == "SET" || command == "HSET" {
			if err := aof.Write(value); err != nil {
				fmt.Println("[ERROR] Could not write the content into AOF:", err)
				continue
			}
		}

		result := handler(args)
		_ = writer.Write(result)
	}
}
