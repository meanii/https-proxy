package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"strings"
)

func main() {
	address := flag.String("address", ":9000", "address you want to run on (default :9000)")
	ln, err := net.Listen("tcp", *address)
	if err != nil {
		log.Fatal(err)
	}
	defer ln.Close()

	log.Printf("tunnel service running on %s", *address)
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("failed to read connection from [%s](%s)", conn.RemoteAddr(), err.Error())
			continue
		}
		go connHandler(conn)
	}
}

func connHandler(conn net.Conn) {
	defer conn.Close()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	handleHttp(ctx, conn)
}

func handleHttp(ctx context.Context, conn net.Conn) error {
	bu := bufio.NewReader(conn)

	httpRequst, err := http.ReadRequest(bu)
	if err != nil {
		if err != io.EOF {
			log.Printf("[Error] Failed to parse HTTP Request: %s", err)
		}
		return err
	}

	request, _ := http.NewRequestWithContext(ctx, httpRequst.Method, httpRequst.URL.String(), httpRequst.Body)

	request.Header = httpRequst.Header.Clone()

	// handling for HTTPS CONNECT method
	if request.Method == http.MethodConnect {
		host := request.Host
		if !strings.Contains(host, ":") {
			host += ":443"
		}

		// Dial the target server
		destConn, err := net.Dial("tcp", host)
		if err != nil {
			return err
		}

		_, err = fmt.Fprint(conn, "HTTP/1.1 200 Connection Established\r\n\r\n")
		if err != nil {
			destConn.Close()
			return err
		}

		// Tunnel data between client and destination
		go io.Copy(destConn, conn)
		io.Copy(conn, destConn)
		return nil
	}

	resp, err := http.DefaultTransport.RoundTrip(request)
	if err != nil {
		resp.Write(conn)
		return err
	}
	defer resp.Body.Close()

	// Write response to the client
	if err := resp.Write(conn); err != nil {
		log.Printf("[Error] failed to write response: %s", err)
		return err
	}

	return nil
}
