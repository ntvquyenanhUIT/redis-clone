package main

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
)

// *2\r\n$5\r\nhello\r\n$5\r\nworld\r\n.

type Value struct {
	typ   string // array, bulk string
	str   string
	array []Value
}

type RespParser struct {
	reader *bufio.Reader
}

// why bufio?

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
