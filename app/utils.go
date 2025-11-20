package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
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

	if strings.HasPrefix(path, "/echo") {
		content := strings.TrimPrefix(path, "/echo")

		// If path is "/echo/abc", remove the "/"
		if strings.HasPrefix(content, "/") {
			content = content[1:]
		}
		return fmt.Sprintf("%s 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", httpVersion, len(content), content)
	}

	return fmt.Sprintf("%v 404 Not Found\r\n\r\n", httpVersion)
}
