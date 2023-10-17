package utils

import (
	"bufio"
	"os"
	"strings"
)

// 读取afid列表文件
func ReadAfidList(filename string) ([]string, error) { //接收文件名，返回字符串切片
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
func WriteAfidList(filename string, afids []string) error { //接受文件名和afid切片
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
