package env

import (
	"fmt"
	"strings"

	"github.com/tommy647/gramarr/internal/config"
	"github.com/tommy647/gramarr/internal/conversation"
	"github.com/tommy647/gramarr/internal/radarr"
	"github.com/tommy647/gramarr/internal/sonarr"
	"github.com/tommy647/gramarr/internal/users"
	"github.com/tommy647/gramarr/internal/util"
	tb "gopkg.in/tucnak/telebot.v2"
)

type Env struct {
	Config *config.Config
	Users  *users.UserDB
	Bot    *tb.Bot
	CM     *conversation.ConversationManager
	Radarr *radarr.Client
	Sonarr *sonarr.Client
}

func (e *Env) RequirePrivate(h func(m *tb.Message)) func(m *tb.Message) {
	return func(m *tb.Message) {
		if !m.Private() {
			return
		}
		h(m)
	}
}

func (e *Env) RequireAuth(access users.UserAccess, h func(m *tb.Message)) func(m *tb.Message) {
	return func(m *tb.Message) {
		user, _ := e.Users.User(m.Sender.ID)
		var msg []string

		// Is Revoked?
		if user.IsRevoked() {
			// Notify User
			msg = append(msg, "Your access has been revoked and you cannot reauthorize.")
			msg = append(msg, "Please reach out to the bot owner for support.")
			util.SendError(e.Bot, m.Sender, strings.Join(msg, "\n"))

			// Notify Admins
			msg = append(msg, fmt.Sprintf("Revoked users %s attempted the following command:", util.DisplayName(m.Sender)))
			msg = append(msg, fmt.Sprintf("`%s`", m.Text))
			util.SendAdmin(e.Bot, e.Users.Admins(), strings.Join(msg, "\n"))
			return
		}

		// Is Not Member?
		isAuthorized := user.IsAdmin() || user.IsMember()
		if !isAuthorized && access != users.UANone {
			// Notify User
			util.SendError(e.Bot, m.Sender, "You are not authorized to use this bot.\n`/auth [password]` to authorize.")

			// Notify Admins
			msg = append(msg, fmt.Sprintf("Unauthorized users %s attempted the following command:", util.DisplayName(m.Sender)))
			msg = append(msg, fmt.Sprintf("`%s`", m.Text))
			util.SendAdmin(e.Bot, e.Users.Admins(), strings.Join(msg, "\n"))
			return
		}

		// Is Non-Admin and requires Admin?
		if !user.IsAdmin() && access == users.UAAdmin {
			// Notify User
			util.SendError(e.Bot, m.Sender, "Only admins can use this command.")

			// Notify Admins
			msg = append(msg, fmt.Sprintf("User %s attempted the following admin command:", util.DisplayName(m.Sender)))
			msg = append(msg, fmt.Sprintf("`%s`", m.Text))
			util.SendAdmin(e.Bot, e.Users.Admins(), strings.Join(msg, "\n"))
			return
		}

		h(m)
	}
}
