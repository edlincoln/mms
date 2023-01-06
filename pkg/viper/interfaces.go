package viper

import (
	"embed"
	"github.com/spf13/viper"
	"io"
)

type Application interface {
	ConfigurationViper(model interface{}, resource embed.FS) (*viper.Viper, error)
}

type IViper interface {
	Configuration(source io.Reader, extension, fileName string) (*viper.Viper, error)
	Unmarshal(v *viper.Viper, m interface{}) error
}
