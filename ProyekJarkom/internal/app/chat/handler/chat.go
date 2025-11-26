// Package handler starts and handles incoming and outgoing messages
package handler

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"strings"

	"github.com/FhmiSddq/ProyekJarkom/internal/infra/env"
)

type ChatHandlerItf interface {
	RequestConnection(address string)
	Announce(message string)
	Listen()
	HandleConnection(connection net.Conn)
	CloseConnection(connection net.Conn)
	GetConnection(address string) net.Conn
}

type ChatHandler struct {
	env         *env.Env
	connections map[net.Conn]string
}

func New(env *env.Env) ChatHandlerItf {
	ChatHandler := ChatHandler{
		env:         env,
		connections: make(map[net.Conn]string),
	}

	return &ChatHandler
}

func (c *ChatHandler) RequestConnection(address string) {
	connection, err := net.Dial("tcp", address)
	if err != nil {
		log.Println(err)

		return
	}

	log.Printf("tcp connection to host %s success\n", address)

	c.connections[connection] = address

	go c.HandleConnection(connection)
}

func (c *ChatHandler) Announce(message string) {
	for connection := range c.connections {
		fmt.Fprintln(connection, message)
	}
}

// groupMembers maps group name to list of addresses of member peers
var groupMembers = make(map[string][]string)

// AssignPeerToGroup assigns a peer address to a group
func (c *ChatHandler) AssignPeerToGroup(group string, address string) {
	if _, exists := groupMembers[group]; !exists {
		groupMembers[group] = []string{}
	}
	// Check if address already in group
	for _, addr := range groupMembers[group] {
		if addr == address {
			return
		}
	}
	groupMembers[group] = append(groupMembers[group], address)
}

// BroadcastToGroup sends message to all peers in the group
func (c *ChatHandler) BroadcastToGroup(group string, message string) {
	members, exists := groupMembers[group]
	if !exists {
		return
	}
	for _, addr := range members {
		conn := c.GetConnection(addr)
		if conn != nil {
			fmt.Fprintln(conn, message)
		}
	}
}

func (c *ChatHandler) Listen() {
	listenAddress := fmt.Sprintf(":%s", strconv.Itoa(int(c.env.Port)))

	listener, err := net.Listen("tcp", listenAddress)
	if err != nil {
		log.Println(err)

		return
	}

	for {
		connection, err := listener.Accept()
		if err != nil {
			log.Println(err)

			continue
		}

		c.connections[connection] = listenAddress

		go c.HandleConnection(connection)
	}
}

func (c *ChatHandler) HandleConnection(connection net.Conn) {
	connectionName := connection.RemoteAddr().String()

	for {
		message, err := bufio.NewReader(connection).ReadString('\n')
		message = strings.TrimSpace(message)

		if err != nil {
			if err == io.EOF {
				log.Printf("disconnected from host: %s\n", connectionName)
			} else {
				log.Println(err)
			}

			break
		}

		log.Printf("[%s]: %s\n", connectionName, message)
	}

	c.CloseConnection(connection)
}

func (c *ChatHandler) CloseConnection(connection net.Conn) {
	connectionName := connection.RemoteAddr().String()

	log.Printf("stopped handling connection: %s\n", connectionName)

	delete(c.connections, connection)

	connection.Close()
}

func (c *ChatHandler) GetConnection(address string) net.Conn {
	for connection := range c.connections {
		if c.connections[connection] == address {
			return connection
		}
	}

	return nil
}
