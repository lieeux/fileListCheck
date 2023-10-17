package main

import (
	"fileListCheck/utils"
	"fmt"
	"github.com/spf13/viper"
	"time"
)

type Config struct {
	ReadFileAddress  string `mapstructure:"read_file_address"`
	WriteFileAddress string `mapstructure:"write_file_address"`
	SeedApiUrl       string `mapstructure:"seed_api_url"`
}

func main() {
	startTime := time.Now() // 记录程序开始时间

	// 读取配置文件
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Failed to read config file: %s", err))
	}

	// 解析配置文件到结构体
	var config Config
	err = viper.Unmarshal(&config)
	if err != nil {
		panic(fmt.Errorf("Failed to parse config file: %s", err))
	}

	// 读取小文件afid列表
	smallFiles, err := utils.ReadAfidList(config.ReadFileAddress + "/afsIndexAfidList.txt")
	if err != nil {
		fmt.Println("无法读取小文件afid列表:", err)
		return
	}

	// 读取大文件afid列表
	largeFiles, err := utils.ReadAfidList(config.ReadFileAddress + "/afsIndexAfidList.txt")
	if err != nil {
		fmt.Println("无法读取大文件afid列表:", err)
		return
	}

	// 读取过期文件afid列表
	expiredFiles, err := utils.ReadAfidList(config.ReadFileAddress + "/afsIndexAfidList.txt")
	if err != nil {
		fmt.Println("无法读取过期文件afid列表:", err)
		return
	}

	// 分类文件
	smallCorrect, smallIncorrect := utils.ClassifyFiles(smallFiles, utils.IsSmallFile)
	largeCorrect, largeIncorrect := utils.ClassifyFiles(largeFiles, utils.IsLargeFile)
	expiredCorrect, expiredIncorrect := utils.ClassifyFiles(expiredFiles, utils.IsExpiredFile)

	// 查询小文件是否为seed文件
	seedFiles, err := utils.FindSeedFiles(smallCorrect)
	if err != nil {
		fmt.Println("查询seed文件时发生错误:", err)
		return
	}

	// 输出结果到文件
	err = utils.WriteAfidList(config.WriteFileAddress+"/rfsData_correct.txt", smallCorrect)
	if err != nil {
		fmt.Println("无法写入小文件afid正确列表:", err)
	}
	err = utils.WriteAfidList(config.WriteFileAddress+"/rfsData_incorrect.txt", smallIncorrect)
	if err != nil {
		fmt.Println("无法写入小文件afid分类异常列表:", err)
	}
	err = utils.WriteAfidList(config.WriteFileAddress+"/raw_correct.txt", largeCorrect)
	if err != nil {
		fmt.Println("无法写入大文件afid正确列表:", err)
	}
	err = utils.WriteAfidList(config.WriteFileAddress+"/raw_incorrect.txt", largeIncorrect)
	if err != nil {
		fmt.Println("无法写入大文件afid分类异常列表:", err)
	}
	err = utils.WriteAfidList(config.WriteFileAddress+"/expired_correct.txt", expiredCorrect)
	if err != nil {
		fmt.Println("无法写入过期文件afid正确列表:", err)
	}
	err = utils.WriteAfidList(config.WriteFileAddress+"/expired_incorrect.txt", expiredIncorrect)
	if err != nil {
		fmt.Println("无法写入过期文件afid异常列表:", err)
	}
	err = utils.WriteAfidList(config.WriteFileAddress+"/seed_files.txt", seedFiles)
	if err != nil {
		fmt.Println("无法写入seed文件列表:", err)
	}

	endTime := time.Now()                 // 记录程序结束时间
	elapsedTime := endTime.Sub(startTime) // 计算程序运行时间

	fmt.Println("程序执行完成！")
	fmt.Println("程序运行时间:", elapsedTime)
}
