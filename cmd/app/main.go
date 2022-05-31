package main

import (
	"github.com/fatih/color"
	"github.com/go-pkgz/lgr"
	migrate "github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres" //for db migration
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	tgbotapi "github.com/mazanur/telegram-bot-api/v6"
	"market-bot/internal/bot/bot_handler"
	"market-bot/internal/bot/view"
	"market-bot/internal/dao"
	"market-bot/internal/service"
	tgbot "market-bot/sdk/tgbot"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

func main() {

	InitConfig()
	lgr.Printf("[INFO] Super users: " + strings.Join(conf.SuperUsers, ","))

	if conf.Dry {
		lgr.Printf("[INFO] Started in dry mode ok\nnBye!")
		os.Exit(0)
	}

	InitLogger()

	pgDb := PgConnInit()
	pgRepository := dao.NewRepository(pgDb)

	bot, err := initTelegramBotApi(pgRepository)
	if err != nil {
		lgr.Fatalf("[ERROR] unable to start app %v", err)
	}

	botConfig := initBotConfig(conf.SuperUsers)

	employeeService := service.NewEmployeeService(pgRepository)

	viewSender := view.NewView(pgRepository, pgRepository, bot)

	application := bot_handler.NewBotApp(viewSender, employeeService, botConfig, bot)

	//application.SetMyCommands()
	go func() {
		err = bot.StartLongPolling(application.Handle)
		if err != nil {
			lgr.Fatalf("[ERROR] unable to start app, %v", err)
		}
	}()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs
}

func initTelegramBotApi(pgRepository *dao.Repository) (*tgbot.Bot, error) {
	l := lgr.Default()
	_ = tgbotapi.SetLogger(lgr.ToStdLogger(l, conf.LogLevel))

	bot, err := tgbot.NewBot(conf.TgToken, pgRepository)
	if err != nil {
		lgr.Fatalf("[ERROR] unable to start app %v", err)
	}
	lgr.Printf("[INFO] BotName %s", bot.Self.UserName)
	bot.Debug = conf.LogLevel == "debug"
	return bot, err
}

func PgConnInit() *sqlx.DB {
	dsn := GetPgDsn()

	if err := MigrateDB(dsn); err != nil {
		lgr.Printf("[ERROR] Database migration failed: %v", err)
	}
	lgr.Print("[INFO] Database migration succeeded")

	db, err := sqlx.Connect("pgx", dsn)
	if err != nil {
		lgr.Printf("Failed to connect to db. dsn='%s': %s", DsnMaskPass(dsn), err.Error())
	}
	db.SetMaxOpenConns(conf.Pg.MaxOpenConn)
	db.SetMaxIdleConns(conf.Pg.MaxIdleConn)
	db.SetConnMaxLifetime(conf.Pg.MaxLifeTime)
	db.SetConnMaxIdleTime(conf.Pg.MaxIdleTime)
	lgr.Print("[INFO] Connected to db")

	return db
}

func MigrateDB(dsn string) error {
	m, err := migrate.New("file://db/migrations", dsn)
	if err != nil {
		return err
	}
	defer m.Close()

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}
	return nil
}

func InitLogger() {
	setupLog(conf.LogLevel == "debug", "")
}

func setupLog(dbg bool, lf string) {
	colorizer := lgr.Mapper{
		ErrorFunc:  func(s string) string { return color.New(color.FgHiRed).Sprint(s) },
		WarnFunc:   func(s string) string { return color.New(color.FgHiYellow).Sprint(s) },
		InfoFunc:   func(s string) string { return color.New(color.FgHiWhite).Sprint(s) },
		DebugFunc:  func(s string) string { return color.New(color.FgWhite).Sprint(s) },
		CallerFunc: func(s string) string { return color.New(color.FgBlue).Sprint(s) },
		TimeFunc:   func(s string) string { return color.New(color.FgCyan).Sprint(s) },
	}

	var stdout, stderr *os.File
	var err error
	if lf != "" {
		stdout, err = os.OpenFile(lf, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
		if err != nil {
			lgr.Printf("error opening log file: %v", err)
			os.Exit(2)
		}
		stderr = stdout
	} else {
		stdout = os.Stdout
		stderr = nil
	}
	if dbg {
		lgr.Setup(
			lgr.Debug,
			lgr.CallerFile,
			lgr.CallerFunc,
			lgr.Msec,
			lgr.LevelBraces,
			lgr.Out(stdout),
			lgr.Err(stderr),
			lgr.Map(colorizer),
			lgr.StackTraceOnError,
		)
	} else {
		lgr.Setup(lgr.Out(stdout), lgr.Err(stderr), lgr.Map(colorizer), lgr.StackTraceOnError)
	}
	lgr.Printf("[INFO] Logger successfully initialized")
}

func initBotConfig(superUsers []string) *bot_handler.Config {
	cfg := &bot_handler.Config{
		SuperUsers: superUsers,
	}
	return cfg
}
