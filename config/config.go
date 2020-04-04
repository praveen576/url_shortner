package config

import (
	"github.com/spf13/viper"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const (
	AppName = "url_shortner"
)

var Config = struct {
	Env      string
	Dir      string
	Hostname string
	Pid      string

	LogFormat string `mapstructure:"log_format"`

	Server struct {
		Network string `mapstructure:"network"`
		Addr    string `mapstructure:"addr"`
	} `mapstructure:"server"`

	Logger struct {
		Level string `mapstructure:"level"`
	}

	Redis struct {
		Addr string `mapstructure:"addr"`
		Db   int    `mapstructure:"db"`
		Pool int    `mapstructure:"pool"`
	} `mapstructure:"redis"`

	InfluxDB struct {
		Addr      string  `mapstructure:"addr"`
		Username  string  `mapstructure:"username"`
		Password  string  `mapstructure:"password"`
		Database  string  `mapstructure:"database"`
		Precision string  `mapstructure:"precision"`
		Interval  float32 `mapstructure:"interval"`
	} `mapstructure:"influxdb"`
}{}

func ReadConfig() {
	dir, err := os.Getwd()

	if err != nil {
		log.Fatal(err)
	}

	hostname, err := os.Hostname()
	if err != nil {
		log.Fatal(err)
	}

	env := os.Getenv(strings.ToUpper(AppName) + "_ENV")
	if env == "" {
		env = "development"
	}

	log.Println("Loading configurations for", env)

	viper.SetEnvPrefix(AppName)
	viper.AddConfigPath(filepath.Join(dir, "config", "env"))
	viper.SetConfigType("yaml")
	viper.SetConfigName(env)
	err = viper.ReadInConfig()
	if err != nil {
		log.Fatal(err)
	}
	err = viper.Unmarshal(&Config)
	if err != nil {
		log.Fatal(err)
	}
	Config.Env = env
	Config.Dir = dir
	Config.Hostname = hostname
	Config.Pid = strconv.Itoa(os.Getpid())

	log.Println(Config)
}
