# https-proxy

A tiny learning-purpose HTTP/HTTPS proxy written in Go. This project shows a small, dependency-free proxy that handles normal HTTP requests and the HTTP CONNECT method for tunneling HTTPS traffic.

This repository is intended for experimentation and learning only — it is not production-ready. It demonstrates how to accept TCP connections, parse HTTP requests, perform CONNECT tunneling, and proxy HTTP responses using Go's standard library.

## Features

- Minimal, single-file implementation (`main.go`).
- Supports plain HTTP proxying (for GET/POST/etc.) and HTTPS tunneling via HTTP CONNECT.
- No external dependencies — uses only the Go standard library.

## Requirements

- Go 1.18+ installed

## Build

From the repository root:

```bash
go build -o https-proxy main.go
```

## Run

Run the binary (default listens on `:9000`):

```bash
./https-proxy -address :9000
```

You can change the listening address/port via `-address` flag.

## Usage examples

- HTTP request via the proxy (using curl):

```bash
curl -x http://localhost:9000 http://example.com/
```

- HTTPS request via the proxy (curl will use CONNECT under the hood):

```bash
curl -x http://localhost:9000 https://www.google.com -I
```

- Configure your browser or system proxy to `localhost:9000` to route traffic through this proxy. Remember this implementation does not perform TLS interception; it simply tunnels bytes for CONNECT requests.

## Limitations / Warnings

- No authentication or access controls.
- No TLS interception (no MITM). CONNECT opens a raw TCP tunnel to the destination.
- Error handling is minimal and there's no request/response filtering or caching.
- Not hardened for production use. Use only in safe learning environments.

## Project layout

- `main.go` — single-file proxy implementation.
- `docs/` — extended documentation (usage, design notes).

## Learning notes

This project is useful to learn about:

- How TCP listeners and connections are handled in Go.
- Parsing HTTP requests with `http.ReadRequest`.
- Implementing an HTTP CONNECT tunnel for HTTPS traffic.
- Using `http.Transport` to perform outbound HTTP requests.

See the `docs` folder for more details, example commands, and a short code walkthrough.

## License

This repo is provided for learning purposes. Use and modify freely.
