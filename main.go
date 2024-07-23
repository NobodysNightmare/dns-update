package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"os"
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

func main() {
	config := readConfig()

	v6 := bestV6IP()
	fmt.Println("Detected IPv6:", v6)

	changed := hasIPv6Changed(v6)
	forcedUpdate := len(os.Args) == 2 && os.Args[1] == "--force"

	if changed || forcedUpdate {
		if changed {
			fmt.Printf("Change detected. ")
		}

		fmt.Printf("Updating AAAA records... ")
		for _, c := range config.V6Records {
			updateRecord(config.Token, c, "AAAA", v6)
		}
		fmt.Println("OK")
		fmt.Println()

		v4 := globalV4IP()
		fmt.Println("Detected IPv4:", v4)

		fmt.Printf("Updating A records... ")
		for _, c := range config.V4Records {
			updateRecord(config.Token, c, "A", v4)
		}
		fmt.Println("OK")

		updateKnownIPv6(v6)
	} else {
		fmt.Println("No change, exiting...")
	}
}

func readConfig() Config {
	file, _ := os.Open("config.json")
	configBytes, _ := io.ReadAll(file)

	var config Config
	if json.Unmarshal(configBytes, &config) != nil {
		fmt.Println("Couldn't read configuration. Using empty config!")
		fmt.Println()
	}

	return config
}

func hasIPv6Changed(currentIP net.IP) bool {
	return !knownIPv6().Equal(currentIP)
}

func knownIPv6() net.IP {
	file, _ := os.Open("known-ipv6.txt")
	defer file.Close()

	addrBytes, _ := io.ReadAll(file)
	return net.ParseIP(string(addrBytes))
}

func updateKnownIPv6(ip net.IP) {
	file, _ := os.Create("known-ipv6.txt")
	defer file.Close()

	file.Write([]byte(ip.String()))
}
