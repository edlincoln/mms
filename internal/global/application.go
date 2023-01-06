package global

import "time"

type Application struct {
	App       App       `mapstructure:"app"`
	Server    Server    `mapstructure:"server"`
	Swagger   Swagger   `mapstructure:"swagger"`
	Database  Database  `mapstructure:"database"`
	Extractor Extractor `mapstructure:"extractor"`
}

type App struct {
	Name        string `mapstructure:"name"`
	Description string `mapstructure:"description"`
	Version     string `mapstructure:"version"`
}

type Server struct {
	Port int `mapstructure:"port"`
}

type Swagger struct {
	Port     int    `mapstructure:"port"`
	Ip       string `mapstructure:"ip"`
	BasePath string `mapstructure:"basePath"`
	Custom   Custom `mapstructure:"custom"`
}

type Custom struct {
	BasePath string `mapstructure:"basePath"`
	Enabled  bool   `mapstructure:"enabled"`
}

type Database struct {
	Connection Connection `mapstructure:"connection"`
}

type Connection struct {
	Host         string `mapstructure:"host"`
	Port         int64  `mapstructure:"port"`
	Dialect      string `mapstructure:"dialect"`
	Schema       string `mapstructure:"schema"`
	User         string `mapstructure:"user"`
	Password     string `mapstructure:"password"`
	DatabaseName string `mapstructure:"databaseName"`
	SSLMode      string `mapstructure:"sslMode"`
}

type Extractor struct {
	Url     string   `mapstructure:"url"`
	Enabled bool     `mapstructure:"enabled"`
	Mocked  bool     `mapstructure:"mocked"`
	Range   int      `mapstructure:"range"`
	Pairs   []string `mapstructure:"pairs"`
	Daily   Daily    `mapstructure:"daily"`
}

type Daily struct {
	MaxPeriod int           `mapstructure:"maxPeriod"`
	Frequency int           `mapstructure:"frequency"`
	Hour      string        `mapstructure:"hour"`
	Retry     int           `mapstructure:"retry"`
	RetryTime time.Duration `mapstructure:"retryTime"`
}
