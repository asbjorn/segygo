# SEG-Y reader in Go

Go-lang library to read and write SEG-Y files


This library tries to replicate the functionality of most other SEG-Y libraries out there. But, for Go-lang that is :-)
It currently supports reading a SEG-Y files (binary header; 3200 bytes) and extract the header info stored there.

The write support is not yet there (started, but not ready yet).

## Features

- [x] Reading and parsing of SEG-Y files
- [x] Writing of an empty SEG-Y file
- [ ] Writing traces to a SEG-Y file


## Usage

```go
package main
 
import (
    "github.com/asbjorn/segygo"
    "log"
)

func main() {
    log.Println("Testing segygo package")
    s, err := segygo.OpenFile("/tmp/myfile.segy")
    if err != nil {
        log.Fatalf("Unable to open SEG-Y file: %s", s.Filename)
    }   
}
```

# Credits

- Asbjørn A. Fellinghaug (asbjorn@fellinghaug.com)
