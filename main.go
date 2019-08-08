package nsocket

import (
	"bufio"
	"errors"
	"log"
	"net"
	"os"
	"runtime"
	"strings"
	"time"
)

// Message data and origin
type Message struct {
	Client *Client
	Data   string
}

// Server nsocket server
type Server struct {
	Server    net.Listener
	Clients   []*Client
	Name      string
	Silence   bool
	OnMessage chan Message
}

// Client of the nsocket server
type Client struct {
	Buf  *bufio.ReadWriter
	Conn net.Conn
	Path string
}

// Write from client
func (client *Client) Write(msg string) error {
	_, err := client.Buf.WriteString(msg + "\n")
	if err != nil {
		return err
	}
	return client.Buf.Flush()
}

// Read client
func (client *Client) Read() (string, error) {
	buf, err := client.Buf.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.Trim(buf, "\n"), nil
}

// Close client
func (client *Client) Close() {
	client.Conn.Close()
}

// CloseClient removes client buffer from the list
func (ns *Server) CloseClient(client *Client) {
	clientIndex := -1
	for i := range ns.Clients {
		if ns.Clients[i] == client {
			clientIndex = i
			break
		}
	}
	if clientIndex != -1 {
		ns.Clients = append(ns.Clients[:clientIndex], ns.Clients[clientIndex+1:]...)
	}
}

// Broadcast to all clients
func (ns *Server) Broadcast(msg string, path string) {
	for _, v := range ns.Clients {
		if v.Path == path {
			err := v.Write(msg)
			if err != nil {
				log.Println("broadcastErr: ", err)
			}
		}
	}
}

// Close the server
func (ns *Server) Close() error {
	if ns.Server != nil {
		err := ns.Server.Close()
		for len(ns.Clients) > 0 {
			time.Sleep(100 * time.Millisecond)
		}
		return err
	}

	return nil
}

func (ns *Server) readClient(client *Client) {
	msg, err := client.Read()
	for err == nil {
		ns.OnMessage <- Message{
			Client: client,
			Data:   msg,
		}
		msg, err = client.Read()
	}
	log.Println("closingClient", err)
	ns.CloseClient(client)
}

// Start a named socket, blocks by reading
func (ns *Server) Start() {
	for {
		newConn, err := ns.Server.Accept()
		if err != nil {
			log.Println("listenErr", err)
			break
		}
		log.Println("newClient")
		newClient := &Client{
			Conn: newConn,
			Buf:  bufio.NewReadWriter(bufio.NewReader(newConn), bufio.NewWriter(newConn)),
		}
		ns.Clients = append(ns.Clients, newClient)
		// handshake message
		msg, err := newClient.Read()
		if err != nil {
			log.Fatal(errors.New("handshake message failed"))
		}
		newClient.Path = msg
		log.Println("path: ", msg)
		go ns.readClient(newClient)
	}
	log.Println("shutdown")
	ns.Clients = []*Client{}
}

// Dial to a named socket
func Dial(name string, path string) (*Client, error) {
	var err error
	var client net.Conn
	if runtime.GOOS == "windows" {
		client, err = DialWindows(`\\.\pipe\`+name+`.sock`, nil)
	} else {
		client, err = net.Dial("unix", `/tmp/`+name+`.sock`)
	}
	if err != nil {
		return nil, err
	}

	newClient := Client{
		Buf:  bufio.NewReadWriter(bufio.NewReader(client), bufio.NewWriter(client)),
		Conn: client,
		Path: path,
	}
	err = newClient.Write(path)
	if err != nil {
		return nil, errors.New("handshake message failed")
	}

	return &newClient, nil
}

// NewServer returns a server pointer
func NewServer(name string) (*Server, error) {
	var err error

	if name == "" {
		return nil, errors.New("the name of the socket server can't be empty")
	}

	ns := &Server{
		Name:      name,
		OnMessage: make(chan Message, 1),
	}
	if runtime.GOOS == "windows" {
		ns.Server, err = Listen(`\\.\pipe\`+ns.Name+`.sock`, nil)
		if err != nil {
			return nil, err
		}
	} else {
		os.RemoveAll("/tmp/" + ns.Name + ".sock")
		ns.Server, err = net.Listen("unix", `/tmp/`+ns.Name+`.sock`)
		if err != nil {
			return nil, err
		}
	}
	return ns, nil
}
