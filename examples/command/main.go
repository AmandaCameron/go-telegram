package main

import (
	"os"

	"github.com/AmandaCameron/go-telegram"
)

func main() {
	bot, err := telegram.NewBot(os.Args[1])
	if err != nil {
		panic(err)
	}

	bot.Help = "Hello Bacon"

	bot.AddCommand(telegram.Command{
		Name:        "test",
		Description: "Test command to test and commands to test.",

		Handle: func(msg telegram.Message) {
			msg.ReplyWith("Bugger off!").Send()
		},
	})

	updates, err := bot.MessagesChan()
	if err != nil {
		panic(err)
	}

	for {
		select {
		case msg := <-updates:
			bot.HandleCommand(msg)
		}
	}
}
