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
		var smallCorrect []string
		var smallIncorrectLarge []string
		var smallIncorrectExpired []string

		for _, afid1 := range smallFiles {
			if utils.IsExpiredFile(afid1) {
				smallIncorrectExpired = append(smallIncorrectExpired, afid1)
			} else if utils.IsLargeFile(afid1) {
				smallIncorrectLarge = append(smallIncorrectLarge, afid1)
			} else {
				smallCorrect = append(smallCorrect, afid1)
			}
		}

		// 输出结果到文件
		err := utils.WriteAfidList(utils.Conf.WriteFileAddress+"/rfsData_correct.txt", smallCorrect)
		if err != nil {
			fmt.Println("无法写入小文件正确列表:", err)
		}
		err = utils.WriteAfidList(utils.Conf.WriteFileAddress+"/rfsData_incorrect_large.txt", smallIncorrectLarge)
		if err != nil {
			fmt.Println("无法写入小文件过大列表:", err)
		}
		err = utils.WriteAfidList(utils.Conf.WriteFileAddress+"/rfsData_incorrect_expired.txt", smallIncorrectExpired)
		if err != nil {
			fmt.Println("无法写入小文件过期列表:", err)
		}
	}()

	go func() {
		defer wg.Done()

		// 分类文件
		var largeCorrect []string
		var largeIncorrectSmall []string
		var largeIncorrectExpired []string

		for _, afid2 := range largeFiles {
			if utils.IsExpiredFile(afid2) {
				largeIncorrectExpired = append(largeIncorrectExpired, afid2)
			} else if utils.IsSmallFile(afid2) {
				largeIncorrectSmall = append(largeIncorrectSmall, afid2)
			} else {
				largeCorrect = append(largeCorrect, afid2)
			}
		}

		// 输出结果到文件
		err := utils.WriteAfidList(utils.Conf.WriteFileAddress+"/raw_correct.txt", largeCorrect)
		if err != nil {
			fmt.Println("无法写入大文件正确列表:", err)
		}
		err = utils.WriteAfidList(utils.Conf.WriteFileAddress+"/raw_incorrect_small.txt", largeIncorrectSmall)
		if err != nil {
			fmt.Println("无法写入大文件过小列表:", err)
		}
		err = utils.WriteAfidList(utils.Conf.WriteFileAddress+"/raw_incorrect_expired.txt", largeIncorrectExpired)
		if err != nil {
			fmt.Println("无法写入大文件过期列表:", err)
		}
	}()

	go func() {
		defer wg.Done()

		// 分类文件
		var expiredCorrect []string
		var expiredIncorrectSmall []string
		var expiredIncorrectLarge []string

		for _, afid3 := range expiredFiles {
			if utils.IsExpiredFile(afid3) {
				expiredCorrect = append(expiredCorrect, afid3)
			} else if utils.IsSmallFile(afid3) {
				expiredIncorrectSmall = append(expiredIncorrectSmall, afid3)
			} else {
				expiredIncorrectLarge = append(expiredIncorrectLarge, afid3)
			}
		}

		// 输出结果到文件
		err := utils.WriteAfidList(utils.Conf.WriteFileAddress+"/expired_correct.txt", expiredCorrect)
		if err != nil {
			fmt.Println("无法写入过期文件正确列表:", err)
		}
		err = utils.WriteAfidList(utils.Conf.WriteFileAddress+"/expired_incorrect_small.txt", expiredIncorrectSmall)
		if err != nil {
			fmt.Println("无法写入没过期小文件列表:", err)
		}
		err = utils.WriteAfidList(utils.Conf.WriteFileAddress+"/expired_incorrect_large.txt", expiredIncorrectLarge)
		if err != nil {
			fmt.Println("无法写入没过期大文件列表:", err)
		}
	}()

	go func() {
		defer wg.Done()

		// 分类文件
		var expiredLargeExpired []string
		var expiredLargeUnexpiredSmall []string
		var expiredLargeUnexpiredLarge []string

		for _, afid4 := range expiredLargeFiles {
			if utils.IsExpiredFile(afid4) {
				expiredLargeExpired = append(expiredLargeExpired, afid4)
			} else if utils.IsSmallFile(afid4) {
				expiredLargeUnexpiredSmall = append(expiredLargeUnexpiredSmall, afid4)
			} else {
				expiredLargeUnexpiredLarge = append(expiredLargeUnexpiredLarge, afid4)
			}
		}

		// 输出结果到文件
		err := utils.WriteAfidList(utils.Conf.WriteFileAddress+"/expiredLarge_expired.txt", expiredLargeExpired)
		if err != nil {
			fmt.Println("无法写入过期大文件过期列表:", err)
		}
		err = utils.WriteAfidList(utils.Conf.WriteFileAddress+"/expiredLarge_unexpired_small.txt", expiredLargeUnexpiredSmall)
		if err != nil {
			fmt.Println("无法写入过期大文件未过期小文件列表:", err)
		}
		err = utils.WriteAfidList(utils.Conf.WriteFileAddress+"/expiredLarge_unexpired_large.txt", expiredLargeUnexpiredLarge)
		if err != nil {
			fmt.Println("无法写入过期大文件未过期大文件列表:", err)
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
