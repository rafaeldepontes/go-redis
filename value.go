package main

import (
	"strconv"
)

const CRLF = "\r\n"

type Value struct {
	typ   string
	str   string
	bulk  string
	array []Value
}

func (v Value) Marshal() []byte {
	switch v.typ {
	case "array":
		return v.marshalArray()
	case "bulk":
		return v.marshalBulk()
	case "string":
		return v.marshalString()
	case "null":
		return v.marshalNull()
	case "error":
		return v.marshalError()
	default:
		return []byte{}
	}
}

func (v Value) marshalArray() []byte {
	size := len(v.array)
	var bytes []byte

	bytes = append(bytes, ARRAY)
	bytes = append(bytes, strconv.Itoa(size)...)
	bytes = append(bytes, CRLF...)

	for i := range size {
		bytes = append(bytes, v.array[i].Marshal()...)
	}
	return bytes
}

func (v Value) marshalBulk() []byte {
	var bytes []byte
	bytes = append(bytes, BULK)
	bytes = append(bytes, strconv.Itoa(len(v.bulk))...)
	bytes = append(bytes, CRLF...)
	bytes = append(bytes, v.bulk...)
	bytes = append(bytes, CRLF...)
	return bytes
}

func (v Value) marshalString() []byte {
	var bytes []byte
	bytes = append(bytes, STRING)
	bytes = append(bytes, v.str...)
	bytes = append(bytes, CRLF...)
	return bytes
}

func (v Value) marshalNull() []byte {
	var bytes []byte
	bytes = append(bytes, ERROR)
	bytes = append(bytes, v.str...)
	bytes = append(bytes, CRLF...)
	return bytes
}

func (v Value) marshalError() []byte {
	return []byte("$-1" + CRLF)
}
