package main

import (

	"github.com/marsofsnow/frpx/pkg/util/log"
	"io/ioutil"
	"net/http"
	"time"
)
var (
	QueryIpUrl = "https://ipw.cn/api/ip/myip"
)

func getPublicIp() string {

	client:=&http.Client{Timeout: 5*time.Second}

	req,err:=http.NewRequest("GET",QueryIpUrl,nil)
	if err!=nil{
		panic(nil)
	}
	resp,err:=client.Do(req)
	if err!=nil{
		log.Error("获取外网 IP 失败，请检查网络:",err)
		panic(err)
	}
	defer resp.Body.Close()

	// 获取 http response 的 body
	body, _ := ioutil.ReadAll(resp.Body)
	externIP := string(body)
	log.Info("外网ip是%s:",externIP)

	return externIP
}