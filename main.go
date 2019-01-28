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

// Server nsocket server
type Server struct {
	Server    net.Listener
	clients   []*Client
	name      string
	silence   bool
	onMessage onMessageCallback
}

// Client of the nsocket server
type Client struct {
	rw   *bufio.ReadWriter
	conn net.Conn
}

type onMessageCallback func(server *Server, client *Client, message string)

// Write from client
func (client *Client) Write(msg string) error {
	_, err := client.rw.WriteString(msg + "\n")
	if err != nil {
		return err
	}
	return client.rw.Flush()
}

// Read client
func (client *Client) Read() (string, error) {
	buf, err := client.rw.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.Trim(buf, "\n"), nil
}

// Close client
func (client *Client) Close() {
	client.conn.Close()
}

// CloseClient removes client buffer from the list
func (ns *Server) CloseClient(client *Client) {
	clientIndex := -1
	for i := range ns.clients {
		if ns.clients[i] == client {
			clientIndex = i
			break
		}
	}
	if clientIndex != -1 {
		ns.clients = append(ns.clients[:clientIndex], ns.clients[clientIndex+1:]...)
	}
}

// Broadcast to all clients
func (ns *Server) Broadcast(msg string) error {
	var err error
	for _, v := range ns.clients {
		err = v.Write(msg)
		if err != nil {
			log.Println("broadcastErr: ", err)
			break
		}
	}
	return err
}

// Close the server
func (ns *Server) Close() error {
	if ns.Server != nil {
		err := ns.Server.Close()
		for len(ns.clients) > 0 {
			time.Sleep(100 * time.Millisecond)
		}
		return err
	}

	return nil
}

func (ns *Server) readClient(client *Client) {
	for {
		buf, err := client.Read()
		if err != nil {
			log.Println("readClientErr", err)
			break
		}
		ns.onMessage(ns, client, buf)
	}
	log.Println("closingClient")
	ns.CloseClient(client)
}

// Start a named socket, blocks by reading
func (ns *Server) Start() {
	log.Println("glad to serve")
	for {
		newConn, err := ns.Server.Accept()
		if err != nil {
			log.Println("listenErr", err)
			break
		}
		log.Println("newClient")
		newClient := &Client{
			conn: newConn,
			rw:   bufio.NewReadWriter(bufio.NewReader(newConn), bufio.NewWriter(newConn)),
		}
		ns.clients = append(ns.clients, newClient)
		// handshake message
		err = newClient.Write("handshake")
		if err != nil {
			log.Println("writeErr", err)
		} else {
			go ns.readClient(newClient)
		}
	}
	log.Println("shutdown")
	ns.clients = []*Client{}
}

// Dial to a named socket
func Dial(name string) (*Client, error) {
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

	rw := bufio.NewReadWriter(bufio.NewReader(client), bufio.NewWriter(client))
	newClient := Client{
		rw:   rw,
		conn: client,
	}
	msg, err := newClient.Read()
	if err != nil {
		return nil, err
	}
	if msg != "handshake" {
		return nil, errors.New("handshake message failed")
	}

	return &newClient, nil
}

// NewServer returns a server pointer
func NewServer(name string, onMessage onMessageCallback) (*Server, error) {
	var err error

	if name == "" {
		return nil, errors.New("the name of the socket server can't be empty")
	}

	// Default to echo broadcast
	if onMessage == nil {
		return nil, errors.New("onMessage fuction can't be empty")
	}

	ns := &Server{
		name:      name,
		onMessage: onMessage,
	}
	if runtime.GOOS == "windows" {
		ns.Server, err = Listen(`\\.\pipe\`+ns.name+`.sock`, nil)
		if err != nil {
			return nil, err
		}
	} else {
		os.RemoveAll("/tmp/" + ns.name + ".sock")
		ns.Server, err = net.Listen("unix", `/tmp/`+ns.name+`.sock`)
		if err != nil {
			return nil, err
		}
	}
	return ns, nil
}
