package server

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"sync"

	"../protocol"
)

type client struct {
	conn   net.Conn
	name   string
	writer *protocol.CommandWriter
}

//TCPChatServer is TCP chat server
type TCPChatServer struct {
	listener net.Listener
	clients  []*client
	mutex    *sync.Mutex
}

var (
	UnknownClient = errors.New("Unknown client")
)

func newServer() *TCPChatServer {
	return &TCPChatServer{
		mutex: &sync.Mutex{},
	}
}

//Listen is to listen on an address
func (s *TCPChatServer) Listen(address string) error {
	l, err := net.Listen("tcp", address)
	if err == nil {
		s.listener = l
	}
	fmt.Printf("Listening on %v", address)
	return err
}

//Close is to close a connection
func (s *TCPChatServer) Close() {
	s.listener.Close()
}

//Start is
func (s *TCPChatServer) Start() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			log.Print(err)
		} else {
			client := s.accept(conn)
			go s.server(client)
		}
	}
}

//Broadcast is
func (s *TCPChatServer) Broadcast(command interface{}) error {
	for _, client := range s.clients {
		client.writer.Write(command)
	}
	return nil
}

//Send is
func (s *TCPChatServer) Send(name string, command interface{}) error {
	for _, client := range s.clients {
		if client.name == name {
			return client.writer.Write(command)
		}
	}
	return UnknownClient
}

func (s *TCPChatServer) accept(conn net.Conn) *client {
	log.Printf("Accepting connection from %v, total clients: %v", conn.RemoteAddr().String(), len(s.clients)+1)

	s.mutex.Lock()
	defer s.mutex.Unlock()

	client := &client{
		conn:   conn,
		writer: protocol.NewCommandWriter(conn),
	}

	s.clients = append(s.clients, client)
	return client
}

func (s *TCPChatServer) remove(client *client) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	for i, check := range s.clients {
		if check == client {
			s.clients = append(s.clients[:i], s.clients[i+1:]...)
		}
	}
	log.Printf("Closing connection frm %v", client.conn.RemoteAddr().String())
	client.conn.Close()
}

func (s *TCPChatServer) server(client *client) {
	cmdReader := protocol.NewCommandReader(client.conn)
	defer s.remove(client)

	for {
		cmd, err := cmdReader.Read()
		if err != nil && err != io.EOF {
			log.Printf("Read error: %v", err)
		}
		if cmd != nil {
			switch v := cmd.(type) {
			case protocol.SendCommand:
				go s.Broadcast(protocol.MessageCommand{
					Message: v.Message,
					Name:    client.name,
				})
			case protocol.NameCommand:
				client.name = v.Name
			}
		}
		if err == io.EOF {
			break
		}
	}
}
