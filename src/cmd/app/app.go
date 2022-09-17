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

var _ chat.ReceiveHandler = &App{}

type App struct {
	builder        builder.Builder
	sessionManager chat.SessionManager
	server         chat.Server
	printer        chan string
}

func (a *App) OnReceive(messages []chat.ChatMessage) {
	if messages == nil {
		logger.Fatal("chat messages is not supposed to be nil, shld be handled before calling this method")
	}
	for _, msg := range messages {
		a.printer <- msg.Data
	}
}

func NewApp(builder builder.Builder) *App {
	app := &App{
		builder:        builder,
		sessionManager: builder.GetSessionManager(),
		printer:        make(chan string, 100),
	}
	app.server = builder.NewServer(app)
	return app
}

func (a *App) Connect(addr string) error {
	client, err := a.builder.NewClient(addr, a)
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
		for {
			msg := <-a.printer
			fmt.Println(msg)
		}
	}()

	scanner := bufio.NewScanner(os.Stdin)

	currId := ""

	for scanner.Scan() {
		cmd := scanner.Text()

		if strings.HasPrefix(cmd, "$") {
			// connect to new client
			s := strings.Split(cmd, " ")
			if len(s) < 2 {
				a.println("invalid command:", cmd)
				continue
			}

			currId = s[1]
			a.println("connecting to:", currId)

			_, err := a.builder.NewClient(currId, a)
			if err != nil {
				logger.Debug(err)
				continue
			}

		} else if strings.HasPrefix(cmd, ">") {
			// switch client
			newId := strings.Split(cmd, " ")
			if len(newId) < 2 {
				a.println("invalid cmd:", cmd)
				continue
			}

			newCurrId := newId[1]
			if !a.sessionManager.ValidateId(newCurrId) {
				a.println("invalid id to switch:", newCurrId)
				continue
			}

			currId = newCurrId

		} else if strings.HasPrefix(cmd, "?") {
			// display connected client
			a.printf("active ids: %v\n", a.sessionManager.GetAllIds())
		} else {
			if currId == "" {
				a.println("please connect to client before sending")
				continue
			} else {
				a.sessionManager.Write(currId, []byte(cmd))
			}
		}
	}
}

func (a *App) printf(msg string, args ...any) {
	a.printer <- fmt.Sprintf(msg, args...)
}

func (a *App) println(msg string, args ...any) {
	a.printer <- fmt.Sprint(msg, args)
}
