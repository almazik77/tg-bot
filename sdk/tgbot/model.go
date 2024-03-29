package tgbot

import (
	"encoding/json"
	"errors"
	"time"
)

type Command string

type Action string

type Data map[string]string

func NewData() Data {
	return map[string]string{}
}

func (e *Data) Scan(src interface{}) error {
	switch val := src.(type) {
	case []uint8:
		err := json.Unmarshal(val, e)
		if err != nil {
			return errors.New("Unable to unmarshall data")
		}
	default:
		return errors.New("Invalid type for Data")
	}
	return nil
}

type Button struct {
	Id          string
	Action      Action
	Data        Data
	CreatedDate time.Time `db:"created_date"`
}

func (b Button) HasAction(action Action) bool {
	return b.Action == action
}

func (b Button) GetData(key string) string {
	if b.Data == nil {
		return ""
	}

	return b.Data[key]
}

type ChatInfo struct {
	ChatId          int64  `db:"chat_id"`
	ActiveChain     string `db:"active_chain"`
	ActiveChainStep string `db:"active_chain_step"`
	ChainData       Data   `db:"chain_data"`
}

type Status string

const (
	UnKnow  = Status("UN_KNOW")
	Checked = Status("CHECKED")
	Active  = Status("ACTIVE")
)

type User struct {
	Id          int64  `db:"id"`
	UserId      *int64 `db:"user_id"`
	DisplayName string `db:"display_name"`
	LastName    string `db:"last_name"`
	Phone       string `db:"phone"`
	UserName    string `db:"user_name"`
	Status      Status `db:"status"`
	District    string `db:"district"` // удалить после релиза 1.0
}
