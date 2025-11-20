package main

import (
	"bufio"
	"fmt"
	"net"
)

func getHeaders(conn *net.Conn) ([]string, error) {
	var headers []string
	reader := bufio.NewReader(*conn)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			return nil, fmt.Errorf("error Reading Header Line: %v", err.Error())
		}

		if line == "\r\n" {
			break
		}
		headers = append(headers, line)
	}
	return headers, nil
}

func handleGet(path string, httpVersion string) string {
	if path == "/" {
		return fmt.Sprintf("%v 200 OK\r\n\r\n", httpVersion)
	}
	return fmt.Sprintf("%v 404 Not Found\r\n\r\n", httpVersion)
}
