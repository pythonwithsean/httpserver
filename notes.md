# OS Concepts - Socket File Descriptors

## How sockets map to my HTTP server code

### Listening Socket FD (`net.Listen`)

`net.Listen("tcp", port)` calls the OS syscalls: `socket()` + `bind()` + `listen()`.

The OS gives the process a **listening socket file descriptor**. This FD is not for reading/writing data — it only exists to accept incoming connections.

### Connection Socket FD (`listener.Accept()`)

`listener.Accept()` calls the OS `accept()` syscall. It blocks until a client connects, then the OS creates a **brand new file descriptor** for that specific connection. This is the FD used to `read()`/`write()` with the remote machine.

### Two types of FDs in the process

| FD | Purpose |
|---|---|
| Listener FD | Waits for new connections (the "door") |
| Connection FD | Reads/writes data with a specific client (the "conversation") |

Each call to `Accept()` produces a new connection FD, which is why the `for` loop can handle many clients — each `conn` is a separate FD pointing to a different remote endpoint. The listener FD stays open the whole time, continuing to accept new connections.

`go handleConn(conn)` hands off that connection FD to a goroutine so the loop can go back to `Accept()` the next one.

## Client Address Format: `[::1]:53741`

When `conn.RemoteAddr()` returns `[::1]:53741`:

**`[::1]`** — IPv6 loopback address (same as `127.0.0.1` in IPv4). Means the connection came from localhost.

**`53741`** — Ephemeral port the OS assigned to the client's side of the connection (not the server's port 5000). Every TCP connection has two endpoints: `client_ip:client_port ↔ server_ip:server_port`. The OS picks a random high port for the client.

## Port Ranges

Total: 0–65535 = **65,536 ports** per IP address.

| Range | Name |
|---|---|
| 0–1023 | Well-known (HTTP=80, SSH=22, etc.) |
| 1024–49151 | Registered |
| 49152–65535 | Ephemeral (OS assigns these to clients) |

## Go Packages: `package main` is not importable

`package main` is the executable/program — it has a `main()` function and gets compiled into a binary. Go **forbids** importing `package main` from other packages. It's not a library.

To test code in `package main`:
- Put test files in the same directory with `package main`
- Call functions directly, no import needed
- Run with `go test .`

To make code importable by other packages:
- Use a different package name (e.g., `package httpserver`)
- Export functions with capital letters
- Import via the module path
