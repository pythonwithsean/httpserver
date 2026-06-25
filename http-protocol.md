# HTTP/1.1 — A Practical Reference for Building This Server

> This is not the actual RFC — it's a distilled, implementer-focused version of
> **RFC 9110** (HTTP Semantics) and **RFC 9112** (HTTP/1.1 message syntax),
> organized around what you actually need to parse/build requests and responses.
> When something here feels incomplete, that section number is where to go read the real thing.

---

## 1. The Message Shape

Both requests and responses share the same skeleton:

```
start-line CRLF
header-field CRLF
header-field CRLF
...
CRLF                  <- blank line: end of headers
[ message body ]
```

**Request start-line:**
```
method SP request-target SP HTTP-version CRLF
GET /users/1 HTTP/1.1\r\n
```

**Response start-line (status-line):**
```
HTTP-version SP status-code SP reason-phrase CRLF
HTTP/1.1 200 OK\r\n
```

There is no `\r\n` after the body. The body's exact length is declared by `Content-Length`
(or delimited by chunked encoding) — nothing follows it.

---

## 2. Mandatory vs Optional Headers

### Request headers

| Header | Required? | Why |
|---|---|---|
| `Host` | **Mandatory** (HTTP/1.1+) | Lets one server/IP serve multiple domains (virtual hosting). RFC 9112 §3.2 — a request without it is invalid. HTTP/1.0 didn't require this. |
| `Content-Length` *or* `Transfer-Encoding: chunked` | Mandatory **if a body is present** | Without one of these, the receiver can't know where the body ends. Never send both. |
| `User-Agent`, `Accept`, `Content-Type`, `Authorization`, `Cookie`, `Connection`, custom `X-*` | Optional | Useful, but the protocol doesn't require them. |

### Response headers

| Header | Required? | Why |
|---|---|---|
| `Content-Length` *or* `Transfer-Encoding: chunked` | Mandatory **if there's a body** and the connection stays open | Same framing reasoning as requests. |
| `Date` | Recommended (SHOULD) | RFC 9110 §6.6.1 — origin servers with a clock should include it. |
| `Content-Type` | Strongly recommended if there's a body | Tells the client how to interpret the bytes. No `Content-Type` means the client has to guess. |
| `Server`, `Cache-Control`, `Set-Cookie`, `ETag`, custom headers | Optional | Purely situational. |

---

## 3. Header Field Syntax

```
field-line = field-name ":" OWS field-value OWS
```
`OWS` = optional whitespace (`*( SP / HTAB )`).

**Field name (`token`)** — letters, digits, and exactly these symbols, nothing else:
```
! # $ % & ' * + - . ^ _ ` | ~
```
No spaces, no colons, no control characters.

**Field value** — visible US-ASCII only: byte range `0x21`-`0x7E`, plus space/tab as separators.
Reject anything outside that range — control characters (including raw `\r`/`\n`), DEL (`0x7F`),
and any byte `≥ 0x80` (non-ASCII). This is what stops CRLF injection / header smuggling: if you
never let a literal `\r\n` exist *inside* what should be a single header's value, an attacker can't
forge an extra header or split your message framing.

See `notes.md` → "Header Injection / Unicode Smuggling" for the byte-range table and the Go
implementation (`IsValidHeaderValue`, `IsValidHeaderKey` in `server/parser.go`).

---

## 4. Body Framing Rules

- If `Content-Length: N` is present, the body is exactly `N` bytes — read exactly that many, no more, no less.
- If `Transfer-Encoding: chunked` is present instead, the body is sent in chunks:
  ```
  <chunk-size-in-hex>\r\n
  <chunk-data>\r\n
  ... (repeat)
  0\r\n
  \r\n
  ```
  A `0`-sized chunk marks the end.
- **Never trust both at once.** A message with both `Content-Length` and `Transfer-Encoding` is
  exactly the classic HTTP request smuggling vector (different systems in a chain may honor
  different headers) — RFC 9112 §6.3 says to reject such messages outright.
- No body at all → no `Content-Length` needed (e.g. most `GET` requests, `204 No Content` responses).

---

## 5. Methods

| Method | Safe? (no side effects) | Idempotent? (same result if repeated) |
|---|---|---|
| `GET` | Yes | Yes |
| `HEAD` | Yes | Yes |
| `OPTIONS` | Yes | Yes |
| `PUT` | No | Yes |
| `DELETE` | No | Yes |
| `POST` | No | No |
| `PATCH` | No | No |

"Safe" = doesn't change server state. "Idempotent" = calling it twice has the same effect as calling it once (important for retries — clients can safely retry idempotent methods after a timeout, but not `POST`).

---

## 6. Status Code Categories

| Range | Meaning | Examples |
|---|---|---|
| `1xx` | Informational | `100 Continue` |
| `2xx` | Success | `200 OK`, `201 Created`, `204 No Content` |
| `3xx` | Redirection | `301 Moved Permanently`, `304 Not Modified` |
| `4xx` | Client error | `400 Bad Request`, `404 Not Found`, `413 Payload Too Large` |
| `5xx` | Server error | `500 Internal Server Error`, `503 Service Unavailable` |

---

## 7. Security Pitfalls to Remember

- **CRLF injection**: unvalidated header values containing raw `\r\n` → forged headers, response splitting. Fixed via byte-range validation (see §3).
- **Request smuggling**: ambiguity between `Content-Length` and `Transfer-Encoding`, or between how two systems in a chain interpret the same bytes. Fixed by strict, single-source-of-truth body framing (see §4).
- **Reflecting unvalidated data into your own response headers** reopens injection even if your *parsing* is solid — validation on the way in doesn't protect data you build on the way out, unless that data already passed through the same validated path (e.g. `req.Host` is safe to reuse because it was validated at parse time).

---

## 8. Where the Real Spec Lives

- **RFC 9110** — HTTP Semantics (headers, methods, status codes — "what things mean")
- **RFC 9112** — HTTP/1.1 (wire format, framing, chunked encoding — "how bytes are laid out")
- **RFC 9111** — HTTP Caching (`Cache-Control`, `ETag`, etc.)

These replaced RFC 7230-7235 (2014), which replaced the original RFC 2616 (1999). Search e.g.
"rfc 9112 host" to jump straight to a section instead of reading linearly.
