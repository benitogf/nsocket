package nsocket

import (
	"bufio"
	"log"
	"net"
	"runtime"
	"strings"
)

var pool = []*bufio.ReadWriter{}

// Read a string from a named socket buffer
func Read(rw *bufio.ReadWriter) (string, error) {
	buf, err := rw.ReadString('\n')
	if err != nil {
		poolIndex := -1
		for i := range pool {
			if pool[i] == rw {
				poolIndex = i
				break
			}
		}
		pool = append(pool[:poolIndex], pool[poolIndex+1:]...)
	}
	return strings.Trim(buf, "\n"), err
}

// Write a string to a named socket buffer
func Write(rw *bufio.ReadWriter, msg string) error {
	_, err := rw.WriteString(msg + "\n")
	if err != nil {
		return err
	}
	return rw.Flush()
}

// Start a named socket, blocks by reading
func Start(name string) {
	var l net.Listener
	var err error
	if runtime.GOOS == "windows" {
		l, err = Listen(`\\.\pipe\`+name+`.sock`, nil)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		l, err = net.Listen("unix", `/tmp/`+name+`.sock`)
		if err != nil {
			log.Fatal(err)
		}
	}
	defer l.Close()
	log.Println("server started")
	for {
		c, err := l.Accept()
		log.Println("new client")
		if err != nil {
			log.Println(err)
			break
		}
		rw := bufio.NewReadWriter(bufio.NewReader(c), bufio.NewWriter(c))
		pool = append(pool, rw)
		go func() {
			for {
				buf, err := Read(rw)
				if err != nil {
					log.Println(err)
					break
				}
				buf = strings.Trim(buf, "\n")
				log.Println(buf)
				if buf == "init" {
					err = Write(rw, buf)
					if err != nil {
						log.Println(err)
						break
					}
				} else {
					for _, v := range pool {
						if v != rw {
							err = Write(v, buf)
							if err != nil {
								log.Println(err)
								break
							}
						}
					}
				}
			}
			c.Close()
		}()
	}
}

// Dial to a named socket
func Dial(name string) (net.Conn, error) {
	if runtime.GOOS == "windows" {
		return DialWindows(`\\.\pipe\`+name+`.sock`, nil)
	}

	return net.Dial("unix", `/tmp/`+name+`.sock`)
}
