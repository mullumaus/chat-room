package protocol

import (
	"bufio"
	"fmt"
	"io"
)

//CommandReader is to read message
type CommandReader struct {
	reader *bufio.Reader
}

//NewCommandReader is to create a CommandReader object
func NewCommandReader(reader io.Reader) *CommandReader {
	return &CommandReader{reader: bufio.NewReader(reader)}
}

func (r *CommandReader) Read() (interface{}, error) {
	commandName, err := r.reader.ReadString(' ')

	if err != nil {
		return nil, err
	}

	switch commandName {
	case "MESSAGE":
		user, err := r.reader.ReadString(' ')
		if err != nil {
			return nil, err
		}
		message, err := r.reader.ReadString('\n')

		if err != nil {
			return nil, err
		}
		return MessageCommand{
			user[:len(user)-1],
			message[:len(message)-1],
		}, nil
	default:
		fmt.Printf("Unknow command :%v", commandName)
	}
	return nil, UnknownCommand
}

//ReadAll is to read all message
func (r *CommandReader) ReadAll() ([]interface{}, error) {
	commands := []interface{}{}

	for {
		command, err := r.Read()

		if command != nil {
			commands = append(commands, command)
		}
		if err == io.EOF {
			break
		} else if err != nil {
			return commands, err
		}
	}
	return commands, nil
}
