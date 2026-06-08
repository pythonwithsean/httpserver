# HTTP Server Roadmap

## Step 1: Define a Request struct
Create a struct to hold: Method, Path, Version, Headers (map), Body

## Step 2: Build the parser function
Takes raw bytes `[]byte`, returns a `Request` struct:
- Find `\r\n\r\n` to split headers from body
- Split headers section by `\r\n` into lines
- First line: split by space → method, path, version
- Remaining lines: split by first `: ` → populate headers map
- Body = everything after `\r\n\r\n`

## Step 3: Handle incomplete reads
Your current single `conn.Read()` might not get the full request. You need to:
- Read in a loop
- Check if you've found `\r\n\r\n`
- If `Content-Length` header exists, keep reading until you have that many body bytes

## Step 4: Build a Response writer
A function that takes a `conn`, status code, headers, and body, then writes a properly formatted HTTP response:
```
HTTP/1.1 200 OK\r\n
Content-Length: 13\r\n
Content-Type: text/plain\r\n
\r\n
Hello, World!
```

## Step 5: Build a router
A map of `path → handler function`. Parse the request path, look it up in the map, call the matching handler or return 404.

## Step 6: Wire it together
In `handleConn`: parse request → route → call handler → write response

Start with Step 1 and 2, get a simple `GET /` returning "Hello World" working before adding complexity.
