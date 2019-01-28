package nsocket

import (
	"log"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStartAndDial(t *testing.T) {
	name := "test"

	// Server
	ns, err := NewServer(name, func(server *Server, client *Client, message string) {
		err := server.Broadcast(message)
		if err != nil {
			log.Fatal("writeErr", err)
		}
	})
	if err != nil {
		log.Fatal(err)
	}
	go ns.Start()

	// Client
	client, err := Dial(name)
	count := 0
	if err != nil {
		log.Fatal(err)
	}

	// Write from client
	for i := 1; i <= 9; i++ {
		err = client.Write("test" + strconv.Itoa(i))
		if err != nil {
			log.Fatal("errClientWrite ", err)
		}
	}

	for {
		msg, err := client.Read()
		if err != nil {
			log.Fatal(err)
			break
		}
		// log.Println(msg)
		count++
		if msg == "test9" {
			client.Close()
			ns.Close()
			break
		}
	}
	require.Equal(t, 9, count)
}
