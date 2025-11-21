package main

import (
	"bufio"
	"fmt"
	"strings"
)

const (
	userAgent = "User-Agent"
)

func getRequestLine(reader *bufio.Reader) (string, string, string, error) {
	requestLineHeader, err := reader.ReadString('\n')
	if err != nil {
		return "", "", "", fmt.Errorf("error Reading Header Line: %v", err.Error())
	}

	requestLineParts := strings.Split(strings.TrimSpace(requestLineHeader), " ")
	// method, path, httpVersion, nil
	return requestLineParts[0], requestLineParts[1], requestLineParts[2], nil
}

func getHeaders(reader *bufio.Reader) (map[string]string, error) {
	headers := map[string]string{}
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			return map[string]string{}, fmt.Errorf("error Reading Header Line: %v", err.Error())
		}

		if line == "\r\n" {
			break
		}

		line = strings.TrimSpace(line)

		header := strings.SplitN(line, ":", 2)
		if len(header) == 2 {
			headers[strings.TrimSpace(header[0])] = strings.TrimSpace(header[1])
		}
	}

	return headers, nil
}

func handleGet(path string, httpVersion string, headers map[string]string) string {
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

	if strings.TrimRight(path, "/") == "/user-agent" {
		userAgent := headers[userAgent]
		return fmt.Sprintf("%s 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", httpVersion, len(userAgent), userAgent)
	}

	return fmt.Sprintf("%v 404 Not Found\r\n\r\n", httpVersion)
}
