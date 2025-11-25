package main

import (
	"bufio"
	"fmt"
	"os"
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
