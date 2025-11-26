package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

const (
	userAgent = "User-Agent"
)

func parseHTTPRequest(reader *bufio.Reader) ([]string, map[string]string, []byte, error) {
	requestLine, err := reader.ReadString('\n')
	if err != nil {
		return nil, nil, nil, fmt.Errorf("error reading the request line: %v", err.Error())
	}

	headers := map[string]string{}
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			return nil, nil, nil, fmt.Errorf("error reading the header: %v", err.Error())
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

	var body []byte
	if contentLength, ok := headers["Content-Length"]; ok {
		length, err := strconv.Atoi(contentLength)
		body = make([]byte, length)
		_, err = io.ReadFull(reader, body)
		if err != nil {
			return nil, nil, nil, fmt.Errorf("error reading the body: %v", err.Error())
		}
	}

	return strings.Split(strings.TrimRight(requestLine, "\r\n"), " "), headers, body, nil
}

func handleGet(req Request, dirName string) string {
	if req.path == "/" {
		return fmt.Sprintf("%v 200 OK\r\n\r\n", req.httpVersion)
	}

	if strings.HasPrefix(req.path, "/echo") {
		content := strings.TrimPrefix(req.path, "/echo")

		if strings.HasPrefix(content, "/") {
			content = content[1:]
		}
		return fmt.Sprintf("%s 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", req.httpVersion, len(content), content)
	}

	if strings.TrimRight(req.path, "/") == "/user-agent" {
		userAgent := req.headers[userAgent]
		return fmt.Sprintf("%s 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", req.httpVersion, len(userAgent), userAgent)
	}

	if strings.HasPrefix(req.path, "/files") {
		fileName := strings.TrimPrefix(req.path, "/files")

		if strings.HasPrefix(fileName, "/") {
			fileName = fileName[1:]
		}

		file, err := os.ReadFile(fmt.Sprintf("%s/%s", dirName, fileName))
		if err != nil {
			if strings.Contains(err.Error(), "no such file or directory") {
				return fmt.Sprintf("%v 404 Not Found\r\n\r\n", req.httpVersion)
			}
			fmt.Println("Failed reading file: ", err.Error())
			os.Exit(1)
		}

		return fmt.Sprintf("%s 200 OK\r\nContent-Type: application/octet-stream\r\nContent-Length: %d\r\n\r\n%s", req.httpVersion, len(string(file)), string(file))
	}

	return fmt.Sprintf("%v 404 Not Found\r\n\r\n", req.httpVersion)
}

func getDirName(args []string, flag string) string {
	var dir string
	for i, arg := range args {
		if arg == flag {
			dir = args[i+1]
		}
	}

	if strings.HasSuffix(dir, "/") {
		dir = dir[:len(dir)-1]
	}

	return dir
}
