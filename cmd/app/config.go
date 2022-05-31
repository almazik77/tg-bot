package main

import (
	"fmt"
	"github.com/caarlos0/env/v6"
	"github.com/go-pkgz/lgr"
	"net/url"
	"strings"
	"time"
)

type pg struct {
	User            string        `env:"USER,notEmpty"  envExpand:"true" envDefault:"postgres"`
	Pass            string        `env:"PASSWORD,notEmpty"  envExpand:"true" envDefault:"postgres"`
	Host            string        `env:"ADDR,notEmpty"  envExpand:"true" envDefault:"localhost"`
	Port            int           `env:"PORT,notEmpty"  envExpand:"true"  envDefault:"5432"`
	Db              string        `env:"DATABASE,notEmpty"  envExpand:"true" envDefault:"technokratos-bot"`
	Params          string        `env:"PARAMS,notEmpty" envDefault:"sslmode=disable&application_name=technokratos-bot"  envExpand:"true"`
	MaxOpenConn     int           `env:"MAX_OPEN_CONN" envDefault:"10"`
	MaxIdleConn     int           `env:"MAX_IDLE_CONN" envDefault:"0"`
	MaxLifeTime     time.Duration `env:"MAX_LIFE_TIME" envDefault:"30m"`
	MaxIdleTime     time.Duration `env:"MAX_IDLE_TIME" envDefault:"1m"`
	PoolConnTimeout time.Duration `env:"POOL_CONNECTION_TIMEOUT" envDefault:"1m"`
}

//goland:noinspection SqlResolve
var conf struct {
	TgToken string `env:"TG_TOKEN,notEmpty"  envExpand:"true" envDefault:""`

	Pg pg `envPrefix:"DB_"`

	LogLevel  string `env:"LOG_LEVEL" envDefault:"debug"`
	LogFormat string `env:"LOG_FORMAT" envDefault:"logstash"`
	Dry       bool   `env:"DRY" envDefault:"false"`

	SuperUsers []string `env:"SUPER_USER" envSeparator:":" envDefault:""`
}

func InitConfig() {
	conf := &conf
	if err := env.Parse(conf); err != nil {
		lgr.Fatalf("[ERROR] Unable to init config, %v", err)
	}
}

func GetPgDsn() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?%s",
		url.QueryEscape(conf.Pg.User),
		url.QueryEscape(conf.Pg.Pass),
		conf.Pg.Host,
		conf.Pg.Port,
		conf.Pg.Db,
		conf.Pg.Params)
}

func DsnMaskPass(dsn string) string {
	at := strings.Index(dsn, "@")
	beforeAt := dsn[:at]
	colon := strings.LastIndex(beforeAt, ":")
	beforeColon := dsn[:colon+1]
	afterAt := dsn[at:]
	return beforeColon + "********" + afterAt
}
