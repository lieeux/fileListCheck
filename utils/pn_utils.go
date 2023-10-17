package utils

import (
	"fmt"
	"gitlab.paradeum.com/bfs-public/bfs-sdk/client/bfs"
	"strconv"
)

// 判断是否为过期文件
func IsExpiredFile(afid string) bool {
	bfsClient, err := bfs.NewClient(Conf.PnUrl, &bfs.Config{Debug: true})
	res, err := bfsClient.DgstIsExist(afid)
	if err == nil && res.Data.IsExist { //存在则没过期
		return false
	} else {
		return true
	}
}

// 判断是否为小文件
func IsSmallFile(afid string) bool {
	bfsClient, err := bfs.NewClient(Conf.PnUrl, &bfs.Config{Debug: true})
	res, err := bfsClient.ReadParamBfs(afid, "file_length")
	if res.Code != 0 {
		fmt.Println(res.Code, res.Msg)
	}
	length, _ := strconv.Atoi(res.Data.Value)
	if err == nil && length <= 10485760 {
		return true
	} else {
		return false
	}
}

// 判断是否为大文件
func IsLargeFile(afid string) bool {
	bfsClient, err := bfs.NewClient(Conf.PnUrl, &bfs.Config{Debug: true})
	res, err := bfsClient.ReadParamBfs(afid, "file_length")
	if res.Code != 0 {
		fmt.Println(res.Code, res.Msg)
	}
	length, _ := strconv.Atoi(res.Data.Value)
	if err == nil && length > 10485760 {
		return true
	} else {
		return false
	}
}
