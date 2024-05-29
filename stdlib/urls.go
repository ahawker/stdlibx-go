package stdlib

import (
	"fmt"
	"net/url"
	"strconv"
)

var (
	// Mapping of url scheme to default well-known port number.
	//
	// https://en.wikipedia.org/wiki/List_of_TCP_and_UDP_port_numbers
	// https://www.iana.org/assignments/uri-schemes/uri-schemes.xhtml
	// https://en.wikipedia.org/wiki/List_of_URI_schemes
	schemaPortMapping = map[string]int{
		"ftp":   21,
		"ssh":   22,
		"smtp":  25,
		"http":  80,
		"https": 443,
	}
)

// URLPort returns a port number for the given URL.
//
// If a port was not defined during URL creation, an attempt is made
// to derive it from the scheme.
func URLPort(u *url.URL) (int, error) {
	port := u.Port()

	switch port {
	case "":
		if port, ok := schemaPortMapping[u.Scheme]; ok {
			return port, nil
		}
		return 0, fmt.Errorf("unknown port for url scheme=%s", u.Scheme)
	default:
		portNum, err := strconv.ParseInt(port, 10, 32)
		if err != nil {
			return 0, err
		}
		return int(portNum), nil
	}
}
