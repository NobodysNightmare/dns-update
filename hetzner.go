package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
)

type UpdateRecordRequest struct {
	ZoneID string `json:"zone_id"`
	Type   string
	Name   string
	Value  string
	TTL    int
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
		fmt.Println("Unexpected status while updating:", resp.StatusCode)
	}
}
