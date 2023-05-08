package utils

import "github.com/spf13/viper"

type StorageConfig struct {
	Host      string `mapstructure:"HOST"`
	PortHttp  string `mapstructure:"PORTHTTP"`
	SecretKey string `mapstructure:"SECRET_KEY"`
}

// LoadConfig Конструктор для создания StorageConfig, который содержит считанные из .env файла данные.
func LoadConfig(path string) (config StorageConfig, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigFile(".env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
