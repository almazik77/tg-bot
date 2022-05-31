package bot_handler

import (
	"github.com/go-pkgz/lgr"
	tgbotapi "github.com/mazanur/telegram-bot-api/v6"
	"market-bot/internal/bot/view"
	"market-bot/internal/service"
	"market-bot/sdk/tgbot"
)

type Config struct {
	BotName                   string
	SuperUsers                []string
	BlackChatId               int64
	MainAnswererTelegramId    int64
	NominationDelayMaxSeconds int64
}

type TelegramApi interface {
	GetInviteLink(config tgbotapi.ChatInviteLinkConfig) (string, error)
	Send(c tgbotapi.Chattable) (tgbotapi.Message, error)
	Request(c tgbotapi.Chattable) (*tgbotapi.APIResponse, error)
}

type BotApp struct {
	view            *view.View
	employeeService *service.EmployeeService
	config          *Config

	api TelegramApi
}

func NewBotApp(view *view.View, employeeProv *service.EmployeeService, config *Config, api TelegramApi) *BotApp {
	return &BotApp{view: view,
		employeeService: employeeProv,
		config:          config,
		api:             api,
	}
}

func (b *BotApp) Handle(u *tgbot.Update) {

	if u.Message != nil && u.Message.WebAppData != nil {
		lgr.Printf("[INFO] WebAppData = %v", u.Message.WebAppData)
		b.view.SendMessage2(u.Message.WebAppData.Data, u)
		return
	}

	b.view.SendMessage(u)

}
