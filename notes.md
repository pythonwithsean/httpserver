# Go Notes

## make() vs nil

`make()` **never** returns nil. It initializes slices, maps, and channels:

```go
s := make([]int, 0)    // non-nil empty slice
m := make(map[string]int) // non-nil empty map
c := make(chan int)       // non-nil channel
```

You get `nil` with the **zero value** (no `make`):

```go
var s []int              // nil slice
var m map[string]int     // nil map
var c chan int           // nil channel
```

| | `nil` | `make()` |
|---|---|---|
| Read from map | returns zero value | returns zero value |
| Write to map | **panics** | works |
| Slice `append` | works | works |
| Slice len/cap | 0/0 | 0/0 (but non-nil) |
| `== nil` | true | false |

---

## Interfaces

```go
type Writer interface {
    Write(p []byte) (n int, err error)
}

type FileWriter struct{}

func (f FileWriter) Write(p []byte) (int, error) {
    return len(p), nil
}

// FileWriter implicitly satisfies Writer
var w Writer = FileWriter{}
```

---

## Method Overriding (Embedding)

Go uses **embedding** instead of inheritance:

```go
type Base struct{}
func (b Base) Greet() string { return "hello" }

type Child struct {
    Base
}
func (c Child) Greet() string { return "hi" } // overrides

c := Child{}
c.Greet()       // "hi"
c.Base.Greet()  // "hello"
```

---

## Generics (Go 1.18+)

```go
func Map[T any, U any](s []T, f func(T) U) []U {
    result := make([]U, len(s))
    for i, v := range s {
        result[i] = f(v)
    }
    return result
}
```

### Type Constraints

```go
func Min[T constraints.Ordered](a, b T) T {
    if a < b { return a }
    return b
}
```

---

## Sorting (`sort` package)

```go
import "sort"

sort.Ints([]int{3,1,2})
sort.Strings([]string{"b","a","c"})
sort.Float64s([]float64{3.1, 1.2})

sort.Slice(people, func(i, j int) bool {
    return people[i].Age < people[j].Age
})

sort.IntsAreSorted([]int{1,2,3}) // true
```

---

## Heaps (`container/heap`)

```go
import "container/heap"

type IntHeap []int
func (h IntHeap) Len() int           { return len(h) }
func (h IntHeap) Less(i, j int) bool { return h[i] < h[j] } // min-heap
func (h IntHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h *IntHeap) Push(x any)        { *h = append(*h, x.(int)) }
func (h *IntHeap) Pop() any {
    old := *h
    n := len(old)
    x := old[n-1]
    *h = old[:n-1]
    return x
}

h := &IntHeap{3, 1, 2}
heap.Init(h)
heap.Push(h, 0)
min := heap.Pop(h) // 0
```

---

## Doubly Linked List (`container/list`)

```go
import "container/list"

l := list.New()
l.PushBack(1)
l.PushBack(2)
l.PushFront(0)

for e := l.Front(); e != nil; e = e.Next() {
    fmt.Println(e.Value)
}

l.Remove(l.Front())
l.Len() // 2
```

---

## Circular List (`container/ring`)

```go
import "container/ring"

r := ring.New(3) // ring of size 3
r.Value = "a"
r = r.Next()
r.Value = "b"
r = r.Next()
r.Value = "c"

r.Do(func(v any) {
    fmt.Println(v) // a, b, c
})

r.Move(1)  // move forward 1
r.Len()    // 3
```

---

## Slice Methods (`slices` package, Go 1.21+)

```go
import "slices"

s := []int{3, 1, 2, 1}
slices.Sort(s)           // [1,1,2,3]
slices.Contains(s, 2)    // true
slices.Index(s, 1)       // 0
slices.Compact(s)        // [1,2,3] removes consecutive dupes
slices.Reverse(s)
slices.Min(s)            // 1
slices.Max(s)            // 3
slices.BinarySearch(s, 2)
```

---

## Map Methods (`maps` package, Go 1.21+)

```go
import "maps"

m1 := map[string]int{"a": 1}
m2 := map[string]int{"b": 2}
maps.Copy(m1, m2)        // m1 = {"a":1, "b":2}
maps.Equal(m1, m2)       // false
maps.DeleteFunc(m1, func(k string, v int) bool {
    return v < 2
})
```

---

## String Methods (`strings` package)

```go
import "strings"

strings.Contains("hello", "ell")   // true
strings.HasPrefix("hello", "he")   // true
strings.HasSuffix("hello", "lo")   // true
strings.Split("a,b,c", ",")        // ["a","b","c"]
strings.Join([]string{"a","b"}, "-") // "a-b"
strings.Replace("hello", "l", "r", 1) // "herlo"
strings.ReplaceAll("hello", "l", "r") // "herro"
strings.Trim("  hi  ", " ")        // "hi"
strings.ToLower("HI")              // "hi"
strings.ToUpper("hi")              // "HI"
strings.Fields("a  b  c")          // ["a","b","c"]
```

### strings.Builder (efficient string building)

```go
var b strings.Builder
b.WriteString("hello")
b.WriteString(" world")
b.String() // "hello world"
```

---

## Runes

A rune is an alias for `int32` representing a Unicode code point:

```go
var r rune = 'A'  // same as int32 = 65
```

### `unicode` package

```go
import "unicode"

unicode.IsLetter('a')    // true
unicode.IsDigit('5')     // true
unicode.IsUpper('A')     // true
unicode.IsLower('a')     // true
unicode.IsSpace(' ')     // true
unicode.IsPunct('!')     // true
unicode.ToUpper('a')     // 'A'
unicode.ToLower('A')     // 'a'
unicode.ToTitle('a')     // 'A'
```

### `unicode/utf8` package

```go
import "unicode/utf8"

utf8.RuneCountInString("héllo")  // 5 (not 6 bytes)
utf8.RuneLen('é')                // 2 bytes
utf8.ValidString("hello")        // true
utf8.DecodeRuneInString("é")     // 'é', 2
```

### Converting between strings and runes

```go
s := "hello"
runes := []rune(s)        // string -> rune slice
str := string(runes)      // rune slice -> string

// Iterating runes in a string
for i, r := range "héllo" {
    fmt.Println(i, r, string(r))
}
```

### Byte vs Rune — which one am I iterating?

| | type | what `for i := 0; i < len(s); i++ { c := s[i] }` gives you | what `for i, r := range s` gives you |
|---|---|---|---|
| **byte** | `uint8` (0-255) | ✅ this — raw byte, one at a time | — |
| **rune** | `int32` (Unicode code point) | ❌ never | ✅ this — decodes UTF-8, may skip multiple bytes per iteration |

Indexing a string with `s[i]` **always** gives a byte, never a rune — even if the string contains multi-byte UTF-8 characters. That's exactly why a byte-range check (like `isValidHeaderValue` above) works for rejecting non-ASCII input: any multi-byte UTF-8 character has at least one byte `≥ 0x80`, so just walking bytes and checking the range catches it without ever needing to decode runes.

### Rune literals (`'A'`) vs byte values

`'A'` in Go is a **rune literal** — its type defaults to `rune` (`int32`), and it holds the Unicode code point value. `r == 'A'` compares code points, not bytes.

Since ASCII (0-127) is literally the first 128 Unicode code points by design, for ASCII characters the rune value and the byte value are numerically identical: `'A'` as a rune is `65`, and the byte `0x41` is also `65`. Same number, just a different *type* (`rune` is `int32`, a byte is `uint8`).

### Number bases — same value, different notation

A byte isn't inherently hex or octal — it's just 8 bits in memory: `01000001`. Decimal, hex, and octal are all just different human-readable ways to *write* that same number; none of them changes what's stored:

```
binary:  01000001
decimal: 65
hex:     0x41
octal:   0101 (or 0o101 in modern Go)
```

All four represent the exact same byte/code point. As long as the numbers match, they point to the same entry in the ASCII/Unicode table — the notation you pick (hex, decimal, octal) doesn't change the underlying value.

**Why hex is the convention for byte ranges**: 1 byte = 8 bits = exactly 2 hex digits (`0x00`-`0xFF`), a clean 1-to-1 mapping (4 bits per hex digit). That's why specs/code dealing with byte ranges (like the header validation range `0x21`-`0x7E`) use hex rather than decimal or octal.

**Octal (base 8)** is mostly a historical leftover from pre-8-bit-byte systems. Still seen today in:
- Unix file permissions: `chmod 755` (each digit = 3 bits = one `rwx` group)
- Old C/Go escape sequences: `'\101'` (octal) == `'\x41'` (hex) == `'A'`
- Go syntax: `0o17` (or legacy `017`) = octal 17 = decimal 15

---

# Reading RFCs (ABNF Notation)

RFCs (7230/9110/9112 for HTTP) define grammar using **ABNF** (Augmented Backus-Naur Form, RFC 5234).

| Symbol | Meaning | Example |
|---|---|---|
| `=` | "is defined as" | `rule = foo` |
| `/` | OR (alternatives) | `method = "GET" / "POST"` |
| `*` | repeat zero or more | `*DIGIT` = zero or more digits |
| `1*` | repeat one or more | `1*DIGIT` = at least one digit |
| `m*n` | repeat m to n times | `2*3DIGIT` = 2 to 3 digits |
| `[ ]` | optional (0 or 1) | `[ ":" port ]` |
| `( )` | grouping | `( "a" / "b" ) "c"` |
| `"text"` | literal string | `"GET"` |
| `%xHH` | a specific byte/char by hex | `%x0D` = CR |

### Example: request line

```
request-line = method SP request-target SP HTTP-version CRLF
```

Read left to right: a request-line is a method, then a space, then the target, then a space, then the version, then CRLF. No `*` or `/`, so each piece appears exactly once, in that order.

### Example: header field (with repetition)

```
field-line = field-name ":" OWS field-value OWS
```

`OWS` = optional whitespace, defined elsewhere as `OWS = *( SP / HTAB )` — zero or more spaces or tabs.

---

## Header Injection / Unicode Smuggling

**Background:** every character is just a number under the hood. **US-ASCII** is the original table — numbers 0-127, covers English letters/digits/punctuation. **Unicode** extends this with thousands more characters (other languages, emoji, lookalike letters) using numbers above 127, often multiple bytes per character.

**The vulnerability:** if a server trusts header values without checking their characters, an attacker can smuggle in control characters or lookalike Unicode characters. Different systems in the chain (your code, a proxy, a browser) may disagree about how to interpret those bytes — this class of bug enables CRLF injection and HTTP request smuggling.

**The fix:** only accept the "boring visible text" range — reject everything else.

| Range (hex) | Range (decimal) | Meaning | Allowed? |
|---|---|---|---|
| `0x00`-`0x1F` | 0-31 | control characters (NUL, CR, LF, etc.) | ❌ rejected |
| `0x20` | 32 | space | ✅ allowed (separator only) |
| `0x21`-`0x7E` | 33-126 | visible ASCII (letters, digits, punctuation) | ✅ allowed |
| `0x7F` | 127 | DEL (control char) | ❌ rejected |
| `0x80`+ | 128+ | any non-ASCII / Unicode byte | ❌ rejected |

`0x` means "hex digits follow" — `0x21` = 33 in decimal.

**Example bad request** (raw bytes, `\r\n` shown literally for clarity):
```
GET / HTTP/1.1\r\n
Host: example.com\r\n
X-Evil: foo\r\nSet-Cookie: admin=true\r\n
\r\n
```
If `X-Evil`'s value weren't checked for raw `\r\n`, the attacker's "header value" actually injects a brand new header (`Set-Cookie: admin=true`) that was never intended — this is CRLF/header injection. A byte-range check (reject anything outside `0x21`-`0x7E` plus space/tab) stops this, because `\r` (`0x0D`) and `\n` (`0x0A`) both fall in the rejected control-character range.

**Go implementation** (operate on bytes, not runes — any multi-byte UTF-8 character has a byte `≥ 0x80`, so this check rejects non-ASCII automatically without decoding Unicode):

```go
func isValidHeaderValue(s string) bool {
    for i := 0; i < len(s); i++ {
        c := s[i]
        if c == ' ' || c == '\t' {
            continue
        }
        if c < 0x21 || c > 0x7E {
            return false // control char, DEL, or non-ASCII byte
        }
    }
    return true
}
```

This also helps `Content-Length` parsing for free: `strconv.Atoi` only recognizes ASCII `'0'`-`'9'`, so it already rejects Unicode digit lookalikes (e.g. full-width `０-９`) — but validating at parse time closes the door earlier, for every header, not just `Content-Length`.

### Header keys are stricter than header values

The ASCII range check above is enough to stop injection, but RFC 7230 actually restricts header **keys** (`field-name`, grammar rule `token`) to a narrower whitelist than values — values can contain spaces and most punctuation, but keys can't.

Allowed in a header key (`tchar` = "token char"):

| Category | Allowed characters |
|---|---|
| Letters | `A-Z` `a-z` |
| Digits | `0-9` |
| Symbols | `! # $ % & ' * + - . ^ _ \` \| ~` |

Anything else — including `( ) < > @ , ; \ / [ ] ? = { } "` and any control/non-ASCII byte — is **not** a valid header key character, even though some of those (like `,` or `"`) are perfectly fine inside a header *value*.
