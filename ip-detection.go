package main

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
)

func bestV6IP() net.IP {
	addresses, err := net.InterfaceAddrs()
	if err != nil {
		return nil
	}

	for _, addr := range addresses {
		if ipnet, ok := addr.(*net.IPNet); ok {
			if ipnet.IP.IsGlobalUnicast() && !ipnet.IP.IsPrivate() {
				return ipnet.IP
			}
		}
	}

	return nil
}

func globalV4IP() net.IP {
	client := &http.Client{
		Transport: &http.Transport{
			DialContext: func(ctx context.Context, network string, addr string) (net.Conn, error) {
				return (&net.Dialer{}).DialContext(ctx, "tcp4", addr)
			},
		},
	}

	resp, err := client.Get("https://ifconfig.co/")
	if err != nil {
		fmt.Println(err)
		return nil
	}

	addrBytes, _ := io.ReadAll(resp.Body)
	addrString := strings.TrimSpace(string(addrBytes))
	return net.ParseIP(addrString)
}
