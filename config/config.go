package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	Server         ServerConfig   `mapstructure:",squash"`
	Database       DatabaseConfig `mapstructure:",squash"`
	Redis          RedisConfig    `mapstructure:",squash"`
	Kafka          KafkaConfig    `mapstructure:",squash"`
	ConfigSettings ConfigSetting  `mapstructure:",squash"`
}

type ServerConfig struct {
	Port string `mapstructure:"SERVER_PORT"`
}

type ConfigSetting struct {
	Settings string `mapstructure:"CONFIG_SETTING"`
}

type DatabaseConfig struct {
	Host     string `mapstructure:"DATABASE_HOST"`
	Port     string `mapstructure:"DATABASE_PORT"`
	Password string `mapstructure:"DATABASE_PASSWORD"`
	User     string `mapstructure:"DATABASE_USER"`
	Name     string `mapstructure:"DATABASE_NAME"`
}

type RedisConfig struct {
	Host string `mapstructure:"REDIS_HOST"`
	Port string `mapstructure:"REDIS_PORT"`
}

type KafkaConfig struct {
	Port string `mapstructure:"KAFKA_PORT"`
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
		Redis: RedisConfig{
			Host: viper.GetString("REDIS_HOST"),
			Port: viper.GetString("REDIS_PORT"),
		},
		Kafka: KafkaConfig{
			Port: viper.GetString("KAFKA_PORT"),
		},
		ConfigSettings: ConfigSetting{
			Settings: viper.GetString("CONFIG_SETTING"),
		},
	}
	return nil
}

// 从文件加载配置
func LoadConfigFile(path string) error {
	viper.SetConfigFile(path)
	if err := viper.ReadInConfig(); err != nil {
		err = fmt.Errorf("读取配置文件失败: %w", err)
		return err
	}
	if err := viper.Unmarshal(&config); err != nil {
		return fmt.Errorf("解析配置文件失败: %w", err)
	}
	return nil
}

var config *Config

func GetConfig(path string) *Config {
	var err error
	// 从环境变量加载配置
	_ = LoadConfigEnv()
	if config.ConfigSettings.Settings == "file" || config.ConfigSettings.Settings == "" {
		//如果没有本地环境变量或者设定从文件获取变量
		err = LoadConfigFile(path)
		if err != nil {
			panic(err)
		}
	}

	if config == nil || config.ConfigSettings.Settings == "" {
		panic("读取配置失败")
	}

	return config
}
