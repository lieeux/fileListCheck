package utils

import (
	"fmt"
	"github.com/spf13/viper"
)

type Config struct {
	ReadFileAddress      string `mapstructure:"read_file_address"`
	WriteFileAddress     string `mapstructure:"write_file_address"`
	SeedApiUrl           string `mapstructure:"seed_api_url"`
	PnUrl                string `mapstructure:"pn_url"`
	SmallFileName        string `mapstructure:"small_file_name"`
	LargeFileName        string `mapstructure:"large_file_name"`
	ExpiredFileName      string `mapstructure:"expired_file_name"`
	ExpiredLargeFileName string `mapstructure:"expiredlarge_file_name"`
	BufferChan           int    `mapstructure:"bufferChan_number"`
	Worker               int    `mapstructure:"worker_number"`
}

var Conf Config

func LoadConfig() {

	// 读取配置文件
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Failed to read config file: %s", err))
	}

	// 解析配置文件到结构体
	//var Conf Config
	err = viper.Unmarshal(&Conf)
	if err != nil {
		panic(fmt.Errorf("Failed to parse config file: %s", err))
	}
}
