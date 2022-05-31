package bot_handler

import (
	"encoding/json"
	tgbotapi "github.com/mazanur/telegram-bot-api/v6"
)

func (b *BotApp) GetInviteLinkByChatId(chatID int64) (string, error) {
	config := tgbotapi.CreateChatInviteLinkConfig{
		ChatConfig: tgbotapi.ChatConfig{ChatID: chatID}, MemberLimit: 1}
	resp, err := b.api.Request(config)
	if err != nil {
		return "", err
	}
	var resMap map[string]interface{}
	err = json.Unmarshal(resp.Result, &resMap)
	if err != nil {
		return "", err
	}
	return resMap["invite_link"].(string), nil
}

func (b *BotApp) DeleteFromChatConfig(userId int64, groupId int64) *tgbotapi.KickChatMemberConfig {
	return &tgbotapi.KickChatMemberConfig{UntilDate: 0,
		RevokeMessages: false,
		ChatMemberConfig: tgbotapi.ChatMemberConfig{
			ChatID: groupId,
			UserID: userId,
		}}
}
func (b *BotApp) MessageWritingStatus(chatId int64) error {
	_, err := b.api.Send(tgbotapi.NewChatAction(chatId, "typing"))
	return err
}

func (b *BotApp) SendInlineArticlesConfig(id string, articles []interface{}) error {
	inlineConfig := tgbotapi.InlineConfig{
		InlineQueryID: id,
		IsPersonal:    true,
		CacheTime:     0,
		Results:       articles,
	}
	_, err := b.api.Request(inlineConfig)
	return err
}
