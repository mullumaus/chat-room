package protocol

import (
	"fmt"
	"io"
)

//CommandWriter is writer command
type CommandWriter struct {
	writer io.Writer
}

//NewCommandWriter is to create a CommandWrite object
func NewCommandWriter(writer io.Writer) *CommandWriter {
	return &CommandWriter{writer: writer}
}

func (w *CommandWriter) writeString(message string) error {
	_, err := w.writer.Write([]byte(message))
	return err
}

func (w *CommandWriter) Write(command interface{}) error {
	var err error

	switch v := command.(type) {
	case SendCommand:
		err = w.writeString(fmt.Sprintf("SEND %v\n", v.Message))
	case MessageCommand:
		err = w.writeString(fmt.Sprintf("MESSAGE %v %v", v.Name, v.Message))
	case NameCommand:
		err = w.writeString(fmt.Sprintf("NAME %v\n", v.Name))
	default:
		err = UnknownCommand
	}
	return err
}
