package view

import (
	"encoding/json"
	"github.com/go-pkgz/lgr"
	"github.com/google/uuid"
	tgbotapi "github.com/mazanur/telegram-bot-api/v6"
	"market-bot/sdk/tgbot"
)

func (v *View) createButton(action tgbot.Action, data map[string]string) *tgbot.Button {
	id := uuid.New()
	button := tgbot.Button{
		Id:     id.String(),
		Action: action,
		Data:   data,
	}
	err := v.chatProv.SaveButton(button)
	if err != nil {
		lgr.Printf("[ERROR] cannot save button, %v", err)
	}
	return &button
}

func logIfError(send tgbotapi.Message, err error) (tgbotapi.Message, error) {
	if err == nil {
		return send, nil
	}
	switch err.(type) {
	default:
		lgr.Printf("[ERROR] cannot send, %v", err)
		return send, err

	case *json.UnmarshalTypeError:
		lgr.Printf("[WARN] unmarshal")
		return send, nil
	}
}

func logIfErrorApiResponse(send *tgbotapi.APIResponse, err error) (*tgbotapi.APIResponse, error) {
	if err == nil {
		return send, nil
	}
	switch err.(type) {
	default:
		lgr.Printf("[ERROR] cannot send, %v", err)
		return send, err

	case *json.UnmarshalTypeError:
		lgr.Printf("[WARN] unmarshal")
		return send, nil
	}
}
