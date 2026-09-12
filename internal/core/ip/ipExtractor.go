package ip

import (
	"fmt"
	"net"
)

func ExtractIPFromRequest(remoteAddr string) (string, error) {
	ip, _, err := net.SplitHostPort(remoteAddr)
	if err != nil {
		return "", fmt.Errorf("error parsing address: %q", remoteAddr)
	}
	return ip, nil
}
