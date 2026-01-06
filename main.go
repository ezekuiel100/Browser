package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

type Url struct {
	scheme string
	host   string
	path   string
}

func main() {
	url, err := parseUrl("http://localhost:3000")

	if err != nil {
		fmt.Println(err)
		return
	}

	request(*url)
}

func parseUrl(raw string) (*Url, error) {
	parts := strings.SplitN(raw, "://", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("URL inválida")
	}

	scheme := parts[0]
	if scheme != "http" {
		return nil, fmt.Errorf("a URL precisa usar o protocolo HTTP")
	}

	res := strings.SplitN(parts[1], "/", 2)
	host := res[0]

	path := "/"
	if len(res) == 2 {
		path = "/" + res[1]
	}

	return &Url{
		scheme: scheme,
		host:   host,
		path:   path,
	}, nil
}

func request(url Url) {
	conn, err := net.Dial("tcp", url.host)

	if err != nil {
		panic(err)
	}
	defer conn.Close()

	request := fmt.Sprintf(
		"GET %s HTTP/1.0\r\nHost: %s\r\n\r\n",
		url.path, url.host,
	)

	conn.Write([]byte(request))

	reader := bufio.NewReader(conn)

	statusLine, err := reader.ReadString('\n')
	if err != nil {
		panic(err)
	}

	parts := strings.SplitN(statusLine, " ", 3)
	fmt.Println(parts)

	response_headers := make(map[string]string)
	for {
		line, _ := reader.ReadString('\n')
		if line == "\r\n" {
			break
		}

		parts := strings.SplitN(line, ":", 2)
		header := strings.ToLower(parts[0])
		value := strings.TrimSpace(parts[1])

		response_headers[header] = value
	}

	fmt.Println(response_headers)
}
