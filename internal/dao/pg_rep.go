package dao

import (
	"fmt"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"market-bot/sdk/tgbot"
)

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

type UserRepository interface {
	SaveUser(user tgbot.User) error
	UpdateUser(user tgbot.User) error
	GetChatInfoByActiveChainAndStep(activeChain string, activeChainStep string) (chat tgbot.ChatInfo, err error)
	GetRandomUser() (tgbot.User, error)
}

func (r *Repository) GetChatInfo(chatId int64) (chat tgbot.ChatInfo, err error) {
	row := r.db.QueryRowx("SELECT * FROM chat_info WHERE chat_id = $1", chatId)

	if err = row.StructScan(&chat); err != nil {
		return tgbot.ChatInfo{}, errors.Wrapf(err, "unable to get chatInfo, chatId: %d", chatId)
	}
	return
}

func (r *Repository) GetChatInfoByActiveChainAndStep(activeChain string, activeChainStep string) (chat tgbot.ChatInfo, err error) {

	row := r.db.QueryRowx("SELECT * FROM chat_info WHERE active_chain = $1 AND active_chain_step = $2", activeChain, activeChainStep)

	if err = row.StructScan(&chat); err != nil {
		return tgbot.ChatInfo{}, errors.Wrapf(err, "unable to get chatInfo, active_chain: %v, active_chain_step: %v", activeChain, activeChainStep)
	}
	return chat, nil
}

func (r *Repository) SaveChatInfo(chat tgbot.ChatInfo) error {

	insert := `INSERT INTO chat_info (chat_id, active_chain, active_chain_step, chain_data)
								VALUES (:chat_id, :active_chain, :active_chain_step, :chain_data)
								ON CONFLICT (chat_id) DO UPDATE SET active_chain      = :active_chain,
														 			active_chain_step = :active_chain_step,
														  			chain_data        = :chain_data`

	if _, err := r.db.NamedExec(insert, chat); err != nil {
		return errors.Wrap(err, "unable to save chatInfo")
	}
	return nil
}

func (r *Repository) GetButton(btnId string) (btn tgbot.Button, err error) {
	row := r.db.QueryRowx("SELECT * FROM button WHERE id = $1", btnId)

	if err = row.StructScan(&btn); err != nil {
		return tgbot.Button{}, fmt.Errorf("unable to get button, btnId: %s, %w", btnId, err)
	}
	return
}

func (r *Repository) SaveButton(button tgbot.Button) error {
	insert := "INSERT INTO button (id, action, data, created_date) VALUES (:id, :action, :data, now())"

	if _, err := r.db.NamedExec(insert, button); err != nil {
		return err
	}
	return nil
}

func (r *Repository) SaveUser(user tgbot.User) error {
	insert := `INSERT INTO profile (user_id, user_name, last_name, display_name, phone, status, district) VALUES (:user_id, :user_name, :last_name, :display_name, :phone, :status, :district)`
	if _, err := r.db.NamedExec(insert, user); err != nil {
		return err
	}
	return nil
}

func (r *Repository) UpdateUser(user tgbot.User) error {
	query := `UPDATE profile SET status = :status WHERE user_id = :user_id`
	if _, err := r.db.NamedExec(query, user); err != nil {
		return err
	}
	return nil
}

func (r *Repository) UpsertUser(user tgbot.User) (tgbot.User, error) {
	oldUser, _ := r.getUserByUserNameOrTelegramId(user.UserName, *user.UserId)

	var query string
	if oldUser.UserId == nil {
		//language=sql
		query = `INSERT INTO profile (user_id, user_name, last_name, display_name, phone, status, district) VALUES (:user_id, :user_name, :last_name, :display_name, :phone, :status, :district)
					ON CONFLICT (user_name) 
					DO UPDATE SET display_name = :display_name, user_id = :user_id, last_name = :last_name`

	} else {
		//language=sql
		query = `INSERT INTO profile (user_id, user_name, last_name, display_name, phone, status, district) VALUES (:user_id, :user_name, :last_name, :display_name, :phone, :status, :district)
					ON CONFLICT (user_id) 
					DO UPDATE SET user_name = :user_name, display_name = :display_name, last_name = :last_name`
	}

	if _, err := r.db.NamedExec(query, user); err != nil {
		return tgbot.User{}, err
	}
	return r.GetUser(*user.UserId)
}

func (r *Repository) GetUser(userId int64) (tgbot.User, error) {
	row := r.db.QueryRowx("SELECT * FROM profile WHERE user_id = $1", userId)

	var user = tgbot.User{}
	if err := row.StructScan(&user); err != nil {
		return tgbot.User{}, errors.Wrapf(err, "unable to get GetUser, userId: %v", user)
	}
	return user, nil
}

func (r *Repository) getUserByUserNameOrTelegramId(userName string, telegramId int64) (tgbot.User, error) {
	row := r.db.QueryRowx("SELECT * FROM profile WHERE user_id = $1 OR user_name = $2", telegramId, userName)

	var user = tgbot.User{}
	if err := row.StructScan(&user); err != nil {
		return tgbot.User{}, errors.Wrapf(err, "unable to get GetUser, userId: %v", user)
	}
	return user, nil
}

func (r *Repository) GetRandomUser() (tgbot.User, error) {
	row := r.db.QueryRowx("SELECT * FROM profile ORDER BY random() LIMIT 1")

	var user tgbot.User
	if err := row.StructScan(&user); err != nil {
		return user, errors.Wrapf(err, "unable to get random user, userId: %v", user)
	}
	return user, nil
}
