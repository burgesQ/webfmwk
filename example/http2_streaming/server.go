package main

import (
	"bufio"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"time"

	fasthttp2 "github.com/dgrr/http2"
	"github.com/valyala/fasthttp"
)

func main() {
	tlsCfg := loadTLS()

	server := &fasthttp.Server{
		Handler: func(ctx *fasthttp.RequestCtx) {
			// requestHandler := func(ctx *fasthttp.RequestCtx) {
			// Définir les en-têtes HTTP appropriés pour le streaming
			ctx.Response.Header.Set("Content-Type", "text/plain")
			// ctx.Response.Header.Set("Pragma", "no-cache")
			// ctx.Response.Header.Set("Expires", "0")
			// ctx.Response.Header.Set("Access-Control-Allow-Origin", "*")
			// ctx.Response.Header.Set("Access-Control-Allow-Headers", "Content-Type")
			// ctx.Response.Header.Set("Access-Control-Allow-Methods", "GET")

			// execute the streaming logic in a separate goroutine
			// go streamContent(ctx.Response.BodyWriter())

			ctx.Response.SetBodyStreamWriter(func(w *bufio.Writer) {
				defer w.Flush()

				defer fmt.Println("body stream done")

				for i := 0; i < 10; i++ {
					data := []byte(fmt.Sprintf("Chunk %d\n", i))

					_, err := w.Write(data)
					if err != nil {
						fmt.Println("Failed to write data:", err)
						return
					}

					// Flush the writer to ensure the data is sent immediately
					err = w.Flush()
					if err != nil {
						fmt.Println("Failed to flush writer:", err)
						return
					}

					time.Sleep(time.Second / 2) // Simulate generation delay
				}
				// ctx.Response.BodyWriteTo(ctx)
			})

			// Create a pipe
			// pr, pw := io.Pipe()

			// Set the response body as the ReadWriteCloser
			// ctx.Response.SetBodyStream(pr, -1)

			//	done := make(chan bool)

			// ctx.Response.SetBodyStreamWriter(func(w *bufio.Writer) {
			// 	for {
			// 		data := []byte{}
			// 		i, _ := pr.Read(data)
			// 		if _, e := w.Write(data[:i]); e != nil {
			// 			fmt.Println("done writting: %v", e)

			// 			return
			// 		}
			// 	}
			// })

			// go func(w io.WriteCloser) {
			// 	writer := bufio.NewWriter(w)

			// 	defer func() {
			// 		writer.Flush()
			// 		pw.Close()
			// 		fmt.Println("flushed closed")

			// 		pr.Close()
			// 		fmt.Println("closing other side")

			// 		ctx.Response.CloseBodyStream() //  Conn().Close()
			// 	}()

			// 	for i := 0; i < 10; i++ {
			// 		data := []byte(fmt.Sprintf("Chunk %d\n", i))

			// 		_, err := writer.Write(data)
			// 		if err != nil {
			// 			fmt.Println("Failed to write data:", err)
			// 			return
			// 		}
			// 		// time.Sleep(time.Second) // Simulate generation delay
			// 	}
			// }(pw)

			// go func() {
			// 	<-done

			// }()
			// ctx.Response.CloseBodyStream()
			// Close the response body after the goroutine has completed
			// defer ctx.Response.BodyClose()

			// // Access the underlying connection
			// ctx.Hijack(func(c net.Conn) {
			// 	// fmt.Fprintf(c, "This message is sent over a hijacked connection to the client %s\n", c.RemoteAddr())
			// 	// fmt.Fprintf(c, "Send me something and I'll echo it to you\n")

			// 	fmt.Println("waiting for gen")
			// 	<-done
			// 	fmt.Println("closing conn")
			// 	// Close the connection manually
			// 	if err := c.Close(); err != nil {
			// 		log.Printf("Failed to close connection: %v", err)
			// 	}
			// })

			// // If you need to perform any additional cleanup after streaming,
			// // you can close the write end of the pipe when the request is done
			// ctx.Response.SetConnectionClose()

			// // ctx.Request.SetConnectionClose()

			// io.ReadCloser

			// ctx.Response.SetBodyStream(, -1)
		},
	}

	// fasthttp2.ConfigureServerAndConfig(server, &tlsCfg)
	fasthttp2.ConfigureServer(server, fasthttp2.ServerConfig{Debug: true})

	// fmt.Printf("\n%+v\n", tlsCfg)

	addr := "192.168.56.1:4243"

	listener, err := net.Listen("tcp4", addr)
	if err != nil {
		panic(err)
	}

	fmt.Printf("listing on https://%q\n", addr)

	tlsListener := tls.NewListener(listener, tlsCfg)

	if err := server.Serve(tlsListener); err != nil {
		panic(err)
	}
}

// Function to generate or stream the content
func streamContent(w io.Writer) {
	// For simplicity, let's use a simple loop here
	for i := 0; i < 10; i++ {
		// Generate a chunk of data
		chunk := []byte(fmt.Sprintf("Chunk %d\n", i))

		// Write the chunk to the writer
		_, err := w.Write(chunk)
		if err != nil {
			return
		}
	}

	// Manually close the connection after streaming is complete
	if closer, ok := w.(io.Closer); ok {
		closer.Close()
	}
}

//
// TLS
//

func loadTLS() *tls.Config {
	cert, err := tls.LoadX509KeyPair(
		"/home/master/repo/certs/pkapman/ssl.crt",
		"/home/master/repo/certs/pkapman/ssl.key")
	if err != nil {
		panic(fmt.Errorf("cannot load cert and key: %w", err))
	}

	tlsCfg := &tls.Config{
		NextProtos:   []string{"h2"},
		MinVersion:   tls.VersionTLS12,
		MaxVersion:   tls.VersionTLS13,
		ClientAuth:   tls.NoClientCert,
		Certificates: []tls.Certificate{cert},
	}

	return tlsCfg
}

func loadMTLSVerifyPeer(tlsCfg *tls.Config) *tls.Config {
	tlsCfg.GetConfigForClient = func(*tls.ClientHelloInfo) (*tls.Config, error) {
		return &tls.Config{
			NextProtos: []string{"h2"},
			MinVersion: tls.VersionTLS12,
			MaxVersion: tls.VersionTLS13,
		}, nil
	}

	return tlsCfg
}

func loadCA(tlsCfg *tls.Config) *tls.Config {
	pool := x509.NewCertPool()

	if caCertPEM, e := os.ReadFile("/home/master/repo/certs/pkapman/cacert.pem"); e != nil {
		panic(fmt.Errorf("cannot load ca cert in pool: %w", e))
	} else if !pool.AppendCertsFromPEM(caCertPEM) {
		panic(errors.New("cannot apend cert to pem"))
	}

	tlsCfg.ClientCAs = pool

	return tlsCfg
}
