package nsocket

import (
	"log"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStartAndDial(t *testing.T) {
	name := "test"
	count := 0

	// Server
	ns, err := NewServer(name)
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
	go ns.Start()

	// Client
	client, err := Dial(name, "one/two/three")
	if err != nil {
		log.Fatal(err)
	}
	go func() {
		// Write from client
		for i := 1; i <= 9; i++ {
			err = client.Write("test" + strconv.Itoa(i))
			if err != nil {
				log.Fatal("errClientWrite ", err)
			}
		}
	}()
	msg, err := client.Read()
	for err == nil {
		log.Println(msg)
		count++
		if msg == "test9" {
			client.Close()
			ns.Close()
			break
		}
		msg, err = client.Read()
	}
	require.Equal(t, 18, count)
}
