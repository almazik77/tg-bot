package model

import (
	"github.com/google/uuid"
	"time"
)

type EmployeeStatus string

const (
	NewEmployee     = EmployeeStatus("NEW")
	OnboardEmployee = EmployeeStatus("ONBOARD")
	ActiveEmployee  = EmployeeStatus("ACTIVE")
	FiredEmployee   = EmployeeStatus("FIRED")
)

type District string

const (
	Java     = District("JAVA")
	Python   = District("PYTHON")
	Android  = District("ANDROID")
	IOS      = District("IOS")
	Front    = District("FRONTEND")
	QA       = District("QA")
	Design   = District("DESIGN")
	Analysts = District("ANALYSTS")
	Devops   = District("DEVOPS")
	PM       = District("PM")
	HR       = District("HR")
	Common   = District("COMMON")
)

type Employee struct {
	Id           uuid.UUID      `db:"id"`
	Status       EmployeeStatus `db:"status"`
	FirstName    string         `db:"first_name"`
	LastName     *string        `db:"last_name"`
	Phone        *string        `db:"phone"`
	TelegramId   int64          `db:"telegram_id"`
	PhotoId      *string        `db:"photo_id"`
	StartMessage *string        `db:"start_message"`
	District     District       `db:"district"`
	CreatedDate  time.Time      `db:"created_date"`
}

type Task struct {
	Id          uuid.UUID `db:"id"`
	Name        string    `db:"name"`
	Url         string    `db:"url"`
	RoomId      uuid.UUID `db:"room_id"`
	Grade       int32     `db:"grade"`
	Finished    bool      `db:"finished"`
	CreatedDate time.Time `db:"created_date"`
}

type Chat struct {
	TelegramId         int64    `db:"telegram_id"`
	Description        *string  `db:"description"`
	Title              string   `db:"title"`
	District           District `db:"district"`
	Active             bool     `db:"active"`
	Required           bool     `db:"required"`
	CanInviteUsers     bool     `db:"can_invite_users"`
	CanRestrictMembers bool     `db:"can_restrict_users"`
}

type Nomination struct {
	Id      int64     `db:"id"`
	Message string    `db:"message"`
	Date    time.Time `db:"create_date"`
}

type Question struct {
	Id             uuid.UUID `db:"id"`
	Title          string    `db:"title"`
	Description    string    `db:"description"`
	Answer         string    `db:"answer"`
	AuthorUserId   int64     `db:"author_user_id"`
	AnswererUserId *int64    `db:"answerer_user_id"`
	Status         string    `db:"status"`
	CreatedDate    time.Time `db:"created_date"`
}

// SysCommand hold one type Triggers from basic.data
type SysCommand struct {
	Triggers []string
	Question string
	Message  string
}
