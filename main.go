package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"strings"
)

type Config struct {
	Token     string
	V6Records []RecordConfig
	V4Records []RecordConfig
}

type RecordConfig struct {
	Zone   string
	Record string
	Name   string
}

type UpdateRecordRequest struct {
	ZoneID string `json:"zone_id"`
	Type   string
	Name   string
	Value  string
	TTL    int
}

func main() {
	file, _ := os.Open("config.json")
	configBytes, _ := io.ReadAll(file)

	var config Config
	json.Unmarshal(configBytes, &config)

	v6 := bestV6IP()
	v4 := globalV4IP()

	fmt.Println("Detected IPv6:", v6)
	fmt.Println("Detected IPv4:", v4)

	// TODO: figure out if there was a change

	for _, c := range config.V6Records {
		updateRecord(config.Token, c, "AAAA", v6)
	}

	for _, c := range config.V4Records {
		updateRecord(config.Token, c, "A", v4)
	}
}

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

func updateRecord(token string, config RecordConfig, recordType string, ip net.IP) {
	reqBody := UpdateRecordRequest{
		ZoneID: config.Zone,
		Type:   recordType,
		TTL:    30,
		Name:   config.Name,
		Value:  ip.String(),
	}
	jsonBytes, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest(http.MethodPut, fmt.Sprintf("https://dns.hetzner.com/api/v1/records/%s", config.Record), bytes.NewReader(jsonBytes))
	req.Header.Set("Auth-API-Token", token)
	resp, _ := http.DefaultClient.Do(req)
	if resp.StatusCode != 200 {
		fmt.Println(resp.StatusCode)
	}
}
