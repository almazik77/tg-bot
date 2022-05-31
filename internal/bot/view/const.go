package view

import (
	"market-bot/sdk/tgbot"
)

const (
	ActionStart = tgbot.Action("START")
)

const (
	ActionAddUser      = tgbot.Action("ADD_USER")
	ActionDelUser      = tgbot.Action("DEL_USER")
	ActionNotImplement = tgbot.Action("NO")
	ActionNewChat      = tgbot.Action("NEW_CHAT")
	ActionNewQuestion  = tgbot.Action("NEW_QUESTION")

	ActionJoinChat = tgbot.Action("JOIN_CHAT")

	ActionChangeChatActive   = tgbot.Action("CHANGE_CHAT_ACTIVE")
	ActionChangeChatRequired = tgbot.Action("CHANGE_CHAT_REQUIRED")

	ActionChangeQuestionStatus = tgbot.Action("CHANGE_QUESTION_STATUS")
	ActionSetAnswerForQuestion = tgbot.Action("SET_ANSWER_QUESTION")
)
