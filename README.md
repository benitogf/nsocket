# nsocket

A fork of (go-winio)[https://github.com/Microsoft/go-winio] that aims to allow windows or unix named sockets.

## how to

```golang
package main

import (
	"bufio"
	"log"
	"strings"
	"time"

	"github.com/benitogf/nsocket"
)

func main() {
	name := "test"
	go nsocket.Start(name)
	time.Sleep(1 * time.Second) // wait for it
	c, err := nsocket.Dial(name)
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()
	rw := bufio.NewReadWriter(bufio.NewReader(c), bufio.NewWriter(c))
	err = nsocket.Write(rw, "init")
	if err != nil {
		log.Fatal(err)
	}
	for {
		buf, err := nsocket.Read(rw)
		if err != nil {
			log.Println(err)
			break
		}
		buf = strings.Trim(buf, "\n")
		log.Println("client", buf)
	}
}
```