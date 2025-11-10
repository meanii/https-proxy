# Usage — examples and quick tests

This document shows simple commands to build, run, and test the proxy locally. These examples assume you run the proxy on `localhost` port `9000`.

1) Build and run

```bash
go build -o https-proxy main.go
./https-proxy -address :9000
```

2) Test HTTP proxying with curl

```bash
curl -x http://localhost:9000 http://httpbin.org/get
```

This will send an HTTP request through the proxy and print the response from the origin server.

3) Test HTTPS tunneling with curl

```bash
curl -x http://localhost:9000 https://www.google.com -I
```

When curl uses an HTTP proxy for an HTTPS URL, it issues a `CONNECT` request to the proxy to establish a tunnel. This proxy implementation responds `HTTP/1.1 200 Connection Established` and then forwards raw bytes between the client and destination.

4) Browser testing

Set your browser's manual proxy settings to `localhost:9000` (HTTP proxy). Try browsing to an HTTP and HTTPS site. Remember this proxy does not inspect TLS content; HTTPS remains end-to-end encrypted between your browser and the origin.

5) Notes

- If you see failures or hangs when establishing tunnels, check firewall or platform proxy settings.
- This proxy logs minimal errors to stdout and is intended for debugging/learning only.
