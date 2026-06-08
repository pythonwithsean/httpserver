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

## Escape Characters

| Escape | Name | What it does |
|---|---|---|
| `\r` | Carriage return | Moves cursor to index 0 (start of line) |
| `\n` | Newline / line feed | Moves cursor down one line |
| `\t` | Tab | Horizontal tab |

`\r` prevents weird spacing by resetting to the beginning of the line before moving down.

HTTP uses `\r\n` (CRLF) as line endings — a holdover from old teletype machines that needed both: return to start (`\r`) then move down (`\n`).

## Go Modules & `go get`

### Module Path (`go.mod`)

The module path like `github.com/pythonwithsean/httpserver` is just a unique identifier to avoid naming collisions. It doesn't make anything public — repo visibility on GitHub controls that.

### How `go get` works

```bash
go get github.com/someone/package
```

1. Fetches the repo from GitHub
2. Reads its `go.mod` to get dependencies
3. Downloads it to your **module cache** (`~/go/pkg/mod/`)
4. Adds it to your `go.mod` under `require`

### Go Binary vs Module Cache

**Go binary location:**
```
/opt/homebrew/bin/go    ← the Go compiler/toolchain itself
```

**Module cache (where `go get` installs packages):**
```
~/go/pkg/mod/           ← downloaded dependencies live here
```

**GOPATH (legacy workspace):**
```
~/go/                   ← default GOPATH
├── bin/                ← compiled binaries from `go install`
├── pkg/mod/            ← module cache (dependencies)
└── src/                ← legacy workspace (rarely used now)
```

When you `go get` a package, it goes into `~/go/pkg/mod/`, not into the Go binary. The Go compiler (`/opt/homebrew/bin/go`) reads from the cache when building your project.

### Private Repos

Make the GitHub repo private. To `go get` from private repos:
```bash
git config --global url."git@github.com:".insteadOf "https://github.com/"
go env -w GOPRIVATE=github.com/yourusername/*
```

## `go get` vs `go install`

| Command | Purpose | Modifies `go.mod`? |
|---|---|---|
| `go get` | Adds a dependency to your project | Yes |
| `go install` | Installs a standalone CLI tool/binary | No |

```bash
go get github.com/some/lib              # adds lib to your go.mod
go install github.com/air-verse/air@latest  # installs the `air` binary to ~/go/bin/
```

**TL;DR:** `go get` = add a library to your project. `go install` = install a standalone tool.

## `go.sum` — Dependency Checksums

`go.sum` records cryptographic hashes of every module version your project depends on.

- `go.mod` = *what* dependencies you need
- `go.sum` = *proof* those dependencies are authentic

Go verifies checksums on every build to ensure downloaded dependencies haven't been tampered with. Never manually edit `go.sum` — it's managed automatically by `go get`, `go mod tidy`, etc.

## PATH & Zsh Config

### What is PATH?

`PATH` is an environment variable — a list of directories (separated by `:`) that tells your shell **where to look for executables** when you type a command.

```bash
echo $PATH
# /opt/homebrew/bin:/usr/local/bin:/usr/bin:/bin:/Users/Sean/go/bin
```

When you type `air`, zsh checks each directory **left to right** until it finds a binary named `air`. If it's not in any of them → `command not found`.

### Adding to PATH

```bash
export PATH="$HOME/go/bin:$PATH"
#           ^^^^^^^^^^^^^ ^^^^^^
#           add this dir   keep everything already in PATH
```

Without `:$PATH` at the end, you'd **replace** the entire PATH and break every command.

### Zsh config load order

| File | When it runs |
|---|---|
| `~/.zshenv` | Always, every shell |
| `~/.zprofile` | Login shells only |
| `~/.zshrc` | Interactive shells (your normal terminal) — **this is the one you edit** |
| `~/.zlogin` | Login shells, after `.zshrc` |

### Permanent vs temporary changes

- **Permanent:** Add `export` lines to `~/.zshrc` — runs every time you open a new terminal
- **Temporary:** Run `export PATH="..."` directly in the terminal — gone when you close it
- **Reload without restarting:** `source ~/.zshrc`

### Useful commands

| Command | What it does |
|---|---|
| `echo $PATH` | See your current PATH |
| `which go` | See which directory `go` was found in |
| `source ~/.zshrc` | Reload config without closing terminal |
| `export KEY=value` | Set a variable (temporary, current session only) |

### Key variables

| Variable | Meaning |
|---|---|
| `$HOME` or `~` | Your home directory (`/Users/Sean`) |
| `$PATH` | Directories to search for commands |
| `$GOPATH` | Go workspace (`~/go`) |
