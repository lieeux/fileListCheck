package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

//type Config struct {
//	ReadFileAddress  string `mapstructure:"read_file_address"`
//	WriteFileAddress string `mapstructure:"write_file_address"`
//	SeedApiUrl       string `mapstructure:"seed_api_url"`
//}

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

	//// 读取配置文件
	//viper.SetConfigName("config")
	//viper.SetConfigType("yaml")
	//viper.AddConfigPath(".")
	//err := viper.ReadInConfig()
	//if err != nil {
	//	panic(fmt.Errorf("Failed to read config file: %s", err))
	//}
	//
	//// 解析配置文件到结构体
	//var config Config
	//err = viper.Unmarshal(&config)
	//if err != nil {
	//	panic(fmt.Errorf("Failed to parse config file: %s", err))
	//}

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
