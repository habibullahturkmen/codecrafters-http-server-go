package main

import (
	"bufio"
	"bytes"
	"compress/gzip"
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

func handleGET(req Request, dirName string) (string, string, error) {
	var responseHeader string
	var responseBody string

	if req.path == "/" {
		responseHeader = fmt.Sprintf("%v 200 OK\r\n\r\n", req.httpVersion)
		return responseHeader, responseBody, nil
	}

	if strings.HasPrefix(req.path, "/echo") {
		content := strings.TrimPrefix(req.path, "/echo")

		if strings.HasPrefix(content, "/") {
			content = content[1:]
		}

		if strings.Contains(req.headers["Accept-Encoding"], "gzip") {
			var buf bytes.Buffer
			gz := gzip.NewWriter(&buf)
			_, err := gz.Write([]byte(content))
			if err != nil {
				return "", "", err
			}
			err = gz.Close()
			if err != nil {
				return "", "", err
			}

			responseHeader = fmt.Sprintf("%s 200 OK\r\nContent-Type: text/plain\r\nContent-Encoding: %s\r\nContent-Length: %d\r\n\r\n%s", req.httpVersion, "gzip", len(content), content)
			responseBody = string(buf.Bytes())
			return responseHeader, responseBody, nil
		}

		responseHeader = fmt.Sprintf("%s 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", req.httpVersion, len(content), content)
		return responseHeader, responseBody, nil
	}

	if strings.TrimRight(req.path, "/") == "/user-agent" {
		userAgent := req.headers[userAgent]
		responseHeader = fmt.Sprintf("%s 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", req.httpVersion, len(userAgent), userAgent)
		return responseHeader, responseBody, nil
	}

	if strings.HasPrefix(req.path, "/files") {
		fileName := strings.TrimPrefix(req.path, "/files")

		if strings.HasPrefix(fileName, "/") {
			fileName = fileName[1:]
		}

		file, err := os.ReadFile(fmt.Sprintf("%s/%s", dirName, fileName))
		if err != nil {
			if strings.Contains(err.Error(), "no such file or directory") {
				responseHeader = fmt.Sprintf("%v 404 Not Found\r\n\r\n", req.httpVersion)
				return responseHeader, responseBody, nil
			}
			return "", "", err
		}

		responseHeader = fmt.Sprintf("%s 200 OK\r\nContent-Type: application/octet-stream\r\nContent-Length: %d\r\n\r\n%s", req.httpVersion, len(string(file)), string(file))
		return responseHeader, responseBody, nil
	}

	responseHeader = fmt.Sprintf("%v 404 Not Found\r\n\r\n", req.httpVersion)
	return responseHeader, responseBody, nil
}

func handlePOST(req Request, dirName string) (string, error) {
	var responseHeader string
	if strings.HasPrefix(req.path, "/files") {
		fileName := strings.TrimPrefix(req.path, "/files")

		if strings.HasPrefix(fileName, "/") {
			fileName = fileName[1:]
		}

		err := os.WriteFile(fmt.Sprintf("%s/%s", dirName, fileName), req.body, 0666)
		if err != nil {
			return "", fmt.Errorf("failed creating file: %v", err.Error())
		}

		responseHeader = fmt.Sprintf("%v 201 Created\r\n\r\n", req.httpVersion)
		return responseHeader, nil
	}
	responseHeader = fmt.Sprintf("%v 404 Not Found\r\n\r\n", req.httpVersion)
	return responseHeader, nil
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
