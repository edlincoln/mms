package configs

import "context"

type Configuration interface {
	Init(ctx context.Context) error
	Close(ctx context.Context) error
}

var (
	data []Configuration
)

func init() {
	data = make([]Configuration, 0)
	data = append(data, &ViperConfiguration{})
	data = append(data, &ServerConfiguration{})
	data = append(data, &DatabaseConfiguration{})
	data = append(data, &DataExtractor{})
}

func GetConfigs() []Configuration {
	return data
}
