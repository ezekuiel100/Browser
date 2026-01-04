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
	data, err := parseUrl("http://exemplo.com/teste")

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(data)
}

func parseUrl(row string) (Url, error) {
	var u Url
	parts := strings.SplitN(row, "://", 2)
	if len(parts) != 2 {
		return u, fmt.Errorf("URL inválida")
	}

	u.scheme = parts[0]
	if u.scheme != "http" {
		return u, fmt.Errorf("a URL precisa usar o protocolo HTTP")
	}

	res := strings.SplitN(parts[1], "/", 2)
	u.host = res[0]

	if len(res) == 1 {
		u.path = "/"
	} else {
		u.path = "/" + res[1]
	}

	return u, nil
}
