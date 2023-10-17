package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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

// 查询小文件是否为seed文件
func FindSeedFiles(afids []string) ([]string, error) { //接受读取的afid切片，返回是seed的afid切片
	var seedFiles []string //声明接收返回值的切片

	for _, afid := range afids {
		seedAPI := Conf.SeedApiUrl + "/tn/location/seedid/" + afid
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
