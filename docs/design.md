# Design notes and code walkthrough

This file explains the small `main.go` implementation and highlights important points to learn from.

1) High-level flow

- The program opens a TCP listener on the configured address (default `:9000`).
- For each accepted connection it calls `connHandler` in a new goroutine.
- `connHandler` creates a cancellable context and calls `handleHttp` which reads an HTTP request from the connection using `http.ReadRequest`.

2) Handling regular HTTP requests

- For non-CONNECT requests the code creates a new `http.Request` with the client's method, URL and body. It copies headers and uses `http.DefaultTransport.RoundTrip` to perform the outbound request. The response is then written back to the client connection with `resp.Write(conn)`.

3) Handling CONNECT (HTTPS tunneling)

- When the method is `CONNECT`, the proxy expects the request `Host` to be the target host (for example `example.com:443`). If a port is missing it appends `:443`.
- It dials a TCP connection to the destination and replies to the client `HTTP/1.1 200 Connection Established\r\n\r\n`.
- After that it copies bytes bidirectionally between the client connection and the destination connection using `io.Copy` (one direction in a goroutine, the other in the current goroutine), forming a raw TCP tunnel.

4) Important learning points & caveats

- Error handling is minimal. For instance, when `http.DefaultTransport.RoundTrip` returns an error, the code attempts to call `resp.Write(conn)` — but `resp` may be nil when an error occurs, which could panic. This is a potential bug to fix.
- No authentication, no logging of sensitive contents, and no rate-limiting are implemented — which are necessary for production.
- The proxy does not inspect or manipulate TLS handshakes. CONNECT simply forwards encrypted bytes.

5) Small improvement ideas (exercise)

- Fix the error-path where `resp` might be nil before writing.
- Add logging hooks to inspect request/response headers safely.
- Add a whitelist or ACL and an optional authentication mechanism.
- Support concurrency limits and timeouts on idle tunnels.

6) Where to start modifying

- Open `main.go` and try the small change that checks `if err != nil { /* handle, don't use resp */ }` before using `resp`.
- Add a contribution-friendly `README` and tests if you want to expand this into a teaching repository.
