package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"sync"
	"time"
)

type Aof struct {
	file *os.File
	rd   *bufio.Reader
	mu   sync.RWMutex
}

func NewAof(path string) (*Aof, error) {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		return nil, err
	}

	aof := &Aof{
		file: f,
		rd:   bufio.NewReader(f),
		mu:   sync.RWMutex{},
	}

	go func() {
		for {
			aof.mu.Lock()
			err := aof.file.Sync()
			if err != nil {
				fmt.Println("[ERROR] Something really bad happened while calling 'aof.file.Sync':", err)
			}
			aof.mu.Unlock()
			time.Sleep(time.Second)
		}
	}()

	return aof, nil
}

func (aof *Aof) Close() error {
	defer aof.mu.Unlock()
	aof.mu.Lock()

	return aof.file.Close()
}

func (aof *Aof) Write(val Value) error {
	defer aof.mu.Unlock()
	aof.mu.Lock()

	_, err := aof.file.Write(val.Marshal())
	return err
}

func (aof *Aof) Read(callback func(value Value)) error {
	aof.mu.Lock()
	defer aof.mu.Unlock()

	resp := NewResp(aof.file)

	for {
		value, err := resp.Read()
		if err == io.EOF {
			break
		}

		if err != nil {
			return err
		}

		callback(value)
	}

	return nil
}
