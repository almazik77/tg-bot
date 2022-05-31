package bot_handler

import (
	"github.com/go-pkgz/lgr"
	tgbotapi "github.com/mazanur/telegram-bot-api/v6"
)

func (b *BotApp) SetMyCommands() {

	var commands []tgbotapi.BotCommand
	commands = append(commands, tgbotapi.BotCommand{Command: "start", Description: "Главный экран"})
	commands = append(commands, tgbotapi.BotCommand{Command: "my_nominations", Description: "Мои номинации"})
	commands = append(commands, tgbotapi.BotCommand{Command: "ping", Description: "Проверка работоспособности"})
	commands = append(commands, tgbotapi.BotCommand{Command: "help", Description: "Помощь"})
	config := tgbotapi.SetMyCommandsConfig{Commands: commands}
	_, err := b.api.Request(config)
	if err != nil {
		lgr.Printf("[ERROR] cant set commands %v", err)
		return
	}
}
