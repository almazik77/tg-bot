package view

import (
	"fmt"
	tgbotapi "github.com/mazanur/telegram-bot-api/v6"
	"market-bot/sdk/tgbot"
	"strings"
)

func (v *View) StartAdminView(u *tgbot.Update) (tgbotapi.Message, error) {

	addUsrBtn := v.createButton(ActionAddUser, nil)
	delUsrBtn := v.createButton(ActionDelUser, nil)
	newChatBtn := v.createButton(ActionNewChat, nil)

	msg := new(tgbot.MessageBuilder).
		Message(u.GetChatId(), u.GetMessageId()).
		Edit(u.IsButton()).
		Text("*Главный экран*").
		AddKeyboardRow().AddButton("➕ Добавить сотрудника", addUsrBtn.Id).
		AddKeyboardRow().AddButton("🗑 Удалить сотрудника", delUsrBtn.Id).
		AddKeyboardRow().AddButton("💬 Добавить чат", newChatBtn.Id).
		Build()

	return logIfError(v.tg.Send(msg))
}

func (v *View) AddHelpAdminMessage(u *tgbot.Update) (tgbotapi.Message, error) {
	bot := v.GetMe()
	botNameWithEscapes := strings.ReplaceAll(bot.UserName, "_", "\\_")
	text := "*Помощь*\n\n"
	text += "Для начала работы введите /start\n"
	text += "Для проверки работоспособности /ping\n"
	text += "Для просмотра своих номинаций, которые были в чёрном чатике /my\\_nominations \n\n"
	text += fmt.Sprintf("Для поиска FAQ @%v начинай вводить свой вопрос\n", botNameWithEscapes)
	text += fmt.Sprintf("Для поиска машин @%v *auto*:номер\\_машины\n", botNameWithEscapes)
	text += fmt.Sprintf("Для поиска сотрудника @%v *user*:имя сотрудника\n", botNameWithEscapes)
	text += fmt.Sprintf("Для поиска чатов @%v *chat*:название чата\n", botNameWithEscapes)
	text += fmt.Sprintf("Для поиска вопросов @%v *question*:название вопроса\n", botNameWithEscapes)

	msg := new(tgbot.MessageBuilder).
		NewMessage(u.GetChatId()).
		Text(text).
		AddKeyboardRow().AddButtonSwitchForCurrentChat("Поиск", "").
		Build()

	return logIfError(v.tg.Send(msg))
}
