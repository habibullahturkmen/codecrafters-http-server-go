package main

import (
	"fmt"
	"net"
	"os"
)

type Server struct {
	listener net.Listener
}

func (s *Server) Listen() {
	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}
	s.listener = l
}

func (s *Server) Accept() net.Conn {
	conn, err := s.listener.Accept()
	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
		os.Exit(1)
	}

	return conn
}

func (s *Server) Close() {
	err := s.listener.Close()
	if err != nil {
		fmt.Println("Failed while closing the connection: ", err.Error())
		os.Exit(1)
	}
}

func (s *Server) Start() {
	s.Listen()
	defer s.Close()

	conn := s.Accept()

	_, err := conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
		os.Exit(1)
	}

	fmt.Println("Connection Established!")
}
