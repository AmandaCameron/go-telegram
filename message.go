package telegram

import (
	"fmt"

	"github.com/AmandaCameron/go-telegram/api"
	"github.com/AmandaCameron/go-telegram/keyboard"
)

type direction uint8

const (
	outgoing direction = iota
	incoming
)

// Message is our encapsulation of the tgbotapi.Message message type.
type Message struct {
	tgbotapi.Message

	context               map[string]interface{}
	replyID               int32
	replyMarkup           interface{}
	disableWebPagePreview bool

	bot *Bot
	dir direction
}

// IsChat returns true if the message is a human-saying-stuff message.
func (msg Message) IsChat() bool {
	return !(msg.Message.DeleteChatPhoto ||
		msg.Message.GroupChatCreated ||
		len(msg.Message.NewChatPhoto) > 0 ||
		len(msg.Message.Photo) > 0 ||
		msg.Message.ReplyToMessage != nil ||
		msg.Document.FileID != "" ||
		msg.Sticker.FileID != "" ||
		msg.Contact.UserID != "" ||
		msg.NewChatParticipant.UserName != "" ||
		msg.LeftChatParticipant.UserName != "" ||
		msg.NewChatTitle != "" ||
		(msg.Message.Location.Latitude != 0.0 && msg.Message.Location.Longitude != 0.0))
}

// SetContext attaches a specified user-defined value to the message.
func (msg *Message) SetContext(name string, to interface{}) {
	if msg.context == nil {
		msg.context = make(map[string]interface{})
	}

	msg.context[name] = to
}

// GetContext returns a specified user-defined value attached to the
// message.
func (msg *Message) GetContext(name string) interface{} {
	if msg.context == nil {
		msg.context = make(map[string]interface{})
	}

	return msg.context[name]
}

// IsGroup returns true if this message took place inside a group chat.
func (msg Message) IsGroup() bool {
	return (msg.Message.Chat.UserName == "")
}

// Message creates a new outbound message to the specified chatID and with
// the specified printf-formatted body
func (bot *Bot) Message(chatID int32, f string, args ...interface{}) *Message {
	return &Message{
		Message: tgbotapi.Message{
			Text: fmt.Sprintf(f, args...),

			Chat: tgbotapi.UserOrGroupChat{
				ID: chatID,
			},
		},

		context: make(map[string]interface{}),
		bot:     bot,
		dir:     outgoing,
	}
}

// ReplyWith generates an outbound reply message, with the formatted string
// from `f` and it's arguments.
func (msg Message) ReplyWith(f string, args ...interface{}) *Message {
	if msg.dir != incoming {
		return nil
	}

	return &Message{
		Message: tgbotapi.Message{
			Text: fmt.Sprintf(f, args...),

			Chat: tgbotapi.UserOrGroupChat{
				ID: msg.Chat.ID,
			},
		},

		context: msg.context,
		replyID: msg.MessageID,

		bot: msg.bot,
		dir: outgoing,
	}
}

// HideKeyboard tells Telegram to hide any existing Custom Keyboards.
// if `selective` is set, it will only go to one user.
func (msg *Message) HideKeyboard(selective bool) Sendable {
	if msg.dir != outgoing {
		return msg
	}

	msg.replyMarkup = tgbotapi.ReplyKeyboardHide{
		HideKeyboard: true,

		Selective: selective,
	}

	return msg
}

// CustomKeyboard tells Telegram of a special keyboard to present to the user when
// responding to this message.
func (msg *Message) CustomKeyboard(mods ...keyboard.Modifier) Sendable {
	if msg.dir != outgoing {
		return msg
	}

	markup := tgbotapi.ReplyKeyboardMarkup{}

	for _, mod := range mods {
		mod(&markup)
	}

	msg.replyMarkup = markup

	return msg
}

// ForceReply tells Telegram that this message, if responded to, should be forced into a

// reply message.
func (msg *Message) ForceReply(selective bool) Sendable {
	msg.replyMarkup = tgbotapi.ForceReply{}

	return msg
}

func (msg *Message) Send() error {
	_, err := msg.bot.api.SendMessage(
		tgbotapi.MessageConfig{
			ChatID: msg.Chat.ID,
			Text:   msg.Text,

			ReplyMarkup: msg.replyMarkup,

			ReplyToMessageID:      msg.replyID,
			DisableWebPagePreview: msg.disableWebPagePreview,
		})

	return err
}
