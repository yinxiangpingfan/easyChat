package config

import (
	"fmt"

	"github.com/spf13/viper"
)

var config *Config

type Config struct {
	Server   ServerConfig   `mapstructure:",squash"`
	Database DatabaseConfig `mapstructure:",squash"`
}

type ServerConfig struct {
	Port string `mapstructure:"SERVER_PORT"`
}

type DatabaseConfig struct {
	Host string `mapstructure:"DATABASE_HOST"`
	Port string `mapstructure:"DATABASE_PORT"`
	User string `mapstructure:"DATABASE_USER"`
	Name string `mapstructure:"DATABASE_NAME"`
}

// 从环境变量加载配置
func LoadConfigEnv() error {
	// 从环境变量加载配置
	viper.AutomaticEnv()
	config = &Config{
		Server: ServerConfig{
			Port: viper.GetString("SERVER_PORT"),
		},
		Database: DatabaseConfig{
			Host: viper.GetString("DATABASE_HOST"),
			Port: viper.GetString("DATABASE_PORT"),
			User: viper.GetString("DATABASE_USER"),
			Name: viper.GetString("DATABASE_NAME"),
		},
	}
	return nil
}

// 从文件加载配置
func LoadConfigFile() error {
	viper.SetConfigFile("local.env")
	if err := viper.ReadInConfig(); err != nil {
		err = fmt.Errorf("读取配置文件失败: %w", err)
		return err
	}
	if err := viper.Unmarshal(&config); err != nil {
		return fmt.Errorf("解析配置文件失败: %w", err)
	}
	return nil
}

func GetConfig(method int) *Config {
	if config == nil {
		if method == 0 {
			_ = LoadConfigEnv()
		} else {
			e := LoadConfigFile()
			fmt.Println(e)
		}
		if config == nil {
			panic("config is nil")
		}
	}
	return config
}
