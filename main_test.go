package nsocket

import (
	"bufio"
	"log"
	"os"
	"strings"
	"testing"
	"time"
)

func TestStartAndDial(t *testing.T) {
	name := "test"
	buf := ""
	var err error
	os.RemoveAll("/tmp/" + name + ".sock")
	go Start(name)
	time.Sleep(1 * time.Second) // wait for it
	c, err := Dial(name)
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()
	rw := bufio.NewReadWriter(bufio.NewReader(c), bufio.NewWriter(c))
	err = Write(rw, "init")
	if err != nil {
		log.Fatal(err)
	}
	for {
		buf, err = Read(rw)
		if err != nil {
			log.Println(err)
			break
		}
		buf = strings.Trim(buf, "\n")
		log.Println("client", buf)
		break
	}
}
