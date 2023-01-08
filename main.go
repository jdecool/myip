package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"strconv"
	"strings"
)

const (
	DEFAULT_HOST = "127.0.0.1"
	DEFAULT_PORT = 8080
)

func getIP(r *http.Request) (string, error) {
	ip := r.Header.Get("X-REAL-IP")
	netIP := net.ParseIP(ip)
	if netIP != nil {
		return ip, nil
	}

	ips := r.Header.Get("X-FORWARDED-FOR")
	splitIps := strings.Split(ips, ",")
	for _, ip := range splitIps {
		netIP := net.ParseIP(ip)
		if netIP != nil {
			return ip, nil
		}
	}

	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return "", err
	}

	netIP = net.ParseIP(ip)
	if netIP != nil {
		return ip, nil
	}

	return "", fmt.Errorf("No valid ip found")
}

func main() {
	var (
		serverHost string
		serverPort int
	)

	flag.StringVar(&serverHost, "host", DEFAULT_HOST, "")
	flag.IntVar(&serverPort, "port", DEFAULT_PORT, "")
	flag.Parse()

	addr := serverHost + ":" + strconv.Itoa(serverPort)
	fmt.Println("Starting server " + addr)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		ip, err := getIP(r)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(ip))
	})
	log.Fatal(http.ListenAndServe(addr, nil))
}
