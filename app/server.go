package main

import (
	"bufio"
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
	fmt.Println("Connection Established!")
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

	// START: Handle Headers
	reader := bufio.NewReader(conn)
	method, path, httpVersion, err := getRequestLine(reader)
	if err != nil {
		fmt.Println("Error Reading The Request Line: ", err.Error())
		os.Exit(1)
	}

	headers, err := getHeaders(reader)
	if err != nil {
		fmt.Println("Error Reading The Headers: ", err.Error())
		os.Exit(1)
	}
	// END: Handle Headers

	switch method {
	case "GET":
		response := handleGet(path, httpVersion, headers)
		_, err = conn.Write([]byte(response))
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}
	case "POST":
	}

}
