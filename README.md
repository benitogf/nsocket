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
	name := "test"
	ns, err := nsocket.NewServer(
		name,
		func(server *nsocket.Server, client *nsocket.Client, message string) {
			log.Println("server:", message)
			// err := client.Write(message)
			err := server.Broadcast(message)
			if err != nil {
				log.Println("broadcastErr: ", err)
			}
		},
	)
	if err != nil {
		log.Println(err)
	}
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
	name := "test"
	client, err := nsocket.Dial(name)
	if err != nil {
		log.Println(err)
	}

	err = client.Write("test")
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
		if message == "test" {
			client.Close()
			break
		}
	}
}
```