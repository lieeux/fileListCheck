package main

import (
	"bufio"
	"fileListCheck/utils"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func main() {
	startTime := time.Now() // 记录程序开始时间

	utils.LoadConfig()

	// 创建等待组,使用go协程并发处理任务
	var wg sync.WaitGroup
	wg.Add(4)

	//go func() {
	//	defer wg.Done()
	//	utils.BatchStreamFor(utils.Conf.ReadFileAddress+"/"+utils.Conf.SmallFileName, utils.Conf.WriteFileAddress+"/rfsData_incorrect_expired.txt", utils.Conf.WriteFileAddress+"/rfsData_incorrect_large.txt", utils.Conf.WriteFileAddress+"/rfsData_correct.txt", utils.IsExpiredFile, utils.IsLargeFile, 1000, 10)
	//}()

	go func() {
		defer wg.Done()

		filer, _ := os.Open(utils.Conf.ReadFileAddress + "/" + utils.Conf.SmallFileName) //打开指定文件，返回文件对象
		defer filer.Close()
		scanner := bufio.NewScanner(filer) //创建一个用于读取文件的扫描器

		filew1, _ := os.Create(utils.Conf.WriteFileAddress + "/rfsData_incorrect_expired.txt") //创建一个指定名称的新文件，返回一个文件对象
		defer filew1.Close()
		writer1 := bufio.NewWriter(filew1) //创建一个用于写入文件的缓冲写入器

		filew2, _ := os.Create(utils.Conf.WriteFileAddress + "/rfsData_incorrect_large.txt") //创建一个指定名称的新文件，返回一个文件对象
		defer filew2.Close()
		writer2 := bufio.NewWriter(filew2) //创建一个用于写入文件的缓冲写入器

		filew3, _ := os.Create(utils.Conf.WriteFileAddress + "/rfsData_correct.txt") //创建一个指定名称的新文件，返回一个文件对象
		defer filew3.Close()
		writer3 := bufio.NewWriter(filew3) //创建一个用于写入文件的缓冲写入器

		filew4, _ := os.Create(utils.Conf.WriteFileAddress + "/seed_files.txt") //创建一个指定名称的新文件，返回一个文件对象
		defer filew4.Close()
		writer4 := bufio.NewWriter(filew4) //创建一个用于写入文件的缓冲写入器

		filew5, _ := os.Create(utils.Conf.WriteFileAddress + "/notseed_files.txt") //创建一个指定名称的新文件，返回一个文件对象
		defer filew5.Close()
		writer5 := bufio.NewWriter(filew5) //创建一个用于写入文件的缓冲写入器

		var wg sync.WaitGroup // 创建等待组

		//var workerNum = 5
		// 创建一个退出通道，用于通知 goroutine 退出
		exitCh := make(chan struct{})
		//建立一个缓冲通道，大小为10
		lines := make(chan string, utils.Conf.BufferChan)
		done := make(chan struct{})

		// 创建生产者，从文件中读取afid
		go func() {
			defer close(done)
			utils.ProduceLines(scanner, lines, exitCh)
		}()

		//创建消费者，消费通道中的afid
		for i := 0; i < utils.Conf.Worker; i++ {
			wg.Add(1)
			go utils.ConsumeLinesPro(writer1, writer2, writer3, writer4, writer5, utils.IsExpiredFile, utils.IsLargeFile, lines, &wg)
		}

		//信号通道，接受中断信号
		signals := make(chan os.Signal, 1)
		signal.Notify(signals, os.Interrupt, syscall.SIGTERM)
		select {
		case <-signals:
			// 接收到中断信号，发送退出信号并等待生产者协程完成
			close(exitCh) // 发送退出信号，通知生产者协程退出
			<-done        // 等待生产者协程完成
		case <-done:
			// 生产者协程已完成，无需等待中断信号
		}
		wg.Wait()

		writer1.Flush() //将缓冲区中的数据写入文件
		writer2.Flush() //将缓冲区中的数据写入文件
		writer3.Flush() //将缓冲区中的数据写入文件
		writer4.Flush() //将缓冲区中的数据写入文件
		writer5.Flush() //将缓冲区中的数据写入文件

	}()

	go func() {
		defer wg.Done()
		utils.BatchStreamFor(utils.Conf.ReadFileAddress+"/"+utils.Conf.LargeFileName, utils.Conf.WriteFileAddress+"/raw_incorrect_expired.txt", utils.Conf.WriteFileAddress+"/raw_incorrect_small.txt", utils.Conf.WriteFileAddress+"/raw_correct.txt", utils.IsExpiredFile, utils.IsSmallFile, utils.Conf.BufferChan, utils.Conf.Worker)
	}()

	go func() {
		defer wg.Done()
		utils.BatchStreamFor(utils.Conf.ReadFileAddress+"/"+utils.Conf.ExpiredFileName, utils.Conf.WriteFileAddress+"/expired_correct.txt", utils.Conf.WriteFileAddress+"/expired_incorrect_small.txt", utils.Conf.WriteFileAddress+"/expired_incorrect_large.txt", utils.IsExpiredFile, utils.IsSmallFile, utils.Conf.BufferChan, utils.Conf.Worker)
	}()

	go func() {
		defer wg.Done()
		utils.BatchStreamFor(utils.Conf.ReadFileAddress+"/"+utils.Conf.ExpiredLargeFileName, utils.Conf.WriteFileAddress+"/expiredLarge_expired.txt", utils.Conf.WriteFileAddress+"/expiredLarge_unexpired_small.txt", utils.Conf.WriteFileAddress+"/expiredLarge_unexpired_large.txt", utils.IsExpiredFile, utils.IsSmallFile, utils.Conf.BufferChan, utils.Conf.Worker)
	}()

	// 等待所有协程完成
	wg.Wait()

	endTime := time.Now()                 // 记录程序结束时间
	elapsedTime := endTime.Sub(startTime) // 计算程序运行时间

	fmt.Println("程序执行完成！")
	fmt.Println("程序运行时间:", elapsedTime)
}
