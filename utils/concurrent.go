package utils

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type Response struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		Afid      string `json:"afid"`
		IsExist   bool   `json:"is_exist"`
		IsBigFile bool   `json:"is_big_file"`
		Rnodes    []struct {
			Rnid      string `json:"rnid"`
			Api       string `json:"api"`
			RnHealthy struct {
				Available struct{} `json:"available"`
			} `json:"rn_healthy"`
		} `json:"rnodes"`
		LargeFileAfid string `json:"large_file_afid"`
	} `json:"data"`
}

// 生产数据
func ProduceLines(scanner *bufio.Scanner, lines chan<- string, exitCh <-chan struct{}) {
	defer close(lines)   //在函数返回前关闭通道 lines，这样可以通知接收方该通道已经没有更多的值可以接收了
	for scanner.Scan() { //循环读取文本文件的一行数据
		line := scanner.Text() //将当前读取的文本行赋值给变量 line
		select {
		case <-exitCh: // 接收到退出通道的数据，退出当前函数
			return
		default:
			lines <- line //将当前读取的文本行数据发送到通道 lines 中，以供其它 goroutine 接收和处理
		}
	}
}

// 消费数据
func ConsumeLines(writer1 *bufio.Writer, writer2 *bufio.Writer, writer3 *bufio.Writer, isFunc1 func(string) bool, isFunc2 func(string) bool, lines <-chan string, wg *sync.WaitGroup) {
	defer wg.Done()
	for line := range lines { // 处理行数据
		if isFunc1(line) {
			writer1.WriteString(line + "\n") //将当前afid写入文件，并在末尾添加一个换行符
		} else if isFunc2(line) {
			writer2.WriteString(line + "\n") //将当前afid写入文件，并在末尾添加一个换行符
		} else {
			writer3.WriteString(line + "\n") //将当前afid写入文件，并在末尾添加一个换行符
		}
	}
}

// 消费数据
func ConsumeLinesPro(writer1 *bufio.Writer, writer2 *bufio.Writer, writer3 *bufio.Writer, writer4 *bufio.Writer, writer5 *bufio.Writer, isFunc1 func(string) bool, isFunc2 func(string) bool, lines <-chan string, wg *sync.WaitGroup) {
	defer wg.Done()
	for line := range lines { // 处理行数据
		if isFunc1(line) {
			writer1.WriteString(line + "\n") //将当前afid写入文件，并在末尾添加一个换行符
		} else if isFunc2(line) {
			writer2.WriteString(line + "\n") //将当前afid写入文件，并在末尾添加一个换行符
		} else {
			writer3.WriteString(line + "\n") //将当前afid写入文件，并在末尾添加一个换行符
		}

		seedAPI := Conf.SeedApiUrl + "/tn/location/seedid/" + line
		resp, _ := http.Get(seedAPI) //使用http.Get向API URL发送GET请求，返回响应对象
		defer resp.Body.Close()

		body, _ := io.ReadAll(resp.Body) //读取响应体并将其存储在名为body的字节数组中

		jsonStr := string(body)

		var respon Response
		err := json.Unmarshal([]byte(jsonStr), &respon)
		if err != nil {
			fmt.Println(err)
		}

		isSeedFile := respon.Data.IsBigFile

		if isSeedFile {
			writer4.WriteString(line + "\n") //将当前afid写入文件，并在末尾添加一个换行符
		} else {
			writer5.WriteString(line + "\n") //将当前afid写入文件，并在末尾添加一个换行符
		}
	}
}

// 并发处理
func BatchStreamFor(filename string, filename1 string, filename2 string, filename3 string, Func1 func(string) bool, Func2 func(string) bool, chanCap, workerNum int) error {
	filer, err := os.Open(filename) //打开指定文件，返回文件对象
	if err != nil {
		return err
	}
	defer filer.Close()
	scanner := bufio.NewScanner(filer) //创建一个用于读取文件的扫描器
	if scanner.Err() != nil {
		return err
	}

	filew1, err := os.Create(filename1) //创建一个指定名称的新文件，返回一个文件对象
	if err != nil {
		return err
	}
	defer filew1.Close()
	writer1 := bufio.NewWriter(filew1) //创建一个用于写入文件的缓冲写入器

	filew2, err := os.Create(filename2) //创建一个指定名称的新文件，返回一个文件对象
	if err != nil {
		return err
	}
	defer filew2.Close()
	writer2 := bufio.NewWriter(filew2) //创建一个用于写入文件的缓冲写入器

	filew3, err := os.Create(filename3) //创建一个指定名称的新文件，返回一个文件对象
	if err != nil {
		return err
	}
	defer filew3.Close()
	writer3 := bufio.NewWriter(filew3) //创建一个用于写入文件的缓冲写入器

	var wg sync.WaitGroup // 创建等待组

	//var workerNum = 5
	// 创建一个退出通道，用于通知 goroutine 退出
	exitCh := make(chan struct{})
	//建立一个缓冲通道，大小为10
	lines := make(chan string, chanCap)
	done := make(chan struct{})

	// 创建生产者，从文件中读取afid
	go func() {
		defer close(done)
		ProduceLines(scanner, lines, exitCh)
	}()

	//创建消费者，消费通道中的afid
	for i := 0; i < workerNum; i++ {
		wg.Add(1)
		go ConsumeLines(writer1, writer2, writer3, Func1, Func2, lines, &wg)
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

	return nil
}
