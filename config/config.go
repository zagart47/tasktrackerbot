package config

import (
	"fmt"
	"log"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type cfg struct {
	MigrationsPath string        `yaml:"migration"`
	Timeout        time.Duration `yaml:"timeout" envDefault:"15s"`
	Bot            struct {
		Token string `yaml:"token"`
	} `yaml:"bot"`
	Postgres struct {
		DSN      string
		Host     string `yaml:"host" env-default:"postgres"`
		Port     string `yaml:"port" env-default:"5432"`
		DBName   string `yaml:"db_name" env-default:"postgres"`
		UserName string `yaml:"user_name" env-default:"postgres"`
		Pwd      string `yaml:"pwd" env-default:"postgres"`
	} `yaml:"postgres"`
	Redis struct {
		Host string `yaml:"host" env-default:"redis"`
		Port string `yaml:"port" env-default:"6379"`
		Pwd  string `yaml:"pwd"`
		Name string `yaml:"name" env-default:"redis"`
	} `yaml:"redis"`
}

func config() cfg {
	configs := cfg{}
	if err := cleanenv.ReadConfig("./config/config.yaml", &configs); err != nil {
		log.Println("cannot read configs")
	}
	configs.Postgres.DSN = fmt.Sprintf("%s://%s:%s@%s:%s/%s?sslmode=disable", configs.Postgres.DBName, configs.Postgres.UserName,
		configs.Postgres.Pwd, configs.Postgres.Host, configs.Postgres.Port, configs.Postgres.DBName)
	return configs
}

var Configs = config()
