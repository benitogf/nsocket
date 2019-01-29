# nsocket

[![Build Status][build-image]][build-url]


[build-url]: https://travis-ci.org/benitogf/nsocket
[build-image]: https://api.travis-ci.org/benitogf/nsocket.svg?branch=master&style=flat-square

A fork of [go-winio](https://github.com/Microsoft/go-winio), windows/unix named socket server and client utility

## how to
	server
```golang
package main

import (
	"log"

	"github.com/benitogf/nsocket"
)

func main() {
	ns, err := NewServer("test")
	if err != nil {
		log.Fatal(err)
	}
	go func() {
		for {
			select {
			case msg := <-ns.onMessage:
				msg.client.Write("_test" + strconv.Itoa(count))
				ns.Broadcast(msg.data, msg.client.path)
			}
		}
	}()
	ns.Start()
}
```
	client
```golang
package main

import (
	"log"

	"github.com/benitogf/nsocket"
)

func main() {
	client, err := nsocket.Dial("test", "one/two/three")
	if err != nil {
		log.Println(err)
	}

	err = client.Write("9")
	if err != nil {
		log.Println("errWrite: ", err)
	}

	for {
		message, err := client.Read()
		if err != nil {
			log.Println(err)
			break
		}
		log.Println("client:", message)
		if message == "9" {
			client.Close()
			break
		}
	}
}
```