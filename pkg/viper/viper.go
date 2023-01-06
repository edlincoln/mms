package viper

import (
	"errors"
	"github.com/spf13/viper"
	"io"
	"strings"
)

type IViperImpl struct {
}

func NewIViper() IViper {
	return IViperImpl{}
}

func (i IViperImpl) Configuration(source io.Reader, extension, fileName string) (*viper.Viper, error) {
	if source == nil {
		return nil, errors.New("nil source reader")
	}

	viperSetup := viper.GetViper()
	viperSetup.SetConfigType(extension)
	viperSetup.SetConfigName(fileName)
	viperSetup.AllowEmptyEnv(true)
	viperSetup.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viperSetup.AutomaticEnv()

	errReadConfig := viper.ReadConfig(source)
	if errReadConfig != nil {
		return nil, errReadConfig
	}

	return viperSetup, nil
}

func (i IViperImpl) Unmarshal(v *viper.Viper, m interface{}) error {
	return v.Unmarshal(&m)
}
