// Package config конфигурация GoMon
package config

import (
	"fmt"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Mailings []Mailing `toml:"mailings"`
	Reports  []Report  `toml:"reports"`
}

type Mailing struct {
	Enable  bool   `toml:"enable"`
	Leave   bool   `toml:"leave"`
	Chat    string `toml:"chat"`
	Message string `toml:"message"`
}

type Report struct {
	Enable  bool   `toml:"enable"`
	Chat    string `toml:"chat"`
	Message string `toml:"message"`
}

func LoadConfig(configPath string) *Config {
	c := New()
	_, err := toml.DecodeFile(configPath, &c)
	if err != nil {
		fmt.Printf("Ошибка при загрузке конфигурации %s\n", err)
		fmt.Println("Применены параметры поумолчанию")
	}
	return c
}

// New инициализация конфигурации приложения
func New() *Config {
	return &Config{}
}
