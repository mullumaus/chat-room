package client

import (
	"log"
	"net"

	"../protocol"
)

//ChatClient is inteface for chat client
type ChatClient interface {
	Dial(address string) error
	Send(command interface{}) error
	SendMessage(message string) error
	SetName(name string) error
	Start()
	Close()
	Incoming() chan protocol.MessageCommand
}

//TCPChatClient is chat client
type TCPChatClient struct {
	conn      net.Conn
	cmdReader *protocol.CommandReader
	CmdWriter *protocol.CommandWriter
	name      string
	error     chan error
	incoming  chan protocol.MessageCommand
}

//NewClient is to create a new client
func NewClient() *TCPChatClient {
	return &TCPChatClient{
		incoming: make(chan protocol.MessageCommand),
	}
}

//Dial is to setup connection
func (c *TCPChatClient) Dial(address string) error {
	conn, err := net.Dial("tcp", address)
	if err == nil {
		c.conn = conn
	}
	c.cmdReader = protocol.NewCommandReader(conn)
	c.CmdWriter = protocol.NewCommandWriter(conn)
	return err
}

//Start is to start receiving data
func (c *TCPChatClient) Start() {
	for {
		cmd, err := c.cmdReader.Read()
		if err != nil {
			c.error <- err
			break
		}
		if cmd != nil {
			switch v := cmd.(type) {
			case protocol.MessageCommand:
				c.incoming <- v
			default:
				log.Printf("Unknown command %v", v)
			}
		}
	}
}

//Close is to close connection
func (c *TCPChatClient) Close() {
	c.conn.Close()
}

//Incoming is to return connnection incoming data
func (c *TCPChatClient) Incoming() chan protocol.MessageCommand {
	return c.incoming
}

func (c *TCPChatClient) Error() chan error {
	return c.error
}

//Send is to write message
func (c *TCPChatClient) Send(command interface{}) error {
	return c.CmdWriter.Write(command)
}

//SetName is to set client name
func (c *TCPChatClient) SetName(name string) error {
	return c.Send(protocol.NameCommand{Name: name})
}

//SendMessage is to send message
func (c *TCPChatClient) SendMessage(message string) error {
	return c.Send(protocol.SendCommand{
		Message: message,
	})
}
