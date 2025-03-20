package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

// 请替换成你的Cookie和要签到的贴吧名称
const (
	USER_COOKIE = "BDUSS=xxxxxxxx;"
	TIEBA_NAME  = "周杰伦"
)

// 获取tbs的响应结构
type TBSResponse struct {
	IsLogin int    `json:"is_login"`
	TBS     string `json:"tbs"`
}

// 签到响应结构
type SignResponse struct {
	No    int    `json:"no"`
	Error string `json:"error"`
	Data  struct {
		Errmsg string `json:"errmsg"`
	} `json:"data"`
}

func main() {
	client := &http.Client{}

	// 第一步：获取tbs令牌
	tbs, err := getTBS(client)
	if err != nil {l
		fmt.Println("获取tbs失败:", err)
		return
	}

	// 第二步：执行签到
	err = signTieba(client, tbs)
	if err != nil {
		fmt.Println("签到失败:", err)
		return
	}

	fmt.Println("签到成功！")
}

func getTBS(client *http.Client) (string, error) {
	req, _ := http.NewRequest("GET", "http://tieba.baidu.com/dc/common/tbs", nil)
	req.Header.Set("Cookie", USER_COOKIE)

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	var tbsResp TBSResponse
	json.Unmarshal(body, &tbsResp)

	if tbsResp.IsLogin != 1 {
		return "", fmt.Errorf("未登录，请检查Cookie有效性")
	}

	return tbsResp.TBS, nil
}

func signTieba(client *http.Client, tbs string) error {
	form := url.Values{}
	form.Add("ie", "utf-8")
	form.Add("kw", TIEBA_NAME)
	form.Add("tbs", tbs)

	req, _ := http.NewRequest(
		"POST",
		"http://tieba.baidu.com/sign/add",
		strings.NewReader(form.Encode()),
	)

	// 设置请求头
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Cookie", USER_COOKIE)
	req.Header.Set("Referer", "https://tieba.baidu.com/")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	var signResp SignResponse
	json.Unmarshal(body, &signResp)

	if signResp.No != 0 {
		return fmt.Errorf("错误码：%d，错误信息：%s", signResp.No, signResp.Error)
	}

	return nil
}
