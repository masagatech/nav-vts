package servers

import (
	"fmt"
	"io"
	"log"
	"net"
	"strconv"

	"github.com/masagatech/nav-vts/app/models"
)

// Server ...
type Server struct {
	host                     string
	port                     string
	onNewClientCallback      func(c *Client)
	onClientConnectionClosed func(c *Client, err error)
	onNewMessage             func(c *Client, message []byte)
}

// Client ...
type Client struct {
	conn   net.Conn
	Server *Server
}

func NewTCPServer(config *models.Config) *Server {
	s := &Server{
		host: config.Servers.Tcp_server.Host,
		port: strconv.Itoa(config.Servers.Tcp_server.Port),
	}
	s.OnNewClient(func(c *Client) {})
	s.OnNewMessage(func(c *Client, message []byte) {})
	s.OnClientConnectionClosed(func(c *Client, err error) {})
	return s
}

// Called right after Server starts listening new client
func (s *Server) OnNewClient(callback func(c *Client)) {
	s.onNewClientCallback = callback
}

// Called right after connection closed
func (s *Server) OnClientConnectionClosed(callback func(c *Client, err error)) {
	s.onClientConnectionClosed = callback
}

// Called when Client receives new message
func (s *Server) OnNewMessage(callback func(c *Client, message []byte)) {
	s.onNewMessage = callback
}

// Run ...
func (s *Server) Start() {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%s", s.host, s.port))
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()
	fmt.Println("Tcp Server started", fmt.Sprintf("%s:%s", s.host, s.port))
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}

		client := &Client{
			conn:   conn,
			Server: s,
		}
		fmt.Println("New Client")
		go client.handleRequest()
		s.onNewClientCallback(client)
	}
}

// Send text message to client
func (c *Client) Send(message string) error {
	_, err := c.conn.Write([]byte(message))
	return err
}

// Send bytes to client
func (c *Client) SendBytes(b []byte) error {
	_, err := c.conn.Write(b)
	return err
}

func (c *Client) Conn() net.Conn {
	return c.conn
}

func (c *Client) Close() error {
	return c.conn.Close()
}

func (c *Client) handleRequest() {
	for {
		buf := make([]byte, 1024)
		// Read the incoming connection into the buffer.
		n, err := c.conn.Read(buf)
		if err != nil {
			//if file/socket is closed remove the socket from list.
			if io.EOF == err {

				//this.removeClient(connection)
				c.conn.Close()
				c.Server.onClientConnectionClosed(c, err)
				break
			} else {
				//fmt.Println(err)
				//protocalHandler.RemoveClient(connection)
				continue
			}
		} else {
			if n > 0 {
				//data := string(buf[:n])
				if err != nil {
					fmt.Println("Error reading:", err.Error())
				}
				// code block to handle incoming data

				//protocalHandler.ParseData(buf, n, connection)
				//connection.Write([]byte(data + "\n"))
				c.Server.onNewMessage(c, buf)
			}
		}
		//fmt.Printf("Message incoming: %s", string(message))
		//client.conn.Write([]byte("Message received.\n"))
	}
}
