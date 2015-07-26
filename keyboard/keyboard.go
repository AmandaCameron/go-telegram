package keyboard

import (	
	"github.com/AmandaCameron/go-telegram/api"
)

type Modifier func(*tgbotapi.ReplyKeyboardMarkup)

func List(items ...string) Modifier {
	return func(kbd *tgbotapi.ReplyKeyboardMarkup) {
		for _, row := range items {
			Row(row)(kbd)
		}
	}
}

func Row(row ...string) Modifier {
	return func(kbd *tgbotapi.ReplyKeyboardMarkup) {
		kbd.Keyboard = append(kbd.Keyboard, row)
	}
}

func Once(kbd *tgbotapi.ReplyKeyboardMarkup) {
	kbd.OneTimeKeyboard = true
}

func Selective(kbd *tgbotapi.ReplyKeyboardMarkup) {
	kbd.Selective = true
}
