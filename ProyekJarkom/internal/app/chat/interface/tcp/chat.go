// Package tcp abstracts the implementation of chat handler
package tcp

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/FhmiSddq/ProyekJarkom/internal/app/chat/handler"
	"github.com/FhmiSddq/ProyekJarkom/internal/infra/env"
)

type Chat struct {
	Env         *env.Env
	ChatHandler handler.ChatHandlerItf
}

func NewChat(env *env.Env, chatHandler handler.ChatHandlerItf) {
	chat := Chat{
		Env:         env,
		ChatHandler: chatHandler,
	}

	chat.Start()
}

func (c *Chat) Start() {
	go c.ChatHandler.Listen()

	reader := bufio.NewReader(os.Stdin)
	running := true

	for running {
		message, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println(err)

			break
		}

		message = strings.TrimSpace(message)
		arguments := strings.Split(message, " ")

		switch arguments[0] {
		case c.Env.Connect:
			if len(arguments) != 2 {
				log.Printf("Usage: %s [HOST]:[PORT]\n", arguments[0])

				continue
			}

			c.ChatHandler.RequestConnection(arguments[1])

		case c.Env.Disconnect:
			if len(arguments) != 2 {
				log.Printf("Usage: %s [HOST]:[PORT]\n", arguments[0])

				continue
			}

			c.ChatHandler.CloseConnection(c.ChatHandler.GetConnection(arguments[1]))

		case c.Env.Exit:
			running = false

		default:
			c.ChatHandler.Announce(message)
		}
	}
}
