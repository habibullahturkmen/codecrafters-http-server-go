package main

import (
	"fmt"
	"net"
	"os"
	"strings"
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
	headers, err := getHeaders(&conn)
	if err != nil {
		fmt.Println("Error Reading Header Line: ", err.Error())
		os.Exit(1)
	}
	requestLineHeader := strings.Fields(headers[0])
	method, path, httpVersion := requestLineHeader[0], requestLineHeader[1], requestLineHeader[2]
	// END: Handle Headers

	switch method {
	case "GET":
		response := handleGet(path, httpVersion)

		_, err = conn.Write([]byte(response))
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}
	case "POST":
	}

}
