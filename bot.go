package telegram

import (
	"github.com/AmandaCameron/go-telegram/api"
)

// Bot represents and holds the meta information for a Telegram Bot.
type Bot struct {
	api *tgbotapi.BotAPI

	Help       string
	LastUpdate int32

	commands Commands
}

// Sendable means you can use this to send a message or file to a user.
type Sendable interface {
	Send() error
}

// Uploadable is like Sendable, but can return a string FileID as well.
type Uploadable interface {
	Upload() (string, error)
}

// NewBot creates a new bot with the specified token.
func NewBot(token string) (*Bot, error) {
	b, err := tgbotapi.NewBotAPI(token)

	if err != nil {
		return nil, err
	}

	bot := &Bot{
		api: b,

		commands: Commands{},
	}

	bot.AddCommand(helpCommand)

	return bot, nil
}

// AddCommand adds a new command to the bot's default command set.
func (bot *Bot) AddCommand(cmd Command) {
	bot.commands.Add(cmd)
}

// HandleCommand handles a command from the bot's default command set.
func (bot *Bot) HandleCommand(msg Message) bool {
	return bot.commands.Handle(msg)
}

// SendTyping sends a message saying that the bot is typing a message.
func (bot *Bot) SendTyping(chatId int32) {
	bot.api.SendChatAction(tgbotapi.NewChatAction(chatId, tgbotapi.ChatTyping))
}

// GetMessages returns the current messages from the Bot API.
func (bot *Bot) GetMessages() ([]Message, error) {
	updates, err := bot.api.GetUpdates(tgbotapi.UpdateConfig{
		Offset: bot.LastUpdate,
	})

	if err != nil {
		return nil, err
	}

	var ret []Message

	for _, update := range updates {
		if update.UpdateID <= bot.LastUpdate {
			continue
		}
		bot.LastUpdate = update.UpdateID

		ret = append(ret, Message{
			Message: update.Message,

			context: make(map[string]interface{}),
			bot:     bot,
			dir:     incoming,
		})
	}

	return ret, nil
}

// MessagesChan returns a channel that will recieve messages periodically from
// the bot's API endpoint.
func (bot *Bot) MessagesChan() (chan Message, error) {
	msgChan := make(chan Message, 100)

	go func() {
		for {
			msgs, err := bot.GetMessages()

			if err != nil {
				// TODO: Something

				return
			}

			for _, msg := range msgs {
				msgChan <- msg
			}
		}
	}()

	return msgChan, nil
}

var helpCommand = Command{
	Name:        "help",
	Description: "Shows this help page.",

	Handle: func(msg Message) {
		msg.ReplyWith("%s\n\n%s",
			msg.bot.Help,
			msg.bot.commands.Help(msg.IsGroup())).Send()
	},
}
