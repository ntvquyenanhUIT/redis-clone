package main

import (
	"fmt"
	"io"
)

type RespWriter struct {
	writer io.Writer
}

func NewRespWriter(w io.Writer) *RespWriter {
	return &RespWriter{
		writer: w,
	}
}

func (w *RespWriter) Write(v Value) error {

	bytes := v.Marshal()

	_, err := w.writer.Write(bytes)
	if err != nil {
		return fmt.Errorf("Error writing response to connection: %v", err)
	}

	return nil
}
