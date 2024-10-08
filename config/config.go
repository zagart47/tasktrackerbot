package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type cfg struct {
	Timeout time.Duration `env:"timeout" envDefault:"15s"`
	Bot     struct {
		Token string `yaml:"token"`
	} `yaml:"bot"`
	PostgreSQL struct {
		DSN      string
		Host     string `yaml:"host" env-default:"localhost"`
		Port     string `yaml:"port" env-default:"5432"`
		DBName   string `yaml:"db_name" env-default:"postgres"`
		UserName string `yaml:"user_name" env-default:"postgres"`
		Pwd      string `yaml:"pwd" env-default:"postgres"`
	} `yaml:"postgreSQL"`
}

func config() cfg {
	configs := cfg{}
	if err := cleanenv.ReadConfig("./config/config.yaml", &configs); err != nil {
		log.Println("cannot read configs")
		os.Exit(1)
	}
	configs.PostgreSQL.DSN = fmt.Sprintf("%s://%s:%s@%s:%s/%s?sslmode=disable", configs.PostgreSQL.DBName, configs.PostgreSQL.UserName,
		configs.PostgreSQL.Pwd, configs.PostgreSQL.Host, configs.PostgreSQL.Port, configs.PostgreSQL.DBName)
	return configs
}

var Configs = config()
