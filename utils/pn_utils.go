package utils

import (
	"gitlab.paradeum.com/bfs-public/bfs-sdk/client/bfs"
	"strconv"
)

//var (
//	bfsClient, err = bfs.NewClient(Conf.PnUrl, &bfs.Config{Debug: true})
//)

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

// 判断是否为非过期小文件
func IsSmallFile(afid string) bool {
	bfsClient, err := bfs.NewClient(Conf.PnUrl, &bfs.Config{Debug: true})
	res, err := bfsClient.ReadParamBfs(afid, "file_length")
	length, _ := strconv.Atoi(res.Data.Value)
	if err == nil && length <= 10485760 && !IsExpiredFile(afid) {
		return true
	} else {
		return false
	}
}

// 判断是否为非过期大文件
func IsLargeFile(afid string) bool {
	bfsClient, err := bfs.NewClient(Conf.PnUrl, &bfs.Config{Debug: true})
	res, err := bfsClient.ReadParamBfs(afid, "file_length")
	length, _ := strconv.Atoi(res.Data.Value)
	if err == nil && length > 10485760 && !IsExpiredFile(afid) {
		return true
	} else {
		return false
	}
}
