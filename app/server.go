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

func (s *Server) start(dirName string) {
	s.listen()
	defer s.close()

	for {
		conn := s.accept()
		go func(conn net.Conn) {
			fmt.Println("New connection from:", conn.RemoteAddr())
			var responseHeader, responseBody string
			reader := bufio.NewReader(conn)

			requestLine, headers, body, err := parseHTTPRequest(reader)
			if err != nil {
				fmt.Println("Error parsing the http request: ", err.Error())
				os.Exit(1)
			}

			req := Request{
				method:      requestLine[0],
				path:        requestLine[1],
				httpVersion: requestLine[2],
				headers:     headers,
				body:        body,
			}

			switch req.method {
			case "GET":
				responseHeader, responseBody, err = handleGET(req, dirName)
				if err != nil {
					fmt.Println(err.Error())
					os.Exit(1)
				}
			case "POST":
				responseHeader, err = handlePOST(req, dirName)
				if err != nil {
					fmt.Println(err.Error())
					os.Exit(1)
				}
			}

			_, err = conn.Write([]byte(responseHeader))
			if err != nil {
				fmt.Println(err.Error())
				os.Exit(1)
			}

			if len(responseBody) > 0 {
				fmt.Println("body len ", len(responseBody))
				_, err := conn.Write([]byte(responseHeader))
				if err != nil {
					fmt.Println(err.Error())
					os.Exit(1)
				}
			}
		}(conn)
	}
}
