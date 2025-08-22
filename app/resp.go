package main

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
)

const (
	STRING   = '+'
	ERROR    = '-'
	INTEGERS = ':'
	BULK     = '$'
	ARRAY    = '*'
)

// *2\r\n$5\r\nhello\r\n$5\r\nworld\r\n.

type Value struct {
	typ   string // "array", "bulk", "string", "error", "integer"
	str   string // used for bulk or simple strings and error
	num   int    // used for integers
	array []Value
}

type RespParser struct {
	reader *bufio.Reader
}

// why bufio?
// simply put, we can't directly deal with raw data. We need something that provides tools to play with raw things.
// however, why not "io" but "bufio". Simply put, more convinient

func NewRespParser(rd io.Reader) *RespParser {
	return &RespParser{reader: bufio.NewReader(rd)}
}

func (p *RespParser) ReadValue() (Value, error) {
	typ, err := p.reader.ReadByte()
	if err != nil {
		return Value{}, err
	}

	switch typ {
	case '*':
		return p.readArray()
	case '$':
		return p.readBulkString()
	default:
		return Value{}, fmt.Errorf("unknown t ype: %v", string(typ))
	}
}

func (p *RespParser) readArray() (Value, error) {

	len, err := p.readInteger()
	if err != nil {
		return Value{}, fmt.Errorf("error reading array length :%v", err)
	}

	v := Value{
		typ:   "array",
		array: make([]Value, 0),
	}

	// *2\r\n$5\r\nhello\r\n$5\r\nworld\r\n.
	for i := 0; i < int(len); i++ {
		val, err := p.ReadValue()
		if err != nil {
			return v, err
		}

		v.array = append(v.array, val)
	}
	return v, nil
}

func (p *RespParser) readBulkString() (Value, error) {
	len, err := p.readInteger()
	if err != nil {
		return Value{}, fmt.Errorf("error reading bulk string length :%v", err)
	}

	data := make([]byte, len+2)

	_, err = io.ReadFull(p.reader, data)
	if err != nil {
		return Value{}, fmt.Errorf("error reading bulk string data: %v", err)
	}

	val := Value{
		typ: "bulk",
		str: string(data[:len]),
	}

	return val, nil
}

func (p *RespParser) readInteger() (int, error) {
	line, err := p.reader.ReadBytes('\n')

	if err != nil {
		return 0, err
	}

	len, err := strconv.Atoi(string(line[:len(line)-2]))
	if err != nil {
		return 0, err
	}
	return len, nil
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
	// *2\r\n$5\r\nhello\r\n$5\r\nworld\r\n.
	var bytes []byte
	len := len(v.array)
	bytes = append(bytes, ARRAY)
	bytes = append(bytes, strconv.Itoa(len)...)
	bytes = append(bytes, '\r', '\n')

	for i := 0; i < len; i++ {
		bytes = append(bytes, v.array[i].Marshal()...)
	}
	return bytes
}

func (v Value) marshalError() []byte {
	var bytes []byte
	bytes = append(bytes, ERROR)
	bytes = append(bytes, v.str...)
	bytes = append(bytes, '\r', '\n')

	return bytes
}

func (v Value) marshalNull() []byte {
	return []byte("$-1\r\n")
}

func (v Value) marshalString() []byte {
	var bytes []byte
	bytes = append(bytes, STRING)
	bytes = append(bytes, v.str...)
	bytes = append(bytes, '\r', '\n')

	return bytes
}

func (v Value) marshalBulk() []byte {
	var bytes []byte
	bytes = append(bytes, BULK)
	bytes = append(bytes, strconv.Itoa(len(v.str))...)
	bytes = append(bytes, '\r', '\n')
	bytes = append(bytes, v.str...)
	bytes = append(bytes, '\r', '\n')

	return bytes
}
