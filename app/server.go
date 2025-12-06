package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

type Server struct {
	Address  string
	listener net.Listener
}

type Request struct {
	method      string
	path        string
	httpVersion string
	headers     map[string]string
	body        []byte
}

func (s *Server) listen() {
	l, err := net.Listen("tcp", s.Address)
	if err != nil {
		fmt.Println(fmt.Errorf("failed to bind to %s: %w", s.Address, err))
		os.Exit(1)
	}
	s.listener = l
	fmt.Println("Server listening on", s.Address)
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
			defer conn.Close()
			fmt.Println("New connection from:", conn.RemoteAddr())

			reader := bufio.NewReader(conn)
			for {
				requestLine, headers, body, err := parseHTTPRequest(reader)
				if err != nil {
					fmt.Println("Error parsing the http request: ", err)
					return
				}

				if len(requestLine) < 3 {
					fmt.Println("Malformed request line")
					return
				}

				req := Request{
					method:      requestLine[0],
					path:        requestLine[1],
					httpVersion: requestLine[2],
					headers:     headers,
					body:        body,
				}

				var (
					responseHeader string
					responseBody   []byte
				)

				switch req.method {
				case "GET":
					responseHeader, responseBody, err = handleGET(req, dirName)
					if err != nil {
						fmt.Println("GET handler error:", err)
						return
					}
				case "POST":
					responseHeader, err = handlePOST(req, dirName)
					if err != nil {
						fmt.Println("POST handler error:", err)
						return
					}
				default:
					fmt.Println("Unsupported method:", req.method)
					return
				}

				// Write headers
				_, err = conn.Write([]byte(responseHeader))
				if err != nil {
					fmt.Println("Error writing header:", err)
					return
				}

				// Write body if present
				if len(responseBody) > 0 {
					_, err := conn.Write(responseBody)
					if err != nil {
						fmt.Println("Error writing body:", err)
						return
					}
				}
			}
		}(conn)
	}
}
