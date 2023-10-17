package main

import (
	"fileListCheck/utils"
	"fmt"
	"sync"
	"time"
)

func main() {
	startTime := time.Now() // 记录程序开始时间

	utils.LoadConfig()

	// 创建等待组
	var wg sync.WaitGroup

	// 使用go协程并发执行读取文件任务
	wg.Add(4)

	var smallFiles []string
	var largeFiles []string
	var expiredFiles []string
	var expiredLargeFiles []string

	go func() {
		defer wg.Done()

		// 读取小文件afid列表
		var err error
		smallFiles, err = utils.ReadAfidList(utils.Conf.ReadFileAddress + "/" + utils.Conf.SmallFileName)
		if err != nil {
			fmt.Println("无法读取小文件afid列表:", err)
		}
	}()

	go func() {
		defer wg.Done()

		// 读取大文件afid列表
		var err error
		largeFiles, err = utils.ReadAfidList(utils.Conf.ReadFileAddress + "/" + utils.Conf.LargeFileName)
		if err != nil {
			fmt.Println("无法读取大文件afid列表:", err)
		}
	}()

	go func() {
		defer wg.Done()

		// 读取过期文件afid列表
		var err error
		expiredFiles, err = utils.ReadAfidList(utils.Conf.ReadFileAddress + "/" + utils.Conf.ExpiredFileName)
		if err != nil {
			fmt.Println("无法读取过期文件afid列表:", err)
		}
	}()

	go func() {
		defer wg.Done()

		// 读取过期大文件afid列表
		var err error
		expiredLargeFiles, err = utils.ReadAfidList(utils.Conf.ReadFileAddress + "/" + utils.Conf.ExpiredLargeFileName)
		if err != nil {
			fmt.Println("无法读取过期大文件afid列表:", err)
		}
	}()

	// 等待所有协程完成读取任务
	wg.Wait()

	// 创建新的等待组,使用go协程并发处理任务
	wg.Add(5)

	go func() {
		defer wg.Done()

		// 分类文件
		smallCorrect, smallIncorrect := utils.ClassifyFiles(smallFiles, utils.IsSmallFile)

		// 输出结果到文件
		err := utils.WriteAfidList(utils.Conf.WriteFileAddress+"/rfsData_correct.txt", smallCorrect)
		if err != nil {
			fmt.Println("无法写入小文件afid正确列表:", err)
		}
		err = utils.WriteAfidList(utils.Conf.WriteFileAddress+"/rfsData_incorrect.txt", smallIncorrect)
		if err != nil {
			fmt.Println("无法写入小文件afid分类异常列表:", err)
		}
	}()

	go func() {
		defer wg.Done()

		// 分类文件
		largeCorrect, largeIncorrect := utils.ClassifyFiles(largeFiles, utils.IsLargeFile)

		// 输出结果到文件
		err := utils.WriteAfidList(utils.Conf.WriteFileAddress+"/raw_correct.txt", largeCorrect)
		if err != nil {
			fmt.Println("无法写入大文件afid正确列表:", err)
		}
		err = utils.WriteAfidList(utils.Conf.WriteFileAddress+"/raw_incorrect.txt", largeIncorrect)
		if err != nil {
			fmt.Println("无法写入大文件afid分类异常列表:", err)
		}
	}()

	go func() {
		defer wg.Done()

		// 分类文件
		expiredCorrect, expiredIncorrect := utils.ClassifyFiles(expiredFiles, utils.IsExpiredFile)

		// 输出结果到文件
		err := utils.WriteAfidList(utils.Conf.WriteFileAddress+"/expired_correct.txt", expiredCorrect)
		if err != nil {
			fmt.Println("无法写入过期文件afid正确列表:", err)
		}
		err = utils.WriteAfidList(utils.Conf.WriteFileAddress+"/expired_incorrect.txt", expiredIncorrect)
		if err != nil {
			fmt.Println("无法写入过期文件afid异常列表:", err)
		}
	}()

	go func() {
		defer wg.Done()

		// 分类文件
		expiredLargeCorrect, expiredLargeIncorrect := utils.ClassifyFiles(expiredLargeFiles, utils.IsExpiredFile)

		// 输出结果到文件
		err := utils.WriteAfidList(utils.Conf.WriteFileAddress+"/expiredLarge_correct.txt", expiredLargeCorrect)
		if err != nil {
			fmt.Println("无法写入过期大文件afid正确列表:", err)
		}
		err = utils.WriteAfidList(utils.Conf.WriteFileAddress+"/expiredLarge_incorrect.txt", expiredLargeIncorrect)
		if err != nil {
			fmt.Println("无法写入过期大文件afid异常列表:", err)
		}
	}()

	go func() {
		defer wg.Done()

		// 查询小文件是否为seed文件
		seedFiles, err := utils.FindSeedFiles(smallFiles)
		if err != nil {
			fmt.Println("查询seed文件时发生错误:", err)
			return
		}

		// 输出结果到文件
		err = utils.WriteAfidList(utils.Conf.WriteFileAddress+"/seed_files.txt", seedFiles)
		if err != nil {
			fmt.Println("无法写入seed文件列表:", err)
		}
	}()

	// 等待所有协程完成
	wg.Wait()

	endTime := time.Now()                 // 记录程序结束时间
	elapsedTime := endTime.Sub(startTime) // 计算程序运行时间

	fmt.Println("程序执行完成！")
	fmt.Println("程序运行时间:", elapsedTime)
}
