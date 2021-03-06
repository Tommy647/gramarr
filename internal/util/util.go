package util

import (
	"fmt"
	"strings"

	"github.com/tommy647/gramarr/internal/users"
	tb "gopkg.in/tucnak/telebot.v2"
)

func Send(bot *tb.Bot, to tb.Recipient, msg string) {
	bot.Send(to, msg, tb.ModeMarkdown)
}

func SendError(bot *tb.Bot, to tb.Recipient, msg string) {
	bot.Send(to, msg, tb.ModeMarkdown)
}

func SendAdmin(bot *tb.Bot, to []users.User, msg string) {
	SendMany(bot, to, fmt.Sprintf("*\\[Admin\\]* %s", msg))
}

func SendKeyboardList(bot *tb.Bot, to tb.Recipient, msg string, list []string) error {
	var buttons []tb.ReplyButton
	for _, item := range list {
		buttons = append(buttons, tb.ReplyButton{Text: item})
	}

	var replyKeys [][]tb.ReplyButton
	for _, b := range buttons {
		replyKeys = append(replyKeys, []tb.ReplyButton{b})
	}

	_, err := bot.Send(to, msg, &tb.ReplyMarkup{
		ReplyKeyboard:   replyKeys,
		OneTimeKeyboard: true,
	})

	return err
}

func SendMany(bot *tb.Bot, to []users.User, msg string) {
	for _, user := range to {
		bot.Send(user, msg, tb.ModeMarkdown)
	}
}

func DisplayName(u *tb.User) string {
	if u.FirstName != "" && u.LastName != "" {
		return EscapeMarkdown(fmt.Sprintf("%s %s", u.FirstName, u.LastName))
	}

	return EscapeMarkdown(u.FirstName)
}

func EscapeMarkdown(s string) string {
	s = strings.Replace(s, "[", "\\[", -1)
	s = strings.Replace(s, "]", "\\]", -1)
	s = strings.Replace(s, "_", "\\_", -1)
	return s
}
