
# Reading a File in Fixed-Size Chunks

## Overview

- Stream a file without loading it all into memory by reusing a fix-size buffer. 

## Key Concepts
1. Buffer reuse: the same 8-bytes slice is overwritten each iteration.
2. File offset: OS tracks the current position; each `Read` advances by `n`.
3. Short reads: `n` maybe less than `len(buffer)`, especially near EOF.

```go
f, err := os.Open("messages.txt")

if err != nil { return }
defer f.Close()

buf := make([]byte, 8)
for {
    n, err := f.Read(buf)
    if err == io.EOF { break }
    if err != nil { return }
    fmt.Printf("read: %s\n", buf[:n])
}
```

## Why Close Files Promptly
- Releases limited OS file descriptors(a small integer the OS uses to refer to an open I/O resource); avoids "too many open files"
- Prevents resource leaks in long-running processes. 
- For writers, ensure buffered data is flushed and durable. 
- `defer f.Close()` guarantees cleanup even on early returns/ erros. 