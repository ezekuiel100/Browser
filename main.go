package main

import (
	"fmt"
	"strings"
)

type Url struct {
	scheme string
	host   string
	path   string
}

func main() {
	url, err := parseUrl("http://exemplo.com/teste")

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(url)
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
