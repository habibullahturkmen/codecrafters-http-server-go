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

type Request struct {
	method      string
	path        string
	httpVersion string
	headers     map[string]string
	body        []uint8
}

func (s *Server) listen() {
	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}
	s.listener = l
	fmt.Println("Connection Established!")
}

func (s *Server) accept() net.Conn {
	conn, err := s.listener.Accept()
	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
		os.Exit(1)
	}

	return conn
}

func (s *Server) close() {
	err := s.listener.Close()
	if err != nil {
		fmt.Println("Failed while closing the connection: ", err.Error())
		os.Exit(1)
	}
}

func (s *Server) start(dir string) {
	s.listen()
	defer s.close()
	for {
		conn := s.accept()
		go func(conn net.Conn) {
			fmt.Println("New connection from:", conn.RemoteAddr())
			reader := bufio.NewReader(conn)

			// START: Handle Headers
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

			req := Request{
				method:      method,
				path:        path,
				httpVersion: httpVersion,
				headers:     headers,
				body:        nil,
			}
			// END: Handle Headers

			switch method {
			case "GET":
				response := handleGet(req, dir)
				_, err = conn.Write([]byte(response))
				if err != nil {
					fmt.Println("Error accepting connection: ", err.Error())
					os.Exit(1)
				}
			case "POST":
			}
		}(conn)
	}
}
