package main

import (
	"fmt"
	"os"
	"zerotrust_chat/cmd/app"
	"zerotrust_chat/cmd/builder"
	"zerotrust_chat/logger"
)

func main() {
	logger.SetLogLevel(logger.DEBUG)
	if len(os.Args) < 2 {
		logger.Fatal("need listener port: i.e. go run src/cmd/main.go <server port>")
	}

	serverAddr := fmt.Sprintf(":%s", os.Args[1])

	builder := builder.NewBuilder(serverAddr)

	app := app.NewApp(builder)
	app.Run()
}
