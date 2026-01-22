package main

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type Url struct {
	scheme string
	host   string
	path   string
}

func main() {
	url, err := parseUrl("https://www.google.com")

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

	text := widget.NewLabel("Carregando...")
	text.Wrapping = fyne.TextWrapWord

	window.SetContent(container.NewVScroll(text))

	go func() {
		var b strings.Builder
		inTag := false

		for _, c := range body {
			if c == '<' {
				inTag = true
				continue
			} else if c == '>' {
				inTag = false
				continue
			}

			if !inTag {
				b.WriteRune(c)
			}
		}

		finalText := strings.Join(strings.Fields(b.String()), " ")
		text.SetText(finalText)

	}()

	window.Resize(fyne.NewSize(1200, 800))
	window.ShowAndRun()
}

func parseUrl(raw string) (*Url, error) {
	parts := strings.SplitN(raw, "://", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("URL inválida")
	}

	scheme := parts[0]
	if scheme != "http" && scheme != "https" {
		return nil, fmt.Errorf("protocolo não suportado: %s", scheme)
	}

	res := strings.SplitN(parts[1], "/", 2)
	host := res[0]

	if scheme == "https" && !strings.Contains(host, ":") {
		host = host + ":443"
	} else if scheme == "http" && !strings.Contains(host, ":") {
		host = host + ":80"
	}

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
	rawConn, err := net.Dial("tcp", url.host)
	if err != nil {
		return nil, err
	}

	var conn net.Conn = rawConn

	if url.scheme == "https" {
		serverName := strings.Split(url.host, ":")[0]

		tlsConn := tls.Client(rawConn, &tls.Config{
			ServerName: serverName,
		})

		err := tlsConn.Handshake()
		if err != nil {
			rawConn.Close()
			return nil, fmt.Errorf("erro no handshake TLS: %v", err)
		}

		conn = tlsConn
	}

	defer conn.Close()

	hostHeader := strings.Split(url.host, ":")[0]
	request := fmt.Sprintf(
		"GET %s HTTP/1.0\r\nHost: %s\r\n\r\n",
		url.path, hostHeader,
	)

	conn.Write([]byte(request))

	reader := bufio.NewReader(conn)

	statusLine, err := reader.ReadString('\n')
	if err != nil {
		return nil, err
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
