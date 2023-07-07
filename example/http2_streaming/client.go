package main

import (
	"crypto/tls"
	"fmt"
	"io"
	"net/http"

	"golang.org/x/net/http2"
)

func main() {
	// Create a custom TLS configuration
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
		// Add any other custom TLS configuration options here
	}

	// Create a new HTTP/2 Transport with the custom TLS configuration
	tr := &http2.Transport{
		TLSClientConfig: tlsConfig,
	}

	// Create an HTTP client using the custom Transport
	client := &http.Client{
		Transport: tr,
	}

	// Create an HTTP GET request
	req, err := http.NewRequest("GET", "https://192.168.56.1:4243/", nil)
	if err != nil {
		fmt.Println("Failed to create request:", err)
		return
	}

	// Send the HTTP request
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Failed to send request:", err)
		returnx
	}
	defer resp.Body.Close()

	// Read the response body in chunks
	buf := make([]byte, 1)
	for {
		_, err := resp.Body.Read(buf)
		if err != nil {
			if err != io.EOF {
				fmt.Println("Failed to read response body:", err)
			}

			break
		}

		fmt.Print(string(buf[0]))
	}
}
