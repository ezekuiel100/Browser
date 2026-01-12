package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/widget"
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

	body, err := request(*url)
	if err != nil {
		log.Fatal(err)
	}

	show(string(body))
}

func show(body string) {
	App := app.New()
	window := App.NewWindow("Browser")

	inTag := false
	var result []rune

	for _, c := range body {
		if c == '<' {
			inTag = true
			continue
		} else if c == '>' {
			inTag = false
			continue
		}

		if !inTag {
			result = append(result, c)
		}
	}

	text := widget.NewLabel(string(result))

	window.SetContent(text)
	window.Resize(fyne.NewSize(800, 600))
	window.ShowAndRun()
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

func request(url Url) ([]byte, error) {
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

	responseHeaders := make(map[string]string)
	for {
		line, _ := reader.ReadString('\n')
		if line == "\r\n" {
			break
		}

		parts := strings.SplitN(line, ":", 2)
		header := strings.ToLower(parts[0])
		value := strings.TrimSpace(parts[1])

		responseHeaders[header] = value
	}

	if _, ok := responseHeaders["transfer-encoding"]; ok {
		return nil, fmt.Errorf("transfer-encoding não suportado")
	}

	if _, ok := responseHeaders["content-encoding"]; ok {
		return nil, fmt.Errorf("content-encoding não suportado")
	}

	body, _ := io.ReadAll(reader)

	return body, nil
}
