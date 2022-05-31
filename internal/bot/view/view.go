package view

import (
	"fmt"
	"github.com/google/uuid"
	tgbotapi "github.com/mazanur/telegram-bot-api/v6"
	"market-bot/internal/service/model"
	"market-bot/sdk/tgbot"
	"strconv"
	"strings"
)

type RoomProvider interface {
}

type UserProvider interface {
	GetUser(userId int64) (tgbot.User, error)
}

type View struct {
	chatProv tgbot.ChatProvider
	userProv UserProvider

	tg *tgbot.Bot
}

func NewView(btnProv tgbot.ChatProvider, userProv UserProvider, tg *tgbot.Bot) *View {
	return &View{
		chatProv: btnProv,
		userProv: userProv,
		tg:       tg}
}

func (v *View) StartViewInPublic(u *tgbot.Update) (tgbotapi.Message, error) {
	msg := new(tgbot.MessageBuilder).
		NewMessage(u.GetChatId()).
		Text("❗️ Команда /start доступна только в приватном чате с ботом").
		Build()

	return logIfError(v.tg.Send(msg))
}
func (v *View) Pong(u *tgbot.Update) (tgbotapi.Message, error) {
	msg := new(tgbot.MessageBuilder).
		NewMessage(u.GetChatId()).
		Edit(u.IsButton()).
		Text("_pong_").
		Build()

	return logIfError(v.tg.Send(msg))
}

func (v *View) ErrorMessage(u *tgbot.Update, text string) (tgbotapi.Message, error) {
	c := &tgbotapi.CallbackConfig{
		CallbackQueryID: u.CallbackQuery.ID,
		Text:            text,
		ShowAlert:       true,
	}
	return logIfError(v.tg.Send(c))
}

func (v *View) WarnMessage(text string, u *tgbot.Update) (tgbotapi.Message, error) {
	c := &tgbotapi.CallbackConfig{
		CallbackQueryID: u.CallbackQuery.ID,
		Text:            text,
		ShowAlert:       false,
	}
	return logIfError(v.tg.Send(c))
}

func (v *View) ErrorMessageText(text string, u *tgbot.Update) (tgbotapi.Message, error) {
	msg := new(tgbot.MessageBuilder).
		Message(u.GetUserId(), u.GetMessageId()).
		InlineId(u.GetInlineId()).
		Edit(u.IsButton()).
		Text(text).
		Build()

	return logIfError(v.tg.Send(msg))
}

func (v *View) NewDeleteMessage(chatID int64, messageID int) (tgbotapi.Message, error) {
	c := tgbotapi.NewDeleteMessage(chatID, messageID)
	return logIfError(v.tg.Send(c))
}

func (v *View) SendChatWritingAction(chatId int64) (tgbotapi.Message, error) {
	msg := tgbotapi.NewChatAction(chatId, tgbotapi.ChatTyping)
	return logIfError(v.tg.Send(msg))
}

func (v *View) ShowAnswers(commands []model.SysCommand, u *tgbot.Update) (tgbotapi.Message, error) {
	inlineRequest := tgbot.NewInlineRequest(u.GetInlineId())
	for _, sys := range commands {
		inlineRequest.AddArticle(uuid.NewString(), sys.Question, "", fmt.Sprintf("*%v*\n\n%v", sys.Question, sys.Message))
	}
	return logIfError(v.tg.Send(inlineRequest.Build()))
}

func (v *View) ShowAuto(auto []string, u *tgbot.Update) (tgbotapi.Message, error) {
	inlineRequest := tgbot.NewInlineRequest(u.GetInlineId())
	for _, res := range auto {
		inlineRequest.AddArticle(uuid.NewString(), res, "", res)
	}
	return logIfError(v.tg.Send(inlineRequest.Build()))
}

func (v *View) ShowEmployees(employees []model.Employee, isAdmin bool, u *tgbot.Update) (tgbotapi.Message, error) {
	inlineRequest := tgbot.NewInlineRequest(u.GetInlineId())
	for _, res := range employees {
		text := fmt.Sprintf("Данные сотрудника\n%v\nтелефон: %v\nстатус: %v", employeeLink(res), *res.Phone, res.Status)
		inlineRequest.AddArticle(uuid.NewString(), fmt.Sprintf("%v %v", res.FirstName, *res.LastName), "", text)
		if isAdmin {
			inlineRequest.AddKeyboardRow().AddButton("Удалить сотрудника", "NO")
		}
	}
	return logIfError(v.tg.Send(inlineRequest.Build()))
}

func (v *View) ShowChats(chats []model.Chat, u *tgbot.Update) (tgbotapi.Message, error) {
	inlineRequest := tgbot.NewInlineRequest(u.GetInlineId())
	for _, chat := range chats {
		text := fmt.Sprintf("Данные чата *%v*\n\nописание: %v\nактиность: %v\nобязательность: %v\nотдел: %v", chat.Title, *chat.Description, chat.Active, chat.Required, chat.District)
		var title string
		if chat.Active {
			title = "\U0001F7E2 " + chat.Title
		} else {
			title = "🔴 " + chat.Title
		}
		inlineRequest.AddArticle(uuid.NewString(), title, "", text)

		changeChatActiveBtn := v.createButton(ActionChangeChatActive, map[string]string{"chatId": strconv.FormatInt(chat.TelegramId, 10)})
		changeChatRequiredBtn := v.createButton(ActionChangeChatRequired, map[string]string{"chatId": strconv.FormatInt(chat.TelegramId, 10)})

		inlineRequest.AddKeyboardRow().AddButton("Сменить статус", changeChatActiveBtn.Id)
		inlineRequest.AddKeyboardRow().AddButton("Сменить обязательность", changeChatRequiredBtn.Id)
	}
	return logIfError(v.tg.Send(inlineRequest.Build()))
}

func (v *View) ShowChat(chat model.Chat, u *tgbot.Update) (*tgbotapi.APIResponse, error) {
	text := fmt.Sprintf("Данные чата *%v*\n\nописание: %v\nактиность: %v\nобязательность: %v\nотдел: %v", chat.Title, *chat.Description, chat.Active, chat.Required, chat.District)
	msg := new(tgbot.MessageBuilder).
		InlineId(u.GetInlineId()).
		Edit(u.IsButton()).
		Text(text)

	changeChatActiveBtn := v.createButton(ActionChangeChatActive, map[string]string{"chatId": strconv.FormatInt(chat.TelegramId, 10)})
	changeChatRequiredBtn := v.createButton(ActionChangeChatRequired, map[string]string{"chatId": strconv.FormatInt(chat.TelegramId, 10)})

	msg.AddKeyboardRow().AddButton("Сменить статус", changeChatActiveBtn.Id).
		AddKeyboardRow().AddButton("Сменить обязательность", changeChatRequiredBtn.Id)

	return logIfErrorApiResponse(v.tg.Request(msg.Build()))
}

func (v *View) GetMe() tgbotapi.User {
	me, _ := v.tg.GetMe()
	return me
}

func employeeLink(e model.Employee) string {
	return fmt.Sprintf("[%s %v](tg://user?id=%d)", e.FirstName, *e.LastName, e.TelegramId)
}

func userLink(u tgbot.User) string {
	if u.UserId == nil {
		return fmt.Sprintf("@%v", escapeCharacters(u.UserName))
	}
	return fmt.Sprintf("[%v %v](tg://user?id=%v)", deleteCharacters(u.DisplayName), deleteCharacters(u.LastName), *u.UserId)
}

func escapeCharacters(string string) string {
	res := strings.ReplaceAll(string, "_", "\\_")
	res = strings.ReplaceAll(res, "*", "\\*")
	res = strings.ReplaceAll(res, "'", "\\'")
	res = strings.ReplaceAll(res, "`", "\\`")
	res = strings.ReplaceAll(res, "’", "\\’")
	return res
}

func deleteCharacters(string string) string {
	res := strings.ReplaceAll(string, "_", "")
	res = strings.ReplaceAll(res, "*", "")
	res = strings.ReplaceAll(res, "'", "")
	res = strings.ReplaceAll(res, "`", "")
	res = strings.ReplaceAll(res, "’", "")
	return res
}

func (v *View) AddNominationMessage(chatId int64, message string, user tgbot.User) (tgbotapi.Message, error) {

	text := fmt.Sprintf("%v - %v !", message, userLink(user))

	msg := new(tgbot.MessageBuilder).
		ChatId(chatId).
		Text(text).
		Build()

	return logIfError(v.tg.Send(msg))
}

func (v *View) ShowNominations(results []model.Nomination, u *tgbot.Update) (tgbotapi.Message, error) {

	text := fmt.Sprintf("*Номинации* %v \n\n", userLink(u.GetUser()))
	if len(results) == 0 {
		text += fmt.Sprintf("Видимо %v еще не был номинириван", userLink(u.GetUser()))
	} else {
		for _, res := range results {
			date := res.Date.Format("02-01-2006")
			text += fmt.Sprintf("%v %v\n", res.Message, date)
		}
	}

	msg := new(tgbot.MessageBuilder).
		ChatId(u.GetChatId()).
		Text(text).
		Build()

	return logIfError(v.tg.Send(msg))
}

func (v *View) SendMessage(u *tgbot.Update) {

	msg := new(tgbot.MessageBuilder).
		AddReplyKeyboardRow().
		ChatId(u.GetChatId()).
		Text("text").
		AddWebAppInfoButton("https://almazik77.github.io/tg-bot/", "button text").
		Build()

	_, _ = logIfError(v.tg.Send(msg))
}

func (v *View) SendMessage2(text string, u *tgbot.Update) {

	msg := new(tgbot.MessageBuilder).
		ChatId(u.GetChatId()).
		Text(text).
		Build()

	_, _ = logIfError(v.tg.Send(msg))
}
