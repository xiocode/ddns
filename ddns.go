package main

import (
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const TIMEOUT = 10

var (
	body = url.Values{
		"login_token": {"0"},
		"format":      {"json"},
		"domain_id":   {"0"},
		"record_id":   {"0"},
		"sub_domain":  {"0"},
		"record_line": {"默认"},
	}
	current_ip     = ""
	check_interval = 30 * time.Second
)

func get_public_ip() (string, error) {
	conn, err := net.DialTimeout("tcp", "ns1.dnspod.net:6666", TIMEOUT*time.Second)
	defer func() {
		if x := recover(); x != nil {
			log.Println("Can't get public ip", x)
		}
		if conn != nil {
			conn.Close()
		}
	}()

	if err == nil {
		deadline := time.Now().Add(TIMEOUT * time.Second)
		err = conn.SetDeadline(deadline)
		if err != nil {
			return "", err
		}
		var bytes []byte
		bytes, err = ioutil.ReadAll(conn)
		if err == nil {
			return string(bytes), nil
		}
	}
	return "", err
}

func update_dnspod(ip string) bool {
	client := new(http.Client)
	body.Set("value", ip)

	req, err := http.NewRequest("POST", "https://dnsapi.cn/Record.Ddns", strings.NewReader(body.Encode()))
	req.Header.Set("Accept", "text/json")
	req.Header.Set("Content-type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	defer resp.Body.Close()
	if err != nil {
		return false
	}

	bytes, _ := ioutil.ReadAll(resp.Body)
	log.Println(string(bytes))
	return resp.StatusCode == 200
}

func main() {
	for {
		ip, err := get_public_ip()
		if err == nil {
			if ip != current_ip {
				if update_dnspod(ip) {
					current_ip = ip
				}
			}
		}
		time.Sleep(check_interval)
	}
}
