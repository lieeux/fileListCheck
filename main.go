package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"

	"gitlab.paradeum.com/bfs-public/bfs-sdk/client/bfs"
)

const (
	seedAPIURL = "http://118.193.47.85:5143/tn/location/seedid/"
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

var (
	bfsClient, err = bfs.NewClient("https://pnode.solarfs.io", &bfs.Config{Debug: true})
)

// 查询小文件是否为seed文件
func findSeedFiles(afids []string) ([]string, error) { //接受读取的afid切片，返回是seed的afid切片
	var seedFiles []string

	for _, afid := range afids {
		seedAPI := seedAPIURL + afid
		resp, err := http.Get(seedAPI) //使用http.Get向API URL发送GET请求，返回响应对象
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body) //读取响应体并将其存储在名为body的字节数组中
		if err != nil {
			return nil, err
		}

		jsonStr := string(body)

		var respon Response
		err = json.Unmarshal([]byte(jsonStr), &respon)
		if err != nil {
			fmt.Println(err)
		}

		isSeedFile := respon.Data.IsBigFile

		if isSeedFile {
			seedFiles = append(seedFiles, afid)
		}
	}
	return seedFiles, nil
}

// 判断是否为过期文件
func isExpiredFile(afid string) bool {
	res, err := bfsClient.DgstIsExist(afid)
	if err == nil && res.Data.IsExist { //存在则没过期
		return false
	} else {
		return true
	}
}

// 判断是否为非过期小文件
func isSmallFile(afid string) bool {
	res, err := bfsClient.ReadParamBfs(afid, "file_length")
	length, _ := strconv.Atoi(res.Data.Value)
	if err == nil && length <= 10485760 && !isExpiredFile(afid) {
		return true
	} else {
		return false
	}
}

// 判断是否为非过期大文件
func isLargeFile(afid string) bool {
	res, err := bfsClient.ReadParamBfs(afid, "file_length")
	length, _ := strconv.Atoi(res.Data.Value)
	if err == nil && length > 10485760 && !isExpiredFile(afid) {
		return true
	} else {
		return false
	}
}

// 分类文件
func classifyFiles(afids []string, classifyFunc func(string) bool) ([]string, []string) {
	var correct []string
	var incorrect []string

	for _, afid := range afids {
		if classifyFunc(afid) {
			correct = append(correct, afid)
		} else {
			incorrect = append(incorrect, afid)
		}
	}

	return correct, incorrect
}

// 读取afid列表文件
func readAfidList(filename string) ([]string, error) { //接收文件名，返回字符串切片
	file, err := os.Open(filename) //打开指定文件，返回文件对象
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file) //创建一个用于读取文件的扫描器
	var afids []string                //定义一个空字符串切片，用于存储读取到的afid
	for scanner.Scan() {              //循环读取文件中的每一行
		afid := strings.TrimSpace(scanner.Text()) //获取当前行的文本内容，并去除首尾的空格
		if afid != "" {                           //如果当前行不为空，则将其添加到afids切片中
			afids = append(afids, afid) //将当前afid添加到afids切片中
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return afids, nil
}

// 写入afid列表文件
func writeAfidList(filename string, afids []string) error { //接受文件名和afid切片
	file, err := os.Create(filename) //创建一个指定名称的新文件，返回一个文件对象
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file) //创建一个用于写入文件的缓冲写入器
	for _, afid := range afids {    //遍历afids切片
		_, err := writer.WriteString(afid + "\n") //将当前afid写入文件，并在末尾添加一个换行符
		if err != nil {
			return err
		}
	}

	return writer.Flush() //将缓冲区中的数据写入文件
}

func main() {
	// 读取小文件afid列表
	smallFiles, err := readAfidList("c:/resource/afsIndexAfidList.txt")
	if err != nil {
		fmt.Println("无法读取小文件afid列表:", err)
		return
	}

	// 读取大文件afid列表
	largeFiles, err := readAfidList("c:/resource/afsIndexAfidList.txt")
	if err != nil {
		fmt.Println("无法读取大文件afid列表:", err)
		return
	}

	// 读取过期文件afid列表
	expiredFiles, err := readAfidList("c:/resource/afsIndexAfidList.txt")
	if err != nil {
		fmt.Println("无法读取过期文件afid列表:", err)
		return
	}

	// 分类文件
	smallCorrect, smallIncorrect := classifyFiles(smallFiles, isSmallFile)
	largeCorrect, largeIncorrect := classifyFiles(largeFiles, isLargeFile)
	expiredCorrect, expiredIncorrect := classifyFiles(expiredFiles, isExpiredFile)

	// 查询小文件是否为seed文件
	seedFiles, err := findSeedFiles(smallCorrect)
	if err != nil {
		fmt.Println("查询seed文件时发生错误:", err)
		return
	}

	// 输出结果到文件
	err = writeAfidList("c:/resource/small_correct.txt", smallCorrect)
	if err != nil {
		fmt.Println("无法写入小文件afid正确列表:", err)
	}
	err = writeAfidList("c:/resource/small_incorrect.txt", smallIncorrect)
	if err != nil {
		fmt.Println("无法写入小文件afid分类异常列表:", err)
	}
	err = writeAfidList("c:/resource/large_correct.txt", largeCorrect)
	if err != nil {
		fmt.Println("无法写入大文件afid正确列表:", err)
	}
	err = writeAfidList("c:/resource/large_incorrect.txt", largeIncorrect)
	if err != nil {
		fmt.Println("无法写入大文件afid分类异常列表:", err)
	}
	err = writeAfidList("c:/resource/expired_correct.txt", expiredCorrect)
	if err != nil {
		fmt.Println("无法写入过期文件afid正确列表:", err)
	}
	err = writeAfidList("c:/resource/expired_incorrect.txt", expiredIncorrect)
	if err != nil {
		fmt.Println("无法写入过期文件afid异常列表:", err)
	}
	err = writeAfidList("c:/resource/seed_files.txt", seedFiles)
	if err != nil {
		fmt.Println("无法写入seed文件列表:", err)
	}

	fmt.Println("程序执行完成！")
}
