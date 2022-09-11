package app

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
	"zerotrust_chat/chat"
	"zerotrust_chat/cmd/builder"
	"zerotrust_chat/logger"
)

var _ chat.ReceiveHandler = &printer{}

type App struct {
	builder        builder.Builder
	sessionManager chat.SessionManager
	server         chat.Server
	printer        *printer
}

type printer struct {
	msgChan chan string
}

func newPrinter() *printer {
	return &printer{
		msgChan: make(chan string, 1),
	}
}

func (p *printer) OnReceive(msg string) {
	p.msgChan <- msg
}

func (p *printer) Run() {
	for {
		msg := <-p.msgChan
		fmt.Println(msg)
	}
}

func NewApp(builder builder.Builder) App {
	printer := newPrinter()
	return App{
		builder:        builder,
		sessionManager: builder.GetSessionManager(),
		server:         builder.NewServer(printer),
		printer:        printer,
	}
}

func (a *App) Connect(addr string) error {
	client, err := a.builder.NewClient(addr, a.printer)
	if err != nil {
		logger.Error(err)
		return err
	}
	if !a.sessionManager.Add(client) {
		return errors.New("client with added:" + client.GetId())
	}

	return nil
}

func (a *App) Disconnect(id string) {
	a.sessionManager.Remove(id)
}

func (a *App) Write(id string, msg string) error {
	return a.sessionManager.Write(id, []byte(msg))
}

func (a *App) Run() {
	go func() {
		a.server.Run()
	}()

	go func() {
		a.printer.Run()
	}()

	scanner := bufio.NewScanner(os.Stdin)

	currId := ""

	for scanner.Scan() {
		cmd := scanner.Text()

		if strings.HasPrefix(cmd, "$") {
			// connect to new client
			s := strings.Split(cmd, " ")
			if len(s) < 2 {
				fmt.Println("invalid command:", cmd)
				continue
			}

			currId = s[1]
			fmt.Println("connecting to:", currId)

			_, err := a.builder.NewClient(currId, a.printer)
			if err != nil {
				logger.Debug(err)
				continue
			}

		} else if strings.HasPrefix(cmd, ">") {
			// switch client
			newId := strings.Split(cmd, " ")
			if len(newId) < 2 {
				fmt.Println("invalid cmd:", cmd)
				continue
			}

			newCurrId := newId[1]
			if !a.sessionManager.ValidateId(newCurrId) {
				fmt.Println("invalid id to switch:", newCurrId)
				continue
			}

			currId = newCurrId

		} else if strings.HasPrefix(cmd, "?") {
			// display connected client
			fmt.Printf("active ids: %v\n", a.sessionManager.GetAllIds())
		} else {
			if currId == "" {
				fmt.Println("please connect to client before sending")
				continue
			} else {
				a.sessionManager.Write(currId, []byte(cmd))
			}
		}
	}
}
